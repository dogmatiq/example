package main

import (
	"context"
	"database/sql"
)

// createSchema creates the schema elements required by the projection handlers.
func createSchema(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(
		ctx,
		`CREATE SCHEMA IF NOT EXISTS bank`,
	)

	return err
}