package database

import (
	"context"
	"database/sql"
)

// CreateSchema creates the schema elements required by the projection handlers.
func CreateSchema(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS customer (
			id   TEXT NOT NULL,
			name TEXT NOT NULL,

			PRIMARY KEY (id)
		);

		CREATE TABLE IF NOT EXISTS account (
			id          TEXT NOT NULL,
			name        TEXT NOT NULL,
			customer_id TEXT NOT NULL,
			balance     INTEGER NOT NULL DEFAULT 0,

			PRIMARY KEY (id)
		);

		CREATE INDEX IF NOT EXISTS idx_account_customer ON account (customer_id);
		`,
	)
	return err
}

// DropSchema drops the schema elements required by the projection handlers.
func DropSchema(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(
		ctx,
		`DROP TABLE IF EXISTS customer;
		DROP TABLE IF EXISTS account;`,
	)
	return err
}
