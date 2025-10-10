package projections

import (
	"context"
	"database/sql"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages/events"
	"github.com/dogmatiq/projectionkit/sqlprojection"
)

// AccountProjectionHandler is a projection that builds a report of accounts
// managed by the bank.
type AccountProjectionHandler struct {
	sqlprojection.NoCompactBehavior
}

// Configure configs the engine for this projection.
func (h *AccountProjectionHandler) Configure(c dogma.ProjectionConfigurer) {
	c.Identity("account-list", "38dcb02a-3d76-4798-9c2a-186f8764ba19")

	c.Routes(
		dogma.HandlesEvent[*events.AccountOpened](),
		dogma.HandlesEvent[*events.AccountCredited](),
		dogma.HandlesEvent[*events.AccountDebited](),
	)
}

// HandleEvent updates the in-memory records to reflect the occurence of m.
func (h *AccountProjectionHandler) HandleEvent(
	ctx context.Context,
	tx *sql.Tx,
	_ dogma.ProjectionEventScope,
	m dogma.Event,
) error {
	switch x := m.(type) {
	case *events.AccountOpened:
		_, err := tx.ExecContext(
			ctx,
			`INSERT INTO account (
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
		_, err := tx.ExecContext(
			ctx,
			`UPDATE account SET balance = balance + ? WHERE id = ?`,
			x.Amount,
			x.AccountID,
		)
		return err

	case *events.AccountDebited:
		_, err := tx.ExecContext(
			ctx,
			`UPDATE account SET balance = balance - ? WHERE id = ?`,
			x.Amount,
			x.AccountID,
		)
		return err

	default:
		panic(dogma.UnexpectedMessage)
	}
}
