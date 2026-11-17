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

	// Each wake creates a new engine and database pool, so it also needs a new UI
	// handler. Requests wait until that handler is ready before using it.
	var uiHandler http.Handler

	startApplication := func(ctx context.Context, ready func()) error {
		return runDogmaApp(
			ctx,
			logger,
			func(engine *runkit.Engine, projectionDB *sql.DB) {
				// Rebuild the UI when its engine and database pool are recreated.
				uiHandler = &ui.Handler{
					DB:              projectionDB,
					CommandExecutor: engine,
				}

				ready()
			},
		)
	}

	serveUI := func(ctx context.Context, wrap func(http.Handler) http.Handler) error {
		handler := http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			uiHandler.ServeHTTP(w, request)
		})

		return runHTTPServer(ctx, logger, wrap(handler))
	}

	// See runDogmaApp() for the Dogma example and runHTTPServer() for the
	// HTTP server that hosts the UI. idletimeout.Run() only manages idle
	// shutdown for the public demo hosted at https://example.dogmatiq.io;
	// you can skip its internals.
	err := idletimeout.Run(
		ctx,
		logger,
		startApplication,
		serveUI,
	)

	if errors.Is(err, context.Canceled) {
		return
	}

	if err != nil {
		logger.Error("bank stopped", "error", err)

		os.Exit(1)
	}
}

// runDogmaApp opens the bank's databases and runs its engine until ctx is
// canceled or an error occurs, then closes both database pools. ready receives
// the initialized engine and projection pool before the engine starts.
func runDogmaApp(
	ctx context.Context,
	logger *slog.Logger,
	ready func(*runkit.Engine, *sql.DB),
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

	ready(engine, projectionDB)

	// Deliver messages until the application is stopped.
	return engine.Run(ctx)
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
		// Bound idle connections so TCP keep-alives do not prevent sleeping.
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
