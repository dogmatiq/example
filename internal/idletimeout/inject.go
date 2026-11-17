package idletimeout

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/dogmatiq/dogma"
)

// Inject adds idle shutdown to the public demo's engine and HTTP handler.
// It replaces runEngine, executor, and handler before the caller starts them.
// The replacement runner manages engine until ctx is canceled or an engine
// fails, using newEngine to create a fresh engine after each idle shutdown.
//
// idle must be a positive duration. Unused database connections expire after
// that duration; the caller retains ownership of the pools and UI handler.
// If idle is invalid, Inject returns an error without changing the arguments.
func Inject(
	ctx context.Context,
	idle string,
	runEngine *func() error,
	executor *dogma.CommandExecutor,
	handler *http.Handler,
	engine Engine,
	newEngine func() Engine,
	logger *slog.Logger,
	databases ...*sql.DB,
) error {
	timeout, err := time.ParseDuration(idle)
	if err != nil || timeout <= 0 {
		return fmt.Errorf("DOGMA_EXAMPLE_IDLE_TIMEOUT must be a positive duration, got %q", idle)
	}

	manager := New(ctx, logger, timeout, engine, newEngine)

	for _, database := range databases {
		database.SetConnMaxIdleTime(timeout)
	}

	*runEngine = manager.Run
	*executor = manager
	*handler = manager.Wrap(*handler)

	return nil
}
