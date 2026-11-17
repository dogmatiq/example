package idletimeout

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/dogmatiq/dogma"
)

// Engine is the part of a Dogma engine used by the lifetime manager. Command
// execution must support waiting for initialization while [Engine.Run] starts.
type Engine interface {
	dogma.CommandExecutor
	Run(context.Context) error
}

// Manager runs one engine at a time and forwards commands to the current engine.
// Request activity keeps it awake; after an idle shutdown, activity starts a new
// engine. Database connection pools belong to the caller and remain unchanged.
type Manager struct {
	ctx           context.Context
	logger        *slog.Logger
	idleTimeout   time.Duration
	initialEngine Engine
	newEngine     func() Engine
	enter         chan chan Engine
	leave         chan struct{}
}

// New returns a manager whose [Manager.Run] method must be started by the caller.
// engine is used for the first start. newEngine must return a fresh engine
// without starting it for each subsequent wake-up.
func New(ctx context.Context, logger *slog.Logger, idleTimeout time.Duration, engine Engine, newEngine func() Engine) *Manager {
	return &Manager{
		ctx:           ctx,
		logger:        logger,
		idleTimeout:   idleTimeout,
		initialEngine: engine,
		newEngine:     newEngine,
		enter:         make(chan chan Engine),
		leave:         make(chan struct{}),
	}
}

// Run starts the engine and manages its lifetime until the context is canceled
// or an engine fails. It waits for a stopped engine before starting another.
func (m *Manager) Run() error {
	timer := time.NewTimer(m.idleTimeout)
	defer timer.Stop()

	var (
		active   int
		waiting  []chan Engine
		current  Engine
		finished <-chan error
		cancel   context.CancelFunc
		stopping bool
	)

	start := func() {
		ctx, stop := context.WithCancel(m.ctx)
		done := make(chan error, 1)
		engine := m.initialEngine
		m.initialEngine = nil

		if engine == nil {
			engine = m.newEngine()
		}

		current, cancel, finished = engine, stop, done

		go func() {
			defer stop()

			done <- engine.Run(ctx)
		}()

		for _, entered := range waiting {
			entered <- engine
		}

		waiting = nil

		if m.idleTimeout > 0 && active == 0 {
			timer.Reset(m.idleTimeout)
		}
	}

	timer.Stop()
	start()

	for {
		select {
		case <-m.ctx.Done():
			if cancel != nil {
				cancel()
				<-finished
			}

			return m.ctx.Err()

		case entered := <-m.enter:
			active++
			timer.Stop()

			if cancel == nil {
				start()
			}

			if stopping {
				waiting = append(waiting, entered)
			} else {
				entered <- current
			}

		case <-m.leave:
			active--

			if active == 0 && cancel != nil && !stopping && m.idleTimeout > 0 {
				timer.Reset(m.idleTimeout)
			}

		case <-timer.C:
			stopping = true
			cancel()

			m.logger.Info("application idle; stopping engine")

		case err := <-finished:
			cancel()
			current, cancel, finished = nil, nil, nil

			if !stopping || (err != nil && !errors.Is(err, context.Canceled)) {
				if err == nil {
					err = errors.New("engine stopped unexpectedly")
				}

				return err
			}

			stopping = false
			timer.Stop()

			if active > 0 {
				start()
			} else {
				waiting = nil
			}
		}
	}
}
