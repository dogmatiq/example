package openaccount

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/message/command"
	"github.com/dogmatiq/example/message/event"
)

// ProcessHandler manages the process of creating accounts for acquired
// customers.
type ProcessHandler struct {
	dogma.StatelessProcessBehavior
	dogma.NoTimeoutBehavior
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (ProcessHandler) Configure(c dogma.ProcessConfigurer) {
	c.Name("new-customer.open-account")

	c.ConsumesEventType(event.CustomerAcquired{})

	c.ProducesCommandType(command.OpenAccount{})
}

// RouteEventToInstance returns the ID of the process instance that is targetted
// by m.
func (ProcessHandler) RouteEventToInstance(_ context.Context, m dogma.Message) (string, bool, error) {
	switch x := m.(type) {
	case event.CustomerAcquired:
		return x.CustomerID, true, nil
	default:
		panic(dogma.UnexpectedMessage)
	}
}

// HandleEvent handles an event message that has been routed to this handler.
func (ProcessHandler) HandleEvent(_ context.Context, s dogma.ProcessEventScope, m dogma.Message) error {
	switch x := m.(type) {
	case event.CustomerAcquired:
		s.Begin()
		s.ExecuteCommand(command.OpenAccount{
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
