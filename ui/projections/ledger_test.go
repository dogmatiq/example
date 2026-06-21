package projections_test

import (
	"testing"
	"time"

	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/messages/commands"
	. "github.com/dogmatiq/testkit"
)

func Test_LedgerProjectionHandler(t *testing.T) {
	t.Run(
		"when an account is opened",
		func(t *testing.T) {
			db := openDB(t)

			Begin(t, &example.App{ReadDB: db}).
				EnableHandlers("account-ledger").
				Prepare(
					ExecuteCommand(
						&commands.OpenAccountForNewCustomer{
							CustomerID:   "C001",
							CustomerName: "Anna Smith",
							AccountID:    "A001",
							AccountName:  "Savings",
						},
					),
				)

			var (
				id         string
				name       string
				customerID string
				balance    int64
			)

			err := db.QueryRow(
				`SELECT id, name, customer_id, balance
				FROM accounts
				WHERE id = ?`,
				"A001",
			).Scan(&id, &name, &customerID, &balance)
			if err != nil {
				t.Fatal(err)
			}

			if id != "A001" {
				t.Fatalf(`expected account ID to be "A001", got %q`, id)
			}
			if name != "Savings" {
				t.Fatalf(`expected account name to be "Savings", got %q`, name)
			}
			if customerID != "C001" {
				t.Fatalf(`expected customer ID to be "C001", got %q`, customerID)
			}
			if balance != 0 {
				t.Fatalf(`expected balance to be 0, got %d`, balance)
			}
		},
	)

	t.Run(
		"when an account is credited",
		func(t *testing.T) {
			db := openDB(t)

			Begin(t, &example.App{ReadDB: db}).
				EnableHandlers("account-ledger").
				Prepare(
					ExecuteCommand(
						&commands.OpenAccountForNewCustomer{
							CustomerID:   "C001",
							CustomerName: "Anna Smith",
							AccountID:    "A001",
							AccountName:  "Savings",
						},
					),
					ExecuteCommand(
						&commands.Deposit{
							TransactionID: "T001",
							AccountID:     "A001",
							Amount:        150,
						},
					),
				)

			var balance int64
			err := db.QueryRow(
				`SELECT balance FROM accounts WHERE id = ?`,
				"A001",
			).Scan(&balance)
			if err != nil {
				t.Fatal(err)
			}
			if balance != 150 {
				t.Fatalf(`expected account balance to be 150, got %d`, balance)
			}

			var (
				description   string
				credit, debit int64
				txBalance     int64
			)
			err = db.QueryRow(
				`SELECT description, debit, credit, balance
				FROM ledger
				WHERE account_id = ?`,
				"A001",
			).Scan(&description, &debit, &credit, &txBalance)
			if err != nil {
				t.Fatal(err)
			}
			if description != "Deposit" {
				t.Fatalf(`expected description to be "Deposit", got %q`, description)
			}
			if credit != 150 {
				t.Fatalf(`expected credit to be 150, got %d`, credit)
			}
			if debit != 0 {
				t.Fatalf(`expected debit to be 0, got %d`, debit)
			}
			if txBalance != 150 {
				t.Fatalf(`expected transaction balance to be 150, got %d`, txBalance)
			}
		},
	)

	t.Run(
		"when an account is debited",
		func(t *testing.T) {
			db := openDB(t)

			Begin(t, &example.App{ReadDB: db}).
				EnableHandlers("account-ledger").
				Prepare(
					ExecuteCommand(
						&commands.OpenAccountForNewCustomer{
							CustomerID:   "C001",
							CustomerName: "Anna Smith",
							AccountID:    "A001",
							AccountName:  "Savings",
						},
					),
					ExecuteCommand(
						&commands.Deposit{
							TransactionID: "T001",
							AccountID:     "A001",
							Amount:        500,
						},
					),
					ExecuteCommand(
						&commands.Withdraw{
							TransactionID: "T002",
							AccountID:     "A001",
							Amount:        150,
							ScheduledTime: time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC),
						},
					),
				)

			var balance int64
			err := db.QueryRow(
				`SELECT balance FROM accounts WHERE id = ?`,
				"A001",
			).Scan(&balance)
			if err != nil {
				t.Fatal(err)
			}
			if balance != 350 {
				t.Fatalf(`expected account balance to be 350, got %d`, balance)
			}

			var (
				description   string
				credit, debit int64
				txBalance     int64
			)
			err = db.QueryRow(
				`SELECT description, debit, credit, balance
				FROM ledger
				WHERE account_id = ? AND description = ?`,
				"A001", "Withdrawal",
			).Scan(&description, &debit, &credit, &txBalance)
			if err != nil {
				t.Fatal(err)
			}
			if description != "Withdrawal" {
				t.Fatalf(`expected description to be "Withdrawal", got %q`, description)
			}
			if debit != 150 {
				t.Fatalf(`expected debit to be 150, got %d`, debit)
			}
			if credit != 0 {
				t.Fatalf(`expected credit to be 0, got %d`, credit)
			}
			if txBalance != 350 {
				t.Fatalf(`expected transaction balance to be 350, got %d`, txBalance)
			}
		},
	)
}
