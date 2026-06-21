package projections

import (
	"context"
	"database/sql"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages"
	"github.com/dogmatiq/example/messages/events"
	"github.com/dogmatiq/projectionkit/sqlprojection"
)

// LedgerProjectionHandler is a projection that maintains account balances
// and ledger entries. It handles both tables in a single projection to
// ensure the balance recorded on each ledger entry is consistent.
type LedgerProjectionHandler struct {
	sqlprojection.NoCompactBehavior
}

// Configure configs the engine for this projection.
func (h *LedgerProjectionHandler) Configure(c dogma.ProjectionConfigurer) {
	c.Identity("account-ledger", "a1b2c3d4-e5f6-7890-abcd-ef1234567890")

	c.Routes(
		dogma.HandlesEvent[*events.AccountOpened](),
		dogma.HandlesEvent[*events.AccountCredited](),
		dogma.HandlesEvent[*events.AccountDebited](),
	)
}

// HandleEvent updates the accounts and ledger tables.
func (h *LedgerProjectionHandler) HandleEvent(
	ctx context.Context,
	tx *sql.Tx,
	s dogma.ProjectionEventScope,
	m dogma.Event,
) error {
	switch x := m.(type) {
	case *events.AccountOpened:
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

	case *events.AccountCredited:
		var balance int64
		err := tx.QueryRowContext(
			ctx,
			`UPDATE accounts SET
				balance = balance + ?
			WHERE id = ?
			RETURNING balance`,
			x.Amount,
			x.AccountID,
		).Scan(&balance)
		if err != nil {
			return err
		}

		desc := creditDescription(x.TransactionType)
		if x.TransactionType == messages.Transfer {
			var exists bool
			_ = tx.QueryRowContext(
				ctx,
				`SELECT 1
				FROM ledger
				WHERE account_id = ?
				AND transaction_id = ?
				AND debit > 0`,
				x.AccountID,
				x.TransactionID,
			).Scan(&exists)
			if exists {
				desc = "Transfer reversal"
			}
		}

		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO ledger (
				account_id,
				transaction_id,
				description,
				credit,
				balance,
				created_at
			) VALUES (
				?,
				?,
				?,
				?,
				?,
				?
			)`,
			x.AccountID,
			x.TransactionID,
			desc,
			x.Amount,
			balance,
			s.RecordedAt(),
		)
		return err

	case *events.AccountDebited:
		var balance int64
		err := tx.QueryRowContext(
			ctx,
			`UPDATE accounts SET
				balance = balance - ?
			WHERE id = ?
			RETURNING balance`,
			x.Amount,
			x.AccountID,
		).Scan(&balance)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO ledger (
				account_id,
				transaction_id,
				description,
				debit,
				balance,
				created_at
			) VALUES (
				?,
				?,
				?,
				?,
				?,
				?
			)`,
			x.AccountID,
			x.TransactionID,
			debitDescription(x.TransactionType),
			x.Amount,
			balance,
			s.RecordedAt(),
		)
		return err

	default:
		panic(dogma.UnexpectedMessage)
	}
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

func creditDescription(t messages.TransactionType) string {
	switch t {
	case messages.Deposit:
		return "Deposit"
	case messages.Transfer:
		return "Incoming transfer"
	default:
		return string(t)
	}
}

func debitDescription(t messages.TransactionType) string {
	switch t {
	case messages.Withdrawal:
		return "Withdrawal"
	case messages.Transfer:
		return "Outgoing transfer"
	default:
		return string(t)
	}
}
