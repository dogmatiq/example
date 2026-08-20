package projections

import (
	"context"
	"database/sql"
	"embed"
	"os"

	"github.com/dogmatiq/projectionkit/sqlprojection"
	_ "github.com/mattn/go-sqlite3" // install "sqlite3" driver
)

var (
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

	file, err := os.CreateTemp("", "bank-*.sqlite3")
	if err != nil {
		return nil, err
	}
	file.Close()

	db, err := sql.Open(
		"sqlite3",
		file.Name()+"?_journal_mode=WAL&_busy_timeout=5000",
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
