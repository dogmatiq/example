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

// Run coordinates application execution, HTTP serving, and optional idle shutdown.
// run calls ready after initialization, then runs until its context is canceled.
// Its resources must be released before it returns. serve wraps its HTTP handler
// with the supplied middleware before accepting requests.
func Run(
	ctx context.Context,
	logger *slog.Logger,
	run func(context.Context, func()) error,
	serve func(context.Context, func(http.Handler) http.Handler) error,
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
	controller := &controller{
		ctx:   ctx,
		enter: make(chan chan struct{}),
		leave: make(chan struct{}),
	}

	group.Go(func() error {
		return controller.run(logger, idleTimeout, run)
	})

	group.Go(func() error {
		return serve(ctx, controller.wrap)
	})

	return group.Wait()
}
