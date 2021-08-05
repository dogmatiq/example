package database

import (
	"context"
	"database/sql"
)

// CreateSchema creates the schema elements required by the projection handlers.
func CreateSchema(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(
		ctx,
		`CREATE SCHEMA IF NOT EXISTS bank`,
	)

	return err
}

// DropSchema drops the schema elements required by the projection handlers.
func DropSchema(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(
		ctx,
		`DROP TABLE IF EXISTS bank CASCADE`,
	)

	return err
}
