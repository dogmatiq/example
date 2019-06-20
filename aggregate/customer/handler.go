package customer

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/message/command"
	"github.com/dogmatiq/example/message/event"
)

// AggregateHandler implements the business logic for a bank customer.
type AggregateHandler struct{}

// New returns a new customer instance.
func (AggregateHandler) New() dogma.AggregateRoot {
	return &root{}
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (AggregateHandler) Configure(c dogma.AggregateConfigurer) {
	c.Name("customer")

	c.ConsumesCommandType(command.OpenAccountForNewCustomer{})
	c.ConsumesCommandType(command.ChangeCustomerEmailAddress{})

	c.ProducesEventType(event.CustomerAcquired{})
	c.ProducesEventType(event.CustomerEmailAddressChanged{})
}

// RouteCommandToInstance returns the ID of the aggregate instance that is
// targetted by m.
func (AggregateHandler) RouteCommandToInstance(m dogma.Message) string {
	switch x := m.(type) {
	case command.OpenAccountForNewCustomer:
		return x.CustomerID
	case command.ChangeCustomerEmailAddress:
		return x.CustomerID
	default:
		panic(dogma.UnexpectedMessage)
	}
}

// HandleCommand handles a command message that has been routed to this handler.
func (AggregateHandler) HandleCommand(s dogma.AggregateCommandScope, m dogma.Message) {
	switch x := m.(type) {
	case command.OpenAccountForNewCustomer:
		acquire(s, x)
	case command.ChangeCustomerEmailAddress:
		changeEmailAddress(s, x)
	default:
		panic(dogma.UnexpectedMessage)
	}
}

func acquire(s dogma.AggregateCommandScope, m command.OpenAccountForNewCustomer) {
	if !s.Create() {
		s.Log("customer has already been acquired")
		return
	}

	s.RecordEvent(event.CustomerAcquired{
		CustomerID:    m.CustomerID,
		CustomerName:  m.CustomerName,
		CustomerEmail: m.CustomerEmail,
		AccountID:     m.AccountID,
		AccountName:   m.AccountName,
	})
}

func changeEmailAddress(s dogma.AggregateCommandScope, m command.ChangeCustomerEmailAddress) {
	r := s.Root().(*root)

	if r.Email != m.CustomerEmail {
		s.RecordEvent(event.CustomerEmailAddressChanged{
			CustomerID:    m.CustomerID,
			CustomerEmail: m.CustomerEmail,
		})
	}
}
