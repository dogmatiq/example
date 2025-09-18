package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dogmatiq/projectionkit/sqlprojection"
)

// New returns an in-memory SQLite database, with database tables necessary to
// run the example application.
//
// It returns an error if the database is unable to be opened, or the schema is
// unable to be created.
func New() (*sql.DB, error) {
	ctx := context.Background()

	const filename = "artifacts/bank.sqlite3"

	if err := os.MkdirAll(
		filepath.Dir(filename),
		0700,
	); err != nil {
		return nil, err
	}

	db, err := sql.Open(
		"sqlite3",
		fmt.Sprintf("file:%s?mode=rwc", filename),
	)
	if err != nil {
		return nil, err
	}

	if err := sqlprojection.SQLiteDriver.CreateSchema(ctx, db); err != nil {
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
