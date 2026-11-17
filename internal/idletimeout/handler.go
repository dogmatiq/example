package idletimeout

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// instance owns an application from startup until its resources are released.
type instance struct {
	ctx      context.Context
	cancel   context.CancelFunc
	ready    chan struct{}
	done     chan struct{}
	handler  http.Handler
	stopping bool
	restart  bool
	requests sync.WaitGroup
}

type handler struct {
	ctx         context.Context
	logger      *slog.Logger
	run         applicationFunc
	idleTimeout time.Duration
	activity    chan struct{}
	failures    chan error

	mu        sync.Mutex
	current   *instance
	active    int
	idleSince time.Time
}

func newHandler(
	ctx context.Context,
	logger *slog.Logger,
	run applicationFunc,
	idleTimeout time.Duration,
) *handler {
	h := &handler{
		ctx:         ctx,
		logger:      logger,
		run:         run,
		idleTimeout: idleTimeout,
		activity:    make(chan struct{}, 1),
		failures:    make(chan error, 1),
	}

	h.start()

	return h
}

func (h *handler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithCancel(request.Context())
	defer cancel()

	h.mu.Lock()
	h.active++
	h.mu.Unlock()
	defer h.release()

	for {
		if ctx.Err() != nil || h.ctx.Err() != nil {
			http.Error(w, "application unavailable", http.StatusServiceUnavailable)

			return
		}

		h.mu.Lock()

		if h.current == nil {
			h.start()
		}

		current := h.current
		stopping := current.stopping
		if !stopping {
			current.requests.Add(1)
		}

		h.mu.Unlock()

		if stopping {
			select {
			case <-ctx.Done():
				http.Error(w, "application unavailable", http.StatusServiceUnavailable)

				return
			case <-current.done:
				if !current.restart {
					http.Error(w, "application stopped", http.StatusServiceUnavailable)

					return
				}

				continue
			}
		}
		defer current.requests.Done()

		select {
		case <-ctx.Done():
			http.Error(w, "application unavailable", http.StatusServiceUnavailable)

			return
		case <-current.done:
			http.Error(w, "application stopped", http.StatusServiceUnavailable)
		case <-current.ready:
			// An engine failure must also cancel requests using that instance.
			stop := context.AfterFunc(current.ctx, cancel)
			defer stop()

			current.handler.ServeHTTP(w, request.WithContext(ctx))
		}

		return
	}
}

// start is called with mu held, except during construction.
func (h *handler) start() {
	ctx, cancel := context.WithCancel(h.ctx)

	current := &instance{
		cancel: cancel,
		ready:  make(chan struct{}),
		done:   make(chan struct{}),
	}

	h.current = current
	h.idleSince = time.Now()

	go func() {
		defer cancel()

		err := h.run(ctx, h.logger, func(serveCtx context.Context, handler http.Handler) error {
			current.ctx = serveCtx
			current.handler = handler

			close(current.ready)

			<-serveCtx.Done()

			// Prevent new requests from entering while this instance drains.
			h.mu.Lock()
			current.stopping = true
			h.mu.Unlock()

			current.requests.Wait()

			return nil
		})

		h.mu.Lock()

		expectedStop := ctx.Err() != nil && (err == nil || errors.Is(err, context.Canceled))
		if expectedStop {
			h.current = nil
		}

		current.restart = expectedStop && h.ctx.Err() == nil

		close(current.done)

		h.mu.Unlock()

		if !expectedStop {
			if err == nil {
				err = errors.New("application stopped unexpectedly")
			}

			h.failures <- err
		}
	}()
}

func (h *handler) release() {
	h.mu.Lock()
	h.active--
	h.idleSince = time.Now()
	h.mu.Unlock()

	select {
	case h.activity <- struct{}{}:
	default:
	}
}

func (h *handler) monitor() error {
	timer := time.NewTimer(h.idleTimeout)
	defer timer.Stop()

	if h.idleTimeout == 0 {
		timer.Stop()
	}

	for {
		select {
		case <-h.ctx.Done():
			return h.ctx.Err()
		case err := <-h.failures:
			return err
		case <-h.activity:
		case <-timer.C:
		}

		h.mu.Lock()

		if h.idleTimeout == 0 || h.active != 0 || h.current == nil || h.current.stopping {
			timer.Stop()
			h.mu.Unlock()

			continue
		}

		remaining := time.Until(h.idleSince.Add(h.idleTimeout))
		if remaining > 0 {
			timer.Reset(remaining)
			h.mu.Unlock()

			continue
		}

		h.current.stopping = true
		h.current.cancel()
		h.mu.Unlock()

		h.logger.Info("application idle; stopping engine and database connections")
	}
}

func (h *handler) close() {
	h.mu.Lock()
	current := h.current

	if current != nil {
		current.stopping = true
		current.cancel()
	}

	h.mu.Unlock()

	if current != nil {
		<-current.done
	}
}
