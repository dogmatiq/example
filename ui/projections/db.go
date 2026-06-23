package projections

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"sync/atomic"

	"github.com/dogmatiq/projectionkit/sqlprojection"
	_ "github.com/mattn/go-sqlite3" // install "sqlite3" driver
)

var (
	// counter is used to build a unique name for each in-memory database
	// instance.
	counter atomic.Int64

	// schema contains the SQL schema files.
	//
	//go:embed *.sql
	schema embed.FS
)

// NewDB returns an in-memory SQLite database, with database tables necessary to
// run the example application.
//
// It returns an error if the database is unable to be opened, or the schema is
// unable to be created.
func NewDB() (*sql.DB, error) {
	ctx := context.Background()

	db, err := sql.Open(
		"sqlite3",
		fmt.Sprintf("file:db%d?mode=memory&cache=shared", counter.Add(1)),
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

// MustNewDB returns an in-memory SQLite database, with database tables
// necessary to run the example application.
//
// It panics if the database is unable to be opened, or the schema is unable to
// be created.
func MustNewDB() *sql.DB {
	db, err := NewDB()
	if err != nil {
		panic(err)
	}

	return db
}

// CreateSchema creates the schema elements required by the projection handlers.
func CreateSchema(ctx context.Context, db *sql.DB) error {
	entries, err := schema.ReadDir(".")
	if err != nil {
		return err
	}

	for _, e := range entries {
		s, err := schema.ReadFile(e.Name())
		if err != nil {
			return err
		}

		if _, err := db.ExecContext(ctx, string(s)); err != nil {
			return err
		}
	}

	return nil
}
