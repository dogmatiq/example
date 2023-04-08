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
	dogma.NoTimeoutMessagesBehavior
	dogma.NoTimeoutHintBehavior
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (OpenAccountForNewCustomerProcessHandler) Configure(c dogma.ProcessConfigurer) {
	c.Identity("open-account-for-new-customer", "89b39176-a57f-4071-afad-e0db62137fd3")

	c.Routes(
		dogma.HandlesEvent[events.CustomerAcquired](),
		dogma.ExecutesCommand[commands.OpenAccount](),
	)
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
func (OpenAccountForNewCustomerProcessHandler) HandleEvent(
	_ context.Context,
	r dogma.ProcessRoot,
	s dogma.ProcessEventScope,
	m dogma.Message,
) error {
	switch x := m.(type) {
	case events.CustomerAcquired:
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
