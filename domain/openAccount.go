package domain

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
)

// OpenAccountForNewCustomerProcessHandler manages the process of opening the
// initial account for a new customer.
type OpenAccountForNewCustomerProcessHandler struct {
	dogma.StatelessProcessBehavior
	dogma.NoTimeoutBehavior
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (OpenAccountForNewCustomerProcessHandler) Configure(c dogma.ProcessConfigurer) {
	c.Name("open-account-for-new-customer")

	c.ConsumesEventType(events.CustomerAcquired{})

	c.ProducesCommandType(commands.OpenAccount{})
}

// RouteEventToInstance returns the ID of the process instance that is targetted
// by m.
func (OpenAccountForNewCustomerProcessHandler) RouteEventToInstance(_ context.Context, m dogma.Message) (string, bool, error) {
	switch x := m.(type) {
	case events.CustomerAcquired:
		return x.CustomerID, true, nil
	default:
		panic(dogma.UnexpectedMessage)
	}
}

// HandleEvent handles an event message that has been routed to this handler.
func (OpenAccountForNewCustomerProcessHandler) HandleEvent(_ context.Context, s dogma.ProcessEventScope, m dogma.Message) error {
	switch x := m.(type) {
	case events.CustomerAcquired:
		s.Begin()
		s.ExecuteCommand(commands.OpenAccount{
			CustomerID:  x.CustomerID,
			AccountID:   x.AccountID,
			AccountName: x.AccountName,
		})
		s.End()

	default:
		panic(dogma.UnexpectedMessage)
	}

	return nil
}
