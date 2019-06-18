package customer

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
)

// customer is the aggregate root for a bank customer.
type customer struct {
	// Email is the customer email address.
	Email string
}

func (c *customer) ApplyEvent(m dogma.Message) {
	switch x := m.(type) {
	case events.CustomerAcquired:
		c.Email = x.CustomerEmail
	case events.CustomerEmailAddressChanged:
		c.Email = x.CustomerEmail
	}
}

// Aggregate implements the business logic for a bank customer.
type Aggregate struct{}

// New returns a new customer instance.
func (Aggregate) New() dogma.AggregateRoot {
	return &customer{}
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (Aggregate) Configure(c dogma.AggregateConfigurer) {
	c.Name("customer")

	c.ConsumesCommandType(commands.OpenAccountForNewCustomer{})
	c.ConsumesCommandType(commands.ChangeCustomerEmailAddress{})

	c.ProducesEventType(events.CustomerAcquired{})
	c.ProducesEventType(events.CustomerEmailAddressChanged{})
}

// RouteCommandToInstance returns the ID of the aggregate instance that is
// targetted by m.
func (Aggregate) RouteCommandToInstance(m dogma.Message) string {
	switch x := m.(type) {
	case commands.OpenAccountForNewCustomer:
		return x.CustomerID
	case commands.ChangeCustomerEmailAddress:
		return x.CustomerID
	default:
		panic(dogma.UnexpectedMessage)
	}
}

// HandleCommand handles a command message that has been routed to this handler.
func (Aggregate) HandleCommand(s dogma.AggregateCommandScope, m dogma.Message) {
	switch x := m.(type) {
	case commands.OpenAccountForNewCustomer:
		acquire(s, x)
	case commands.ChangeCustomerEmailAddress:
		changeEmailAddress(s, x)
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
		CustomerID:    m.CustomerID,
		CustomerName:  m.CustomerName,
		CustomerEmail: m.CustomerEmail,
		AccountID:     m.AccountID,
		AccountName:   m.AccountName,
	})
}

func changeEmailAddress(s dogma.AggregateCommandScope, m commands.ChangeCustomerEmailAddress) {
	r := s.Root().(*customer)

	if r.Email != m.CustomerEmail {
		s.RecordEvent(events.CustomerEmailAddressChanged{
			CustomerID:    m.CustomerID,
			CustomerEmail: m.CustomerEmail,
		})
	}
}
