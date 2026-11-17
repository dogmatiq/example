package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dogmatiq/dogma"
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

	err := run(ctx, logger)

	if errors.Is(err, context.Canceled) {
		return
	}

	if err != nil {
		logger.Error("bank stopped", "error", err)

		os.Exit(1)
	}
}

// run hosts the bank's [Dogma] application using the [runkit] package, which
// provides a reference engine implementation with PostgreSQL storage.
//
// [example.App] defines the application's message handlers and routes.
// The engine executes the application, while an HTTP server serves [ui.Handler].
// The UI submits commands to the engine and reads projections of the bank's
// events to display its state.
//
// run blocks until ctx is canceled or the engine or HTTP server fails.
//
// [Dogma]: https://github.com/dogmatiq/dogma
func run(ctx context.Context, logger *slog.Logger) error {
	// The UI reads projections of the bank's events from this database.
	projectionDB, err := projections.NewDB(
		ctx,
		os.Getenv("DOGMA_EXAMPLE_PROJECTION_DSN"),
	)
	if err != nil {
		return err
	}
	defer projectionDB.Close()

	// The engine stores events and pending work using its own connection pool.
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

	application := newDogmaApp(projectionDB)
	group, ctx := errgroup.WithContext(ctx)

	var (
		executor dogma.CommandExecutor
		manager  *idletimeout.Manager
	)

	// The public demo hosted at https://example.dogmatiq.io opts into idle
	// shutdown. Its handler submits commands through the current-engine proxy.
	if idle := os.Getenv("DOGMA_EXAMPLE_IDLE_TIMEOUT"); idle != "" {
		timeout, err := time.ParseDuration(idle)
		if err != nil || timeout <= 0 {
			return fmt.Errorf("DOGMA_EXAMPLE_IDLE_TIMEOUT must be a positive duration, got %q", idle)
		}

		// Expire idle connections without replacing the app and UI pools.
		projectionDB.SetConnMaxIdleTime(timeout)
		engineDB.SetConnMaxIdleTime(timeout)

		manager = idletimeout.New(ctx, logger, timeout, func() idletimeout.Engine {
			return newEngine(application, engineDB, logger)
		})
		executor = manager
		group.Go(manager.Run)
	} else {
		// The regular example submits commands directly to its engine.
		engine := newEngine(application, engineDB, logger)
		executor = engine
		group.Go(func() error {
			return engine.Run(ctx)
		})
	}

	handler := newHTTPHandler(projectionDB, executor)
	if manager != nil {
		handler = manager.Wrap(handler)
	}

	listenAddress := os.Getenv("DOGMA_EXAMPLE_LISTEN_ADDRESS")
	if listenAddress == "" {
		listenAddress = ":8080"
	}

	server := &http.Server{
		Addr:        listenAddress,
		ReadTimeout: 5 * time.Second,
		IdleTimeout: time.Minute,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			// Event observers require a deadline; bound each request's wait.
			ctx, cancel := context.WithTimeout(
				request.Context(),
				5*time.Second,
			)
			defer cancel()

			handler.ServeHTTP(w, request.WithContext(ctx))
		}),
	}
	defer server.Close()

	// Start accepting HTTP requests and report the actual listening address.
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

	// Drain HTTP requests when the application stops or the process is signaled.
	group.Go(func() error {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		return server.Shutdown(shutdownCtx)
	})

	// Wait for the engine and HTTP server before closing their pools.
	return group.Wait()
}

// newDogmaApp describes the bank's message handlers and their routes.
func newDogmaApp(projectionDB *sql.DB) *example.App {
	return &example.App{ReadDB: projectionDB}
}

// newEngine creates an engine that executes the application and persists its state.
func newEngine(application *example.App, engineDB *sql.DB, logger *slog.Logger) *runkit.Engine {
	return runkit.New(
		application,
		engineDB,
		runkit.WithLogger(logger),
	)
}

// newHTTPHandler connects the UI to its projection database and command executor.
func newHTTPHandler(projectionDB *sql.DB, executor dogma.CommandExecutor) http.Handler {
	return &ui.Handler{
		DB:              projectionDB,
		CommandExecutor: executor,
	}
}
