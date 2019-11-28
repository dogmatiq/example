package domain

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
)

// CustomerHandler implements the business logic for a bank customer.
type CustomerHandler struct {
	dogma.StatelessAggregateBehavior
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (CustomerHandler) Configure(c dogma.AggregateConfigurer) {
	c.Identity("customer", "f30111d5-f100-4495-90ad-b09746ba8477")

	c.ConsumesCommandType(commands.OpenAccountForNewCustomer{})

	c.ProducesEventType(events.CustomerAcquired{})
}

// RouteCommandToInstance returns the ID of the aggregate instance that is
// targetted by m.
func (CustomerHandler) RouteCommandToInstance(m dogma.Message) string {
	switch x := m.(type) {
	case commands.OpenAccountForNewCustomer:
		return x.CustomerID
	default:
		panic(dogma.UnexpectedMessage)
	}
}

// HandleCommand handles a command message that has been routed to this handler.
func (CustomerHandler) HandleCommand(s dogma.AggregateCommandScope, m dogma.Message) {
	switch x := m.(type) {
	case commands.OpenAccountForNewCustomer:
		acquire(s, x)
	default:
		panic(dogma.UnexpectedMessage)
	}
}

func acquire(s dogma.AggregateCommandScope, m commands.OpenAccountForNewCustomer) {
	if !s.Create() {
		s.Log("customer has already been acquired")
		return
	}

	s.RecordEvent(events.CustomerAcquired{
		CustomerID:   m.CustomerID,
		CustomerName: m.CustomerName,
		AccountID:    m.AccountID,
		AccountName:  m.AccountName,
	})
}
