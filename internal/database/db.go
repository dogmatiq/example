package database

import (
	"context"
	"database/sql"

	"github.com/dogmatiq/projectionkit/sql/sqlite"
)

// New returns an in memory database with pre-created schema.
func New() *sql.DB {
	ctx := context.Background()

	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		panic(err)
	}
	defer db.Close()

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
		)`,
	); err != nil {
		panic(err)
	}

	return db
}
