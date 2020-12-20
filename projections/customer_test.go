package projections_test

import (
	"context"
	"testing"

	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/messages/commands"
	. "github.com/dogmatiq/testkit"
	"github.com/dogmatiq/testkit/engine"
)

func Test_CustomerProjectionHandler(t *testing.T) {
	t.Run(
		"when an account is opened for a new customer",
		func(t *testing.T) {
			database, db := openDB(context.Background())
			defer database.Close()

			Begin(
				t,
				&example.App{ReadDB: db},
				WithUnsafeOperationOptions(
					engine.EnableProjections(true),
				),
			).Prepare(
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
					name
				FROM customer`,
			)
			if err != nil {
				t.Fatal(err)
			}
			defer rows.Close()

			if !rows.Next() {
				t.Fatal("expected a database row")
			}

			var (
				id   string
				name string
			)

			if err := rows.Scan(
				&id,
				&name,
			); err != nil {
				t.Fatal(err)
			}

			if id != "C001" {
				t.Fatalf(
					`expected customer ID to be "C001", got "%s"`,
					id,
				)
			}

			if name != "Anna Smith" {
				t.Fatalf(
					`expected customer name to be "Anna Smith", got "%s"`,
					name,
				)
			}

			if rows.Next() {
				t.Fatal("expected no more rows")
			}
		},
	)
}
