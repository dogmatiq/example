package idletimeout

import (
	"context"
	"fmt"
	"log/slog"
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

// Run manages the application's idle lifecycle while runHTTPServer serves HTTP.
//
// run owns one application instance, including its engine and database pools.
// It passes its context and UI handler to the supplied callback, which blocks
// until its context is canceled. run releases its resources before returning.
func Run(
	ctx context.Context,
	logger *slog.Logger,
	run applicationFunc,
	runHTTPServer func(context.Context, *slog.Logger, http.Handler) error,
) error {
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

	group.Go(func() error {
		return runHTTPServer(ctx, logger, handler)
	})

	group.Go(func() error {
		return handler.monitor()
	})

	return group.Wait()
}
