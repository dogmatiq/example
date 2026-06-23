package projections_test

import (
	"testing"
	"time"

	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/ui/projections"
	. "github.com/dogmatiq/testkit"
)

func Test_LedgerProjectionHandler(t *testing.T) {
	t.Run(
		"when an account is opened",
		func(t *testing.T) {
			db := projections.MustNewDB()
			t.Cleanup(func() { db.Close() })

			Begin(t, &example.App{ReadDB: db}).
				EnableHandlers("ledger").
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
				name       string
				customerID string
				balance    int64
			)

			if err := db.QueryRow(
				`SELECT
					name,
					customer_id,
					balance
				FROM accounts
				WHERE id = "A001"`,
			).Scan(
				&name,
				&customerID,
				&balance,
			); err != nil {
				t.Fatal(err)
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
			db := projections.MustNewDB()
			t.Cleanup(func() { db.Close() })

			Begin(t, &example.App{ReadDB: db}).
				EnableHandlers("ledger").
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

			var accountBalance int64

			if err := db.QueryRow(
				`SELECT balance
				FROM accounts
				WHERE id = "A001"`,
			).Scan(&accountBalance); err != nil {
				t.Fatal(err)
			}

			if accountBalance != 150 {
				t.Fatalf(`expected account balance to be 150, got %d`, accountBalance)
			}

			var (
				description   string
				credit, debit int64
				ledgerBalance int64
			)

			if err := db.QueryRow(
				`SELECT
					description,
					debit,
					credit,
					balance
				FROM ledger
				WHERE account_id = "A001"`,
			).Scan(
				&description,
				&debit,
				&credit,
				&ledgerBalance,
			); err != nil {
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
			if ledgerBalance != 150 {
				t.Fatalf(`expected transaction balance to be 150, got %d`, ledgerBalance)
			}
		},
	)

	t.Run(
		"when an account is debited",
		func(t *testing.T) {
			db := projections.MustNewDB()
			t.Cleanup(func() { db.Close() })

			Begin(t, &example.App{ReadDB: db}).
				EnableHandlers("ledger").
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

			var accountBalance int64

			if err := db.QueryRow(
				`SELECT balance
				FROM accounts
				WHERE id = "A001"`,
			).Scan(&accountBalance); err != nil {
				t.Fatal(err)
			}

			if accountBalance != 350 {
				t.Fatalf(`expected account balance to be 350, got %d`, accountBalance)
			}

			var (
				description   string
				credit, debit int64
				ledgerBalance int64
			)

			if err := db.QueryRow(
				`SELECT
					description,
					debit,
					credit,
					balance
				FROM ledger
				WHERE account_id = "A001"
					AND transaction_id = "T002"`,
			).Scan(
				&description,
				&debit,
				&credit,
				&ledgerBalance,
			); err != nil {
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
			if ledgerBalance != 350 {
				t.Fatalf(`expected transaction balance to be 350, got %d`, ledgerBalance)
			}
		},
	)

	t.Run(
		"when a transfer is approved",
		func(t *testing.T) {
			db := projections.MustNewDB()
			t.Cleanup(func() { db.Close() })

			Begin(t, &example.App{ReadDB: db}).
				EnableHandlers("ledger").
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
						&commands.OpenAccountForNewCustomer{
							CustomerID:   "C002",
							CustomerName: "Bob Jones",
							AccountID:    "A002",
							AccountName:  "Checking",
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
						&commands.Transfer{
							TransactionID: "T002",
							FromAccountID: "A001",
							ToAccountID:   "A002",
							Amount:        200,
							ScheduledTime: time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC),
						},
					),
				)

			var description string

			if err := db.QueryRow(
				`SELECT description
				FROM ledger
				WHERE account_id = "A001"
					AND transaction_id = "T002"
					AND debit > 0`,
			).Scan(&description); err != nil {
				t.Fatal(err)
			}

			if description != "Outgoing transfer" {
				t.Fatalf(`expected description to be "Outgoing transfer", got %q`, description)
			}

			if err := db.QueryRow(
				`SELECT description
				FROM ledger
				WHERE account_id = "A002"
					AND transaction_id = "T002"
					AND credit > 0`,
			).Scan(&description); err != nil {
				t.Fatal(err)
			}

			if description != "Incoming transfer" {
				t.Fatalf(`expected description to be "Incoming transfer", got %q`, description)
			}
		},
	)

	t.Run(
		"when a transfer is declined after debiting",
		func(t *testing.T) {
			db := projections.MustNewDB()
			t.Cleanup(func() { db.Close() })

			Begin(t, &example.App{ReadDB: db}).
				EnableHandlers("ledger").
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
						&commands.OpenAccountForNewCustomer{
							CustomerID:   "C002",
							CustomerName: "Bob Jones",
							AccountID:    "A002",
							AccountName:  "Checking",
						},
					),
					ExecuteCommand(
						&commands.Deposit{
							TransactionID: "T001",
							AccountID:     "A001",
							Amount:        1000000,
						},
					),
					ExecuteCommand(
						&commands.Transfer{
							TransactionID: "T002",
							FromAccountID: "A001",
							ToAccountID:   "A002",
							Amount:        900001,
							ScheduledTime: time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC),
						},
					),
				)

			var description string

			if err := db.QueryRow(
				`SELECT description
				FROM ledger
				WHERE account_id = "A001"
					AND transaction_id = "T002"
					AND credit > 0`,
			).Scan(&description); err != nil {
				t.Fatal(err)
			}

			if description != "Outgoing transfer refund" {
				t.Fatalf(`expected description to be "Outgoing transfer refund", got %q`, description)
			}
		},
	)
}
