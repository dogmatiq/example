package projections

import (
	"context"
	"database/sql"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages/events"
)

// AccountProjectionHandler is a projection that builds a report of accounts
// managed by the bank.
type AccountProjectionHandler struct {
	dogma.NoTimeoutHintBehavior
}

// Configure configs the engine for this projection.
func (h *AccountProjectionHandler) Configure(c dogma.ProjectionConfigurer) {
	c.Identity("account-projection", "38dcb02a-3d76-4798-9c2a-186f8764ba19")

	c.ConsumesEventType(events.AccountOpened{})
}

// HandleEvent updates the in-memory records to reflect the occurence of m.
func (h *AccountProjectionHandler) HandleEvent(
	ctx context.Context,
	tx *sql.Tx,
	s dogma.ProjectionEventScope,
	m dogma.Message,
) error {
	switch x := m.(type) {
	case events.AccountOpened:
		_, err := tx.ExecContext(
			ctx,
			`INSERT INTO account (
				id,
				name,
				customer_id,
				balance
			) VALUES (
				?,
				?,
				?,
				0
			)`,
			x.AccountID,
			x.AccountName,
			x.CustomerID,
		)
		return err

	default:
		panic(dogma.UnexpectedMessage)
	}
}
