package database

import (
	"context"
	"database/sql"
	"os"

	"github.com/dogmatiq/projectionkit/sqlprojection"
)

// New returns an in-memory SQLite database, with database tables necessary to
// run the example application.
//
// It returns an error if the database is unable to be opened, or the schema is
// unable to be created.
func New() (*sql.DB, error) {
	ctx := context.Background()

	dsn := os.Getenv("DSN")
	if dsn == "" {
		// The default DSN is configured for use with a PostgreSQL server
		// running under docker as the docker stack configuratin in
		// https://github.com/dogmatiq/sqltest.
		dsn = "user=postgres password=rootpass sslmode=disable host=127.0.0.1 port=25432"
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := sqlprojection.CreateSchema(ctx, db); err != nil {
		db.Close()
		return nil, err
	}

	if err := CreateSchema(ctx, db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// MustNew returns an in-memory SQLite database, with database tables necessary
// to run the example application.
//
// It panics if the database is unable to be opened, or the schema is unable to
// be created.
func MustNew() *sql.DB {
	db, err := New()
	if err != nil {
		panic(err)
	}

	return db
}
