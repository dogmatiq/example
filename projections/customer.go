package projections

import (
	"context"
	"database/sql"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages/events"
)

// CustomerProjectionHandler is a projection that builds a report of customers
// acquired by the bank.
type CustomerProjectionHandler struct {
	dogma.NoTimeoutHintBehavior
}

// Configure configs the engine for this projection.
func (h *CustomerProjectionHandler) Configure(c dogma.ProjectionConfigurer) {
	c.Identity("customer-list", "70b58269-931a-46d6-b745-286a670fb6f7")

	c.ConsumesEventType(events.CustomerAcquired{})
}

// HandleEvent updates the in-memory records to reflect the occurence of m.
func (h *CustomerProjectionHandler) HandleEvent(
	ctx context.Context,
	tx *sql.Tx,
	s dogma.ProjectionEventScope,
	m dogma.Message,
) error {
	switch x := m.(type) {
	case events.CustomerAcquired:
		_, err := tx.ExecContext(
			ctx,
			`INSERT INTO customer (
				id,
				name
			) VALUES (
				?,
				?
			)`,
			x.CustomerID,
			x.CustomerName,
		)
		return err

	default:
		panic(dogma.UnexpectedMessage)
	}
}
