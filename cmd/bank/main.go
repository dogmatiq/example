package main

import (
	"context"
	"database/sql"
	"errors"
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

// main handles logging and the process exit status.
//
// Application and HTTP server setup lives in [run] so errors can be returned
// after deferred cleanup. [os.Exit] skips deferred calls, so [main] calls it
// only after [run] returns.
func main() {
	// Use text logs locally and JSON for the public demo's log collector.
	newLogHandler := func() slog.Handler {
		return slog.NewTextHandler(os.Stdout, nil)
	}
	if os.Getenv("DOGMA_EXAMPLE_IDLE_TIMEOUT") != "" {
		newLogHandler = func() slog.Handler {
			return slog.NewJSONHandler(os.Stdout, nil)
		}
	}

	// Use one logger for engine, application, and HTTP lifecycle messages.
	logger := slog.New(newLogHandler())

	if err := run(logger); err != nil {
		logger.Error("bank stopped", "error", err)
		os.Exit(1)
	}
}

// run runs the bank's [Dogma] application and serves its web UI.
//
// [example.App] defines the application's message handlers and routes.
// The [runkit] package provides an engine that executes the application and
// persists its state in PostgreSQL. An HTTP server serves [ui.Handler], which
// submits commands to the engine and reads projections to display the bank's
// state.
//
// [run] blocks until a shutdown signal is received or the engine or HTTP server
// fails. Shutdown signals are handled as a successful exit.
//
// [Dogma]: https://github.com/dogmatiq/dogma
func run(logger *slog.Logger) (err error) {
	// Stop the engine and HTTP server together on interrupt or termination.
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer func() {
		stop()

		// A shutdown signal cancels ctx. Treat cancellation as a successful
		// shutdown, even if it occurs during startup.
		if errors.Is(err, context.Canceled) {
			err = nil
		}
	}()

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

	// Describe the bank's message handlers and routes, then create its engine.
	application := newDogmaApp(projectionDB)
	engine := newEngine(application, engineDB, logger)

	// Connect the UI to the engine and the projections it maintains.
	uiHandler := newHTTPHandler(projectionDB, engine)
	var handler http.Handler = uiHandler

	group, ctx := errgroup.WithContext(ctx)
	runEngine := func() error {
		return engine.Run(ctx)
	}

	// The public demo hosted at https://example.dogmatiq.io uses
	// [idletimeout.Inject] to stop its engine between visits. This optional
	// hosting setup can be skipped when reading the regular Dogma example.
	if idle := os.Getenv("DOGMA_EXAMPLE_IDLE_TIMEOUT"); idle != "" {
		if err := idletimeout.Inject(
			ctx,
			idle,
			&runEngine,
			&uiHandler.CommandExecutor,
			&handler,
			engine,
			func() idletimeout.Engine {
				return newEngine(application, engineDB, logger)
			},
			logger,
			engineDB,
			projectionDB,
		); err != nil {
			return err
		}
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

	// Run the engine alongside the HTTP server that hosts the UI.
	group.Go(runEngine)

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
func newHTTPHandler(projectionDB *sql.DB, executor dogma.CommandExecutor) *ui.Handler {
	return &ui.Handler{
		DB:              projectionDB,
		CommandExecutor: executor,
	}
}
