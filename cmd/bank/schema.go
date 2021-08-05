package main

import (
	"context"
	"database/sql"
	_ "embed"
)

//go:embed schema.sql
var schemaDDL string

// createSchema creates the schema elements required by the projection handlers.
func createSchema(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, schemaDDL)
	return err
}
