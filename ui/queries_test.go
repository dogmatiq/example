package ui

import (
	"testing"
	"time"

	"github.com/dogmatiq/example/internal/xtesting"
)

func Test_queries(t *testing.T) {
	db := xtesting.NewDatabase(t)
	ctx := t.Context()
	if _, err := db.ExecContext(ctx, `INSERT INTO customers (id, name) VALUES ('C001', 'Anna'), ('C002', 'Bob')`); err != nil {
		t.Fatal(err)
	}
	// Balances above the PostgreSQL INTEGER limit must round-trip as int64.
	const balance = 3_000_000_000
	if _, err := db.ExecContext(ctx,
		`INSERT INTO accounts (id, name, customer_id, balance) VALUES
		('A001', 'Savings', 'C001', $1),
		('A002', 'Checking', 'C001', 0),
		('A003', 'Savings', 'C002', 0)`, balance,
	); err != nil {
		t.Fatal(err)
	}
	occurredAt := time.Date(2026, time.November, 16, 10, 20, 30, 123456000, time.FixedZone("Brisbane", 10*60*60))
	if _, err := db.ExecContext(ctx,
		`INSERT INTO ledger (account_id, transaction_id, description, credit, balance, created_at)
		VALUES ('A001', 'T001', 'Deposit', $1, $1, $2)`, balance, occurredAt,
	); err != nil {
		t.Fatal(err)
	}

	h := &Handler{DB: db}
	name, err := h.queryCustomerName(ctx, "C001")
	if err != nil || name != "Anna" {
		t.Fatalf("customer name = %q, %v; want Anna", name, err)
	}
	name, gotBalance, err := h.queryAccountDetails(ctx, "A001")
	if err != nil || name != "Savings" || gotBalance != balance {
		t.Fatalf("account details = %q, %d, %v", name, gotBalance, err)
	}
	accounts, err := h.queryAccounts(ctx, "C001")
	if err != nil {
		t.Fatal(err)
	}
	if len(accounts) != 2 || accounts[0].ID != "A002" || accounts[1].Balance != balance {
		t.Fatalf("unexpected accounts: %+v", accounts)
	}
	groups, err := h.queryAllAccountsGrouped(ctx, "C001", "A001")
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 2 || !groups[0].IsOwn || groups[0].Accounts[0].ID != "A002" || groups[1].Accounts[0].ID != "A003" {
		t.Fatalf("unexpected account groups: %+v", groups)
	}
	transactions, err := h.queryTransactions(ctx, "A001")
	if err != nil {
		t.Fatal(err)
	}
	if len(transactions) != 1 {
		t.Fatalf("got %d transactions, want 1", len(transactions))
	}
	tx := transactions[0]
	if !tx.OccurredAt.Equal(occurredAt) || tx.Credit != balance || tx.Balance != balance || tx.Description != "Deposit" {
		t.Fatalf("unexpected transaction: %+v", tx)
	}
}
