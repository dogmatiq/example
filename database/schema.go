package database

import (
	"context"
	"database/sql"
)

// CreateSchema creates the schema elements required by the projection handlers.
func CreateSchema(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS customers (
			id   TEXT NOT NULL,
			name TEXT NOT NULL,

			PRIMARY KEY (id)
		);

		CREATE TABLE IF NOT EXISTS accounts (
			id          TEXT NOT NULL,
			name        TEXT NOT NULL,
			customer_id TEXT NOT NULL,
			balance     INTEGER NOT NULL DEFAULT 0,

			PRIMARY KEY (id)
		);

		CREATE INDEX IF NOT EXISTS idx_accounts_customer ON accounts (customer_id);

		CREATE TABLE IF NOT EXISTS ledger (
			entry_id    INTEGER   PRIMARY KEY,
			account_id  TEXT      NOT NULL,
			description TEXT      NOT NULL,
			debit       INTEGER   NOT NULL DEFAULT 0,
			credit      INTEGER   NOT NULL DEFAULT 0,
			balance     INTEGER   NOT NULL,
			created_at  TIMESTAMP NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_ledger_account ON ledger (account_id);
		`,
	)
	return err
}

// DropSchema drops the schema elements required by the projection handlers.
func DropSchema(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(
		ctx,
		`DROP TABLE IF EXISTS ledger;
		DROP TABLE IF EXISTS customers;
		DROP TABLE IF EXISTS accounts;`,
	)
	return err
}
