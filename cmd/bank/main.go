package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/dogmatiq/enginekit/config/runtimeconfig"
	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/database"
	"github.com/dogmatiq/example/ui"
	"github.com/dogmatiq/testkit/engine"
	"github.com/dogmatiq/testkit/fact"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	db := database.MustNew()
	defer db.Close()

	app := &example.App{
		ReadDB: db,
	}

	e, err := engine.New(runtimeconfig.FromApplication(app))
	if err != nil {
		panic(err)
	}

	logger := fact.NewLogger(func(s string) {
		fmt.Println(s)
	})

	observer := fact.ObserverFunc(func(f fact.Fact) {
		switch f.(type) {
		case fact.DispatchBegun,
			fact.HandlingBegun,
			fact.HandlingCompleted,
			fact.EventRecordedByAggregate,
			fact.EventRecordedByIntegration,
			fact.CommandExecutedByProcess,
			fact.TimeoutScheduledByProcess,
			fact.MessageLoggedByAggregate,
			fact.MessageLoggedByIntegration,
			fact.MessageLoggedByProcess,
			fact.MessageLoggedByProjection:
			logger.Notify(f)
		}
	})

	opts := []engine.OperationOption{
		engine.EnableProjections(true),
		engine.WithObserver(observer),
	}

	// Run the engine in the background. This processes timeouts and scheduled
	// events that are triggered by process managers.
	go func() {
		if err := engine.Run(ctx, e, 0, opts...); err != nil {
			fmt.Fprintln(os.Stderr, "engine error:", err)
		}
	}()

	server := &http.Server{
		Addr: ":8080",
		Handler: &ui.Handler{
			DB: db,
			CommandExecutor: engine.CommandExecutor{
				Engine:  e,
				Options: opts,
			},
		},
	}

	// Shut down the HTTP server when the context is canceled.
	go func() {
		<-ctx.Done()
		server.Shutdown(context.Background())
	}()

	fmt.Println("Dogmatiq Bank is running at http://localhost:8080")

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Fprintln(os.Stderr, "server error:", err)
		os.Exit(1)
	}
}
