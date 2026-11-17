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

	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/internal/publicdemo"
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

	var err error

	// The public demo hosted at https://example.dogmatiq.io uses its own
	// startup path. run() below is the regular Dogma example.
	if idle := os.Getenv("DOGMA_EXAMPLE_IDLE_TIMEOUT"); idle != "" {
		err = publicdemo.Run(ctx, logger, idle)
	} else {
		err = run(ctx, logger)
	}

	if errors.Is(err, context.Canceled) {
		return
	}

	if err != nil {
		logger.Error("bank stopped", "error", err)

		os.Exit(1)
	}
}

// run starts the bank's engine and the HTTP server that hosts its UI.
// Both database pools close after the engine and server have stopped.
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
		return runHTTPServer(ctx, logger, uiHandler)
	})

	// Wait for shutdown before closing either database pool.
	return group.Wait()
}

// runHTTPServer serves the UI until ctx is canceled or the server fails.
// DOGMA_EXAMPLE_LISTEN_ADDRESS overrides the default listen address, ":8080".
func runHTTPServer(ctx context.Context, logger *slog.Logger, handler http.Handler) error {
	listenAddress := os.Getenv("DOGMA_EXAMPLE_LISTEN_ADDRESS")
	if listenAddress == "" {
		listenAddress = ":8080"
	}

	group, ctx := errgroup.WithContext(ctx)

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

	return group.Wait()
}
