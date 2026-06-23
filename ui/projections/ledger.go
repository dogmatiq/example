package projections

import (
	"context"
	"database/sql"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages"
	"github.com/dogmatiq/example/messages/events"
	"github.com/dogmatiq/projectionkit/sqlprojection"
)

// LedgerProjectionHandler maintains account balances and a ledger of
// transactions against each account.
//
// It handles both the accounts and ledger tables within a single projection to
// ensure the running balance recorded on each ledger entry is consistent with
// the account balance.
//
// The UI queries the accounts table to list a customer's accounts and their
// balances, and the ledger table to display transaction history.
type LedgerProjectionHandler struct {
	sqlprojection.NoCompactBehavior
}

// Configure configs the engine for this projection.
func (h *LedgerProjectionHandler) Configure(c dogma.ProjectionConfigurer) {
	c.Identity("ledger", "a1b2c3d4-e5f6-7890-abcd-ef1234567890")

	c.Routes(
		dogma.HandlesEvent[*events.AccountOpened](),
		dogma.HandlesEvent[*events.AccountCredited](),
		dogma.HandlesEvent[*events.AccountDebited](),
	)
}

// HandleEvent inserts into the "ledger" table whenever an account is credited
// or debited, and updates the "accounts" table to reflect the current balance.
func (h *LedgerProjectionHandler) HandleEvent(
	ctx context.Context,
	tx *sql.Tx,
	s dogma.ProjectionEventScope,
	m dogma.Event,
) error {
	switch x := m.(type) {
	case *events.AccountOpened:
		return h.accountOpened(ctx, tx, x)
	case *events.AccountCredited:
		return h.accountCredited(ctx, tx, s, x)
	case *events.AccountDebited:
		return h.accountDebited(ctx, tx, s, x)
	default:
		panic(dogma.UnexpectedMessage)
	}
}

func (h *LedgerProjectionHandler) accountOpened(
	ctx context.Context,
	tx *sql.Tx,
	x *events.AccountOpened,
) error {
	_, err := tx.ExecContext(
		ctx,
		`INSERT INTO accounts (
			id,
			name,
			customer_id
		) VALUES (
			?,
			?,
			?
		)`,
		x.AccountID,
		x.AccountName,
		x.CustomerID,
	)
	return err
}

func (h *LedgerProjectionHandler) accountCredited(
	ctx context.Context,
	tx *sql.Tx,
	s dogma.ProjectionEventScope,
	x *events.AccountCredited,
) error {
	var balance int64

	if err := tx.QueryRowContext(
		ctx,
		`UPDATE accounts SET
			balance = balance + ?
		WHERE id = ?
		RETURNING balance`,
		x.Amount,
		x.AccountID,
	).Scan(&balance); err != nil {
		return err
	}

	var isRefund bool

	if err := tx.QueryRowContext(
		ctx,
		`SELECT EXISTS (
			SELECT 1 FROM ledger
			WHERE account_id = ?
				AND transaction_id = ?
				AND debit > 0
		)`,
		x.AccountID,
		x.TransactionID,
	).Scan(&isRefund); err != nil {
		return err
	}

	_, err := tx.ExecContext(
		ctx,
		`INSERT INTO ledger (
			account_id,
			transaction_id,
			transaction_order,
			description,
			credit,
			balance,
			created_at
		) VALUES (
			?,
			?,
			(SELECT COUNT(*) FROM ledger WHERE transaction_id = ?),
			?,
			?,
			?,
			?
		)`,
		x.AccountID,
		x.TransactionID,
		x.TransactionID,
		creditDescription(x.TransactionType, isRefund),
		x.Amount,
		balance,
		s.RecordedAt(),
	)
	return err
}

func (h *LedgerProjectionHandler) accountDebited(
	ctx context.Context,
	tx *sql.Tx,
	s dogma.ProjectionEventScope,
	x *events.AccountDebited,
) error {
	var balance int64

	if err := tx.QueryRowContext(
		ctx,
		`UPDATE accounts SET
			balance = balance - ?
		WHERE id = ?
		RETURNING balance`,
		x.Amount,
		x.AccountID,
	).Scan(&balance); err != nil {
		return err
	}

	_, err := tx.ExecContext(
		ctx,
		`INSERT INTO ledger (
			account_id,
			transaction_id,
			transaction_order,
			description,
			debit,
			balance,
			created_at
		) VALUES (
			?,
			?,
			(SELECT COUNT(*) FROM ledger WHERE transaction_id = ?),
			?,
			?,
			?,
			?
		)`,
		x.AccountID,
		x.TransactionID,
		x.TransactionID,
		debitDescription(x.TransactionType),
		x.Amount,
		balance,
		s.RecordedAt(),
	)
	return err
}

// Reset clears all projection data.
func (h *LedgerProjectionHandler) Reset(
	ctx context.Context,
	tx *sql.Tx,
	_ dogma.ProjectionResetScope,
) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM ledger`); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM accounts`); err != nil {
		return err
	}

	return nil
}

// creditDescription returns the human-readable description for a credit entry
// in the ledger.
//
// If isRefund is true, the description indicates that the credit is a refund of
// a prior debit within the same transaction.
func creditDescription(t messages.TransactionType, isRefund bool) string {
	if isRefund {
		return debitDescription(t) + " refund"
	}

	switch t {
	case messages.Deposit:
		return "Deposit"
	case messages.Transfer:
		return "Incoming transfer"
	default:
		panic("unrecognized transaction type for credit: " + string(t))
	}
}

// debitDescription returns the human-readable description for a debit entry in
// the ledger.
func debitDescription(t messages.TransactionType) string {
	switch t {
	case messages.Withdrawal:
		return "Withdrawal"
	case messages.Transfer:
		return "Outgoing transfer"
	default:
		panic("unrecognized transaction type for debit: " + string(t))
	}
}
