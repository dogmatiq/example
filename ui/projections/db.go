package projections

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/dogmatiq/projectionkit/sqlprojection"
	_ "github.com/jackc/pgx/v5/stdlib" // install "pgx" driver
)

// schema contains the SQL schema files.
//
//go:embed *.sql
var schema embed.FS

// NewDB opens a PostgreSQL database and creates the schema needed by the UI
// projections. dsn is a PostgreSQL connection string.
//
// The caller must close the returned database.
func NewDB(ctx context.Context, dsn string) (*sql.DB, error) {
	if dsn == "" {
		return nil, errors.New("PostgreSQL connection string is required (set DATABASE_URL)")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("unable to connect to PostgreSQL: %w", err)
	}

	if err := CreateSchema(ctx, db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// CreateSchema creates the schema elements required by the projection adapter
// and handlers. It can be called again on an existing database.
func CreateSchema(ctx context.Context, db *sql.DB) error {
	if err := sqlprojection.PostgresDriver.CreateSchema(ctx, db); err != nil {
		return err
	}

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
