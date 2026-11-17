package idletimeout

import (
	"context"
	"log/slog"
	"net/http"
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
	idleTimeout time.Duration,
	run func(context.Context, func()) error,
	serve func(context.Context, func(http.Handler) http.Handler) error,
) error {
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
