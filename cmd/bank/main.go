package main

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/api"
	"github.com/dogmatiq/projectionkit/sqlprojection"
	"github.com/dogmatiq/verity"
	"github.com/dogmatiq/verity/persistence/sqlpersistence"
	_ "github.com/jackc/pgx/v4/stdlib"
	"golang.org/x/sync/errgroup"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	ctx := context.Background()

	if err := run(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		// The default DSN is configured for use with a PostgreSQL server
		// running under docker as the docker stack configuratin in
		// https://github.com/dogmatiq/sqltest.
		dsn = "user=postgres password=rootpass sslmode=disable host=127.0.0.1 port=25432"
	}

	// Open a connection to the PostgreSQL database.
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	// Create the SQL schema elements required by dogmatiq/verity.
	if err := sqlpersistence.CreateSchema(ctx, db); err != nil {
		return err
	}

	// Create the SQL schema elements required by dogmatiq/projectionkit.
	if err := sqlprojection.CreateSchema(ctx, db); err != nil {
		return err
	}

	// Create the SQL schema elements required by the example bank application.
	if err := createSchema(ctx, db); err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(ctx)

	// Run the verity engine using the PostgreSQL database for persistence.
	g.Go(func() error {
		return verity.Run(
			ctx,
			&example.App{
				ReadDB: db,
			},
			verity.WithPersistence(
				&sqlpersistence.Provider{
					DB: db,
				},
			),
		)
	})

	// Run the JSON-RPC API server.
	g.Go(func() error {
		port := os.Getenv("API_PORT")
		if port == "" {
			port = "3001"
		}

		server := &http.Server{
			Addr:    net.JoinHostPort("", port),
			Handler: api.NewHandler(),
		}

		go func() {
			<-ctx.Done()
			server.Close()
		}()

		fmt.Printf("listening for API requests on %s\n", server.Addr)
		return server.ListenAndServe()
	})

	return g.Wait()
}
