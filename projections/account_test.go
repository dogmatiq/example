package projections_test

import (
	"testing"

	"github.com/dogmatiq/example/internal/database"
	"github.com/dogmatiq/example/internal/testrunner"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/testkit"
	"github.com/dogmatiq/testkit/engine"
)

func Test_AccountProjectionHandler(t *testing.T) {
	t.Run(
		"when an account is opened for a new customer",
		func(t *testing.T) {
			db := database.New()
			defer db.Close()

			testrunner.New(db).
				Begin(
					t,
					testkit.WithOperationOptions(
						engine.EnableProjections(true),
					),
				).
				Prepare(
					commands.OpenAccountForNewCustomer{
						CustomerID:   "C001",
						CustomerName: "Anna Smith",
						AccountID:    "A001",
						AccountName:  "Savings",
					},
				)

			rows, err := db.Query(
				`SELECT
					id,
					name,
					customer_id
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
			)

			if err := rows.Scan(
				&id,
				&name,
				&customerID,
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

			if rows.Next() {
				t.Fatal("expected no more rows")
			}
		},
	)
}
