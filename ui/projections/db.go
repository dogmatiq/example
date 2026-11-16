package projections

import (
	"context"
	"database/sql"
	"embed"
	"errors"

	"github.com/dogmatiq/projectionkit/sqlprojection"
	_ "github.com/jackc/pgx/v5/stdlib" // install "pgx" driver
)

var (
	// schema contains the SQL schema files.
	//
	//go:embed *.sql
	schema embed.FS
)

// NewDB returns a PostgreSQL database, with database tables necessary to
// run the example application.
//
// It returns an error if the database is unable to be opened, or the schema is
// unable to be created.
func NewDB(ctx context.Context, dsn string) (*sql.DB, error) {
	if dsn == "" {
		return nil, errors.New("DOGMA_EXAMPLE_PROJECTION_DSN is required")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := sqlprojection.PostgresDriver.CreateSchema(ctx, db); err != nil {
		db.Close()
		return nil, err
	}

	if err := CreateSchema(ctx, db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// MustNewDB returns a PostgreSQL database, with database tables
// necessary to run the example application.
//
// It panics if the database is unable to be opened, or the schema is unable to
// be created.
func MustNewDB(ctx context.Context, dsn string) *sql.DB {
	db, err := NewDB(ctx, dsn)
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
