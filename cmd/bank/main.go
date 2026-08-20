package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/ui"
	"github.com/dogmatiq/example/ui/projections"
	"github.com/dogmatiq/runkit"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lmittmann/tint"
)

func main() {
	logger := slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level: slog.LevelDebug,
		}),
	)

	if err := run(logger); err != nil {
		logger.Error(
			"application error",
			slog.String("error", err.Error()),
		)

		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	engineDB, err := sql.Open("pgx", os.Getenv("RUNKIT_DSN"))
	if err != nil {
		return err
	}
	defer engineDB.Close()

	if err := runkit.CreateSchema(ctx, engineDB); err != nil {
		return err
	}

	readDB := projections.MustNewDB()
	defer readDB.Close()

	app := &example.App{
		ReadDB: readDB,
	}

	engine := &runkit.Engine{
		DB:     engineDB,
		App:    app,
		Logger: logger,
	}

	server := &http.Server{
		Addr: ":8080",
		Handler: http.TimeoutHandler(
			&ui.Handler{
				DB:              readDB,
				CommandExecutor: engine,
			},
			10*time.Second,
			"request timed out",
		),
	}

	logger.InfoContext(
		ctx,
		"Dogmatiq Bank is running",
		slog.String("url", "http://localhost:8080"),
	)

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		context.AfterFunc(ctx, func() {
			server.Shutdown(context.Background())
		})
		return server.ListenAndServe()
	})

	group.Go(func() error {

		return engine.Run(ctx)
	})

	return group.Wait()
}
