package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/ui"
	"github.com/dogmatiq/example/ui/projections"
	"github.com/dogmatiq/runkit"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	db, err := projections.NewDB(ctx, os.Getenv("DOGMA_EXAMPLE_PROJECTION_DSN"))
	if err != nil {
		return err
	}
	defer db.Close()

	dsn := os.Getenv("DOGMA_EXAMPLE_ENGINE_DSN")
	if dsn == "" {
		return errors.New("DOGMA_EXAMPLE_ENGINE_DSN is required")
	}
	engineDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer engineDB.Close()

	if err := runkit.CreateSchema(ctx, engineDB); err != nil {
		return err
	}

	e := runkit.New(&example.App{ReadDB: db}, engineDB)
	group, ctx := errgroup.WithContext(ctx)
	group.Go(func() error { return e.Run(ctx) })

	server := &http.Server{
		Addr: ":8080",
		Handler: &ui.Handler{
			DB:              db,
			CommandExecutor: e,
		},
	}

	group.Go(func() error {
		fmt.Println("Dogmatiq Bank is running at http://localhost:8080")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	group.Go(func() error {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	})

	if err := group.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}
