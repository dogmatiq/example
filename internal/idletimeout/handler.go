package idletimeout

import (
	"context"
	"net/http"

	"github.com/dogmatiq/dogma"
)

// ExecuteCommand submits a command to the current engine, waking it if needed.
// All command options, including event observers, are passed through unchanged.
func (m *Manager) ExecuteCommand(ctx context.Context, command dogma.Command, options ...dogma.ExecuteCommandOption) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	stop := context.AfterFunc(m.ctx, cancel)
	defer stop()

	engine, err := m.acquire(ctx)
	if err != nil {
		return err
	}
	defer m.release()

	return engine.ExecuteCommand(ctx, command, options...)
}

// Wrap keeps the engine running for each HTTP request, including requests that
// only read projections and never submit commands.
func (m *Manager) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		ctx, cancel := context.WithCancel(request.Context())
		defer cancel()

		stop := context.AfterFunc(m.ctx, cancel)
		defer stop()

		if _, err := m.acquire(ctx); err != nil {
			http.Error(w, "application unavailable", http.StatusServiceUnavailable)

			return
		}
		defer m.release()

		next.ServeHTTP(w, request.WithContext(ctx))
	})
}

func (m *Manager) acquire(ctx context.Context) (Engine, error) {
	entered := make(chan Engine, 1)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-m.ctx.Done():
		return nil, m.ctx.Err()
	case m.enter <- entered:
	}

	select {
	case <-ctx.Done():
		m.release()

		return nil, ctx.Err()
	case <-m.ctx.Done():
		m.release()

		return nil, m.ctx.Err()
	case engine := <-entered:
		if err := ctx.Err(); err != nil {
			m.release()

			return nil, err
		}

		return engine, nil
	}
}

func (m *Manager) release() {
	select {
	case m.leave <- struct{}{}:
	case <-m.ctx.Done():
	}
}
