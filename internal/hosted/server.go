package hosted

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"golang.org/x/sync/errgroup"
)

type applicationFunc func(
	context.Context,
	*slog.Logger,
	func(context.Context, http.Handler) error,
) error

// Run serves the demo on port 8080, optionally stopping the application when idle.
//
// run owns one application instance, including its engine and database pools.
// It calls serveHTTP with its context and handler, and releases its resources
// before returning. serveHTTP blocks until its context is canceled.
func Run(ctx context.Context, logger *slog.Logger, run applicationFunc) error {
	var idleTimeout time.Duration

	if value := os.Getenv("DOGMA_EXAMPLE_IDLE_TIMEOUT"); value != "" {
		var err error

		idleTimeout, err = time.ParseDuration(value)
		if err != nil || idleTimeout <= 0 {
			return fmt.Errorf("DOGMA_EXAMPLE_IDLE_TIMEOUT must be a positive duration, got %q", value)
		}
	}

	group, ctx := errgroup.WithContext(ctx)

	handler := newHandler(ctx, logger, run, idleTimeout)
	defer handler.close()

	server := &http.Server{
		Addr:        ":8080",
		Handler:     handler,
		ReadTimeout: 5 * time.Second,
		IdleTimeout: idleTimeout,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}
	defer server.Close()

	group.Go(func() error {
		listener, err := net.Listen("tcp", server.Addr)
		if err != nil {
			return err
		}

		logger.Info(
			"HTTP server listening",
			"address",
			listener.Addr(),
		)

		if err := server.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	})

	group.Go(func() error {
		return handler.monitor()
	})

	group.Go(func() error {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		return server.Shutdown(shutdownCtx)
	})

	return group.Wait()
}
