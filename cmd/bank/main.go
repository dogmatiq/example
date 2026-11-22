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

// main configures logging and reports failures through the process exit status.
//
// Application startup and shutdown are in [run] so its deferred cleanup can
// finish before [main] calls [os.Exit], which skips deferred calls.
func main() {
	// Use text logs locally and JSON for the public demo's log collector.
	var logHandler slog.Handler = slog.NewTextHandler(os.Stdout, nil)
	if os.Getenv("DOGMA_EXAMPLE_IDLE_TIMEOUT") != "" {
		logHandler = slog.NewJSONHandler(os.Stdout, nil)
	}

	// Share this logger throughout the example.
	logger := slog.New(logHandler)

	if err := run(logger); err != nil {
		logger.Error("bank stopped", "error", err)
		os.Exit(1)
	}
}

// run starts the bank's [Dogma] application and the HTTP server for its web UI.
//
// [example.App] defines the application's message handlers and routes.
// An engine from the [runkit] package delivers messages to those handlers and
// stores events and pending work in PostgreSQL. [ui.Handler] submits commands
// to the engine and queries projections to display the bank's state.
//
// [run] returns nil after a shutdown signal, or an error if startup or execution
// fails.
//
// [Dogma]: https://github.com/dogmatiq/dogma
func run(logger *slog.Logger) (err error) {
	// Use a shared context to stop the engine and HTTP server on a shutdown signal.
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer func() {
		stop()

		// A shutdown signal requests a normal exit. Suppress the resulting
		// [context.Canceled] error, including when it interrupts startup.
		if errors.Is(err, context.Canceled) {
			err = nil
		}
	}()

	// Open the database that projection handlers update and the UI reads.
	projectionDB, err := projections.NewDB(
		ctx,
		os.Getenv("DOGMA_EXAMPLE_PROJECTION_DSN"),
	)
	if err != nil {
		return err
	}
	defer projectionDB.Close()

	// Open a separate connection pool for the engine's events and pending work.
	engineDSN := os.Getenv("DOGMA_EXAMPLE_ENGINE_DSN")
	if engineDSN == "" {
		return errors.New("DOGMA_EXAMPLE_ENGINE_DSN is required")
	}

	engineDB, err := sql.Open("pgx", engineDSN)
	if err != nil {
		return err
	}
	defer engineDB.Close()

	// Create the engine's schema before it starts processing messages.
	if err := runkit.CreateSchema(ctx, engineDB); err != nil {
		return err
	}

	// [example.App.Configure] defines the bank's message handlers and routes.
	application := newDogmaApp(projectionDB)

	// Use the [runkit] package to create an engine that executes the application.
	engine := newEngine(application, engineDB, logger)

	// The UI submits commands to the engine and reads the projection database.
	uiHandler := newHTTPHandler(projectionDB, engine)
	var handler http.Handler = uiHandler

	// Share cancellation so an engine or HTTP server failure stops both.
	group, ctx := errgroup.WithContext(ctx)
	runEngine := func() error {
		return engine.Run(ctx)
	}

	// The public demo hosted at https://example.dogmatiq.io uses
	// [idletimeout.Inject] to stop the engine while the UI is idle.
	// You can skip this block when following the application and HTTP server setup.
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

	// Serve process liveness independently of the UI and its engine.
	mux := http.NewServeMux()
	mux.HandleFunc(
		"GET /_internal/healthz",
		func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte("ok\n"))
		},
	)
	mux.Handle("/", handler)

	// Configure the HTTP server that hosts the bank's UI.
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
			// Bound the UI's waits for events and projections. Event observers
			// created with [dogma.WithEventObserver] require a deadline.
			ctx, cancel := context.WithTimeout(
				request.Context(),
				5*time.Second,
			)
			defer cancel()

			mux.ServeHTTP(w, request.WithContext(ctx))
		}),
	}
	defer server.Close()

	// Start the engine so it can accept commands and deliver messages.
	group.Go(runEngine)

	// Start serving the UI and report its listening address.
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

	// Allow active HTTP requests to finish when shutdown is requested.
	group.Go(func() error {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		return server.Shutdown(shutdownCtx)
	})

	// Wait for shutdown to finish before closing the database pools.
	return group.Wait()
}

// newDogmaApp constructs the bank's [example.App] with its projection database.
func newDogmaApp(projectionDB *sql.DB) *example.App {
	return &example.App{ReadDB: projectionDB}
}

// newEngine creates a [runkit.Engine] to execute the application.
func newEngine(application *example.App, engineDB *sql.DB, logger *slog.Logger) *runkit.Engine {
	return runkit.New(
		application,
		engineDB,
		runkit.WithLogger(logger),
	)
}

// newHTTPHandler creates a [ui.Handler] that queries projections and submits commands.
func newHTTPHandler(projectionDB *sql.DB, executor dogma.CommandExecutor) *ui.Handler {
	return &ui.Handler{
		DB:              projectionDB,
		CommandExecutor: executor,
	}
}
