package publicdemo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/internal/idletimeout"
	"github.com/dogmatiq/example/ui"
	"github.com/dogmatiq/example/ui/projections"
	"github.com/dogmatiq/runkit"
	"golang.org/x/sync/errgroup"
)

// Run starts the hosted public demo with the given idle timeout. The caller
// owns signal handling, logging configuration, and reporting the returned error.
func Run(ctx context.Context, logger *slog.Logger, idle string) error {
	idleTimeout, err := time.ParseDuration(idle)
	if err != nil || idleTimeout <= 0 {
		return fmt.Errorf("DOGMA_EXAMPLE_IDLE_TIMEOUT must be a positive duration, got %q", idle)
	}

	// Each wake recreates the engine, database pools, and UI handler. The idle
	// middleware makes requests wait until the new handler is ready.
	var uiHandler http.Handler

	startApplication := func(ctx context.Context, ready func()) error {
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

		uiHandler = &ui.Handler{
			DB:              projectionDB,
			CommandExecutor: engine,
		}

		ready()

		return engine.Run(ctx)
	}

	serveUI := func(ctx context.Context, wrap func(http.Handler) http.Handler) error {
		handler := wrap(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			uiHandler.ServeHTTP(w, request)
		}))

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

		// Drain HTTP requests when the process is stopping.
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

	return idletimeout.Run(
		ctx,
		logger,
		idleTimeout,
		startApplication,
		serveUI,
	)
}
