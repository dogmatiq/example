package idletimeout

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

// controller serializes request activity and application start/stop transitions.
type controller struct {
	ctx   context.Context
	enter chan chan struct{}
	leave chan struct{}
}

func (c *controller) wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		ctx, cancel := context.WithCancel(request.Context())
		defer cancel()

		stop := context.AfterFunc(c.ctx, cancel)
		defer stop()

		entered := make(chan struct{})

		select {
		case <-ctx.Done():
			http.Error(w, "application unavailable", http.StatusServiceUnavailable)

			return
		case c.enter <- entered:
		}
		defer func() {
			select {
			case c.leave <- struct{}{}:
			case <-c.ctx.Done():
			}
		}()

		select {
		case <-ctx.Done():
			http.Error(w, "application unavailable", http.StatusServiceUnavailable)
		case <-entered:
			if ctx.Err() == nil {
				next.ServeHTTP(w, request.WithContext(ctx))
			} else {
				http.Error(w, "application unavailable", http.StatusServiceUnavailable)
			}
		}
	})
}

func (c *controller) run(
	logger *slog.Logger,
	idleTimeout time.Duration,
	run func(context.Context, func()) error,
) error {
	timer := time.NewTimer(idleTimeout)
	defer timer.Stop()

	var (
		active      int
		waiting     []chan struct{}
		ready       <-chan struct{}
		finished    <-chan error
		cancel      context.CancelFunc
		initialized bool
		stopping    bool
	)

	start := func() {
		ctx, stop := context.WithCancel(c.ctx)
		started := make(chan struct{})
		done := make(chan error, 1)

		cancel, ready, finished = stop, started, done
		initialized = false

		go func() {
			defer stop()

			done <- run(ctx, func() {
				close(started)
			})
		}()

		if idleTimeout > 0 && active == 0 {
			timer.Reset(idleTimeout)
		}
	}

	timer.Stop()
	start()

	for {
		select {
		case <-c.ctx.Done():
			if cancel != nil {
				cancel()
				<-finished
			}

			return c.ctx.Err()

		case entered := <-c.enter:
			active++
			timer.Stop()

			if cancel == nil {
				start()
			}

			if initialized && !stopping {
				close(entered)
			} else {
				waiting = append(waiting, entered)
			}

		case <-c.leave:
			active--

			if active == 0 && cancel != nil && !stopping && idleTimeout > 0 {
				timer.Reset(idleTimeout)
			}

		case <-ready:
			ready = nil
			initialized = true

			for _, entered := range waiting {
				close(entered)
			}

			waiting = nil

		case <-timer.C:
			stopping = true
			ready = nil
			cancel()

			logger.Info("application idle; stopping engine and database connections")

		case err := <-finished:
			cancel()
			cancel, ready, finished = nil, nil, nil

			if !stopping || (err != nil && !errors.Is(err, context.Canceled)) {
				if err == nil {
					err = errors.New("application stopped unexpectedly")
				}

				return err
			}

			stopping = false
			initialized = false
			timer.Stop()

			if active > 0 {
				start()
			} else {
				waiting = nil
			}
		}
	}
}
