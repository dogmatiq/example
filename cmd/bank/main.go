package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/internal/idletimeout"
	"github.com/dogmatiq/example/ui"
	"github.com/dogmatiq/example/ui/projections"
	"github.com/dogmatiq/runkit"
	_ "golang.org/x/crypto/x509roots/fallback" // Provide root certificates in scratch containers.
	"golang.org/x/sync/errgroup"
)

func main() {
	// Use one logger for engine, application, and HTTP lifecycle messages.
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	// Stop the engine and HTTP server together on interrupt or termination.
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// See runApplication() below for the Dogma example. This wrapper handles HTTP
	// serving and idle shutdown for the public demo; you can skip its internals.
	err := idletimeout.Run(ctx, logger, runApplication)

	if errors.Is(err, context.Canceled) {
		return
	}

	if err != nil {
		logger.Error("bank stopped", "error", err)

		os.Exit(1)
	}
}

func runApplication(
	ctx context.Context,
	logger *slog.Logger,
	serveHTTP func(context.Context, http.Handler) error,
) error {
	// The UI reads projections of the bank's events from this database.
	projectionDB, err := projections.NewDB(
		ctx,
		os.Getenv("DOGMA_EXAMPLE_PROJECTION_DSN"),
	)
	if err != nil {
		return err
	}
	defer projectionDB.Close()

	// Runkit stores events and pending work using its own connection pool.
	engineDSN := os.Getenv("DOGMA_EXAMPLE_ENGINE_DSN")
	if engineDSN == "" {
		return errors.New("DOGMA_EXAMPLE_ENGINE_DSN is required")
	}

	engineDB, err := sql.Open("pgx", engineDSN)
	if err != nil {
		return err
	}
	defer engineDB.Close()

	// Create the engine's schema before starting message delivery.
	if err := runkit.CreateSchema(ctx, engineDB); err != nil {
		return err
	}

	// The application describes the bank's message handlers and their routes.
	application := &example.App{ReadDB: projectionDB}

	// Create an engine to execute the application and persist its state.
	engine := runkit.New(
		application,
		engineDB,
		runkit.WithLogger(logger),
	)

	// Run the engine in the background to deliver messages to the handlers.
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return engine.Run(ctx)
	})

	// The UI submits commands to the engine and reads the resulting projections.
	uiHandler := &ui.Handler{
		DB:              projectionDB,
		CommandExecutor: engine,
	}

	group.Go(func() error {
		return serveHTTP(ctx, uiHandler)
	})

	// Wait for shutdown before closing either database pool.
	return group.Wait()
}
