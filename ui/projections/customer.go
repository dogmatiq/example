package projections

import (
	"context"
	"database/sql"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages/events"
	"github.com/dogmatiq/projectionkit/sqlprojection"
)

// CustomerProjectionHandler maintains a list of the bank's customers.
//
// The UI queries the customers table to populate the login page, which simply
// lets the user pick an existing customer rather than performing any real
// authentication. It is also used to display the customer's name throughout the
// interface.
type CustomerProjectionHandler struct {
	sqlprojection.NoCompactBehavior
}

// Configure configs the engine for this projection.
func (h *CustomerProjectionHandler) Configure(c dogma.ProjectionConfigurer) {
	c.Identity("customers", "70b58269-931a-46d6-b745-286a670fb6f7")

	c.Routes(
		dogma.HandlesEvent[*events.CustomerAcquired](),
	)
}

// HandleEvent inserts into the "customers" table whenever the bank acquires a
// new customer.
func (h *CustomerProjectionHandler) HandleEvent(
	ctx context.Context,
	tx *sql.Tx,
	_ dogma.ProjectionEventScope,
	m dogma.Event,
) error {
	switch x := m.(type) {
	case *events.CustomerAcquired:
		_, err := tx.ExecContext(
			ctx,
			`INSERT INTO customers (
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

// Reset clears all projection data.
func (h *CustomerProjectionHandler) Reset(
	ctx context.Context,
	tx *sql.Tx,
	_ dogma.ProjectionResetScope,
) error {
	_, err := tx.ExecContext(
		ctx,
		`DELETE FROM customers`,
	)
	return err
}
