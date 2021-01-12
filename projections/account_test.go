package projections_test

import (
	"context"
	"testing"
	"time"

	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/messages/commands"
	. "github.com/dogmatiq/testkit"
)

func Test_AccountProjectionHandler(t *testing.T) {
	t.Run(
		"when an account is opened",
		func(t *testing.T) {
			database, db := openDB(context.Background())
			defer database.Close()

			Begin(t, &example.App{ReadDB: db}).
				EnableHandlers("account-list").
				Prepare(
					ExecuteCommand(
						commands.OpenAccountForNewCustomer{
							CustomerID:   "C001",
							CustomerName: "Anna Smith",
							AccountID:    "A001",
							AccountName:  "Savings",
						},
					),
				)

			rows, err := db.Query(
				`SELECT
					id,
					name,
					customer_id,
					balance
				FROM account`,
			)
			if err != nil {
				t.Fatal(err)
			}
			defer rows.Close()

			if !rows.Next() {
				t.Fatal("expected a database row")
			}

			var (
				id         string
				name       string
				customerID string
				balance    int64
			)

			if err := rows.Scan(
				&id,
				&name,
				&customerID,
				&balance,
			); err != nil {
				t.Fatal(err)
			}

			if id != "A001" {
				t.Fatalf(
					`expected account ID to be "A001", got "%s"`,
					id,
				)
			}

			if name != "Savings" {
				t.Fatalf(
					`expected account name to be "Savings", got "%s"`,
					name,
				)
			}

			if customerID != "C001" {
				t.Fatalf(
					`expected customer ID to be "C001", got "%s"`,
					customerID,
				)
			}

			if balance != 0 {
				t.Fatalf(
					`expected balance to be 0, got "%d"`,
					balance,
				)
			}

			if rows.Next() {
				t.Fatal("expected no more rows")
			}
		},
	)

	t.Run(
		"when an account is credited",
		func(t *testing.T) {
			database, db := openDB(context.Background())
			defer database.Close()

			Begin(t, &example.App{ReadDB: db}).
				EnableHandlers("account-list").
				Prepare(
					ExecuteCommand(
						commands.OpenAccountForNewCustomer{
							CustomerID:   "C001",
							CustomerName: "Anna Smith",
							AccountID:    "A001",
							AccountName:  "Savings",
						},
					),
					ExecuteCommand(
						commands.Deposit{
							TransactionID: "T001",
							AccountID:     "A001",
							Amount:        150,
						},
					),
				)

			rows, err := db.Query(
				`SELECT
					id,
					balance
				FROM account
				WHERE id = "A001"`,
			)
			if err != nil {
				t.Fatal(err)
			}
			defer rows.Close()

			if !rows.Next() {
				t.Fatal("expected a database row")
			}

			var (
				id      string
				balance int64
			)

			if err := rows.Scan(
				&id,
				&balance,
			); err != nil {
				t.Fatal(err)
			}

			if balance != 150 {
				t.Fatalf(
					`expected balance to be 150, got "%d"`,
					balance,
				)
			}

			if rows.Next() {
				t.Fatal("expected no more rows")
			}
		},
	)

	t.Run(
		"when an account is debited",
		func(t *testing.T) {
			database, db := openDB(context.Background())
			defer database.Close()

			Begin(t, &example.App{ReadDB: db}).
				EnableHandlers("account-list").
				Prepare(
					ExecuteCommand(
						commands.OpenAccountForNewCustomer{
							CustomerID:   "C001",
							CustomerName: "Anna Smith",
							AccountID:    "A001",
							AccountName:  "Savings",
						},
					),
					ExecuteCommand(
						commands.Deposit{
							TransactionID: "T001",
							AccountID:     "A001",
							Amount:        500,
						},
					),
					ExecuteCommand(
						commands.Withdraw{
							TransactionID: "T002",
							AccountID:     "A001",
							Amount:        150,
							ScheduledDate: time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC),
						},
					),
				)

			rows, err := db.Query(
				`SELECT
					id,
					balance
				FROM account
				WHERE id = "A001"`,
			)
			if err != nil {
				t.Fatal(err)
			}
			defer rows.Close()

			if !rows.Next() {
				t.Fatal("expected a database row")
			}

			var (
				id      string
				balance int64
			)

			if err := rows.Scan(
				&id,
				&balance,
			); err != nil {
				t.Fatal(err)
			}

			if balance != 350 {
				t.Fatalf(
					`expected balance to be 350, got "%d"`,
					balance,
				)
			}

			if rows.Next() {
				t.Fatal("expected no more rows")
			}
		},
	)
}
