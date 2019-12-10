package database

import (
	"context"
	"database/sql"

	"github.com/dogmatiq/projectionkit/sql/sqlite"
)

// New returns an in-memory SQLite database, with database tables necessary to
// run the example application.
//
// It panics if the database is unable to be opened, or the schema is unable to be
// created.
func New() *sql.DB {
	ctx := context.Background()

	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		panic(err)
	}

	// Setup the pool to ensure the memory database survives when all
	// connections are returned to the pool.
	//
	// See https://github.com/mattn/go-sqlite3#faq.
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(-1)

	if err := sqlite.CreateSchema(ctx, db); err != nil {
		panic(err)
	}

	if _, err := db.ExecContext(
		ctx,
		`CREATE TABLE customer (
			id   TEXT NOT NULL,
			name TEXT NOT NULL,

			PRIMARY KEY (id)
		);

		CREATE TABLE account (
			id          TEXT NOT NULL,
			name        TEXT NOT NULL,
			customer_id TEXT NOT NULL,
			balance     INTEGER NOT NULL DEFAULT 0,

			PRIMARY KEY (id)
		);

		CREATE INDEX idx_account_customer ON account (customer_id);
		`,
	); err != nil {
		panic(err)
	}

	return db
}
