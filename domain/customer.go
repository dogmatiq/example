package domain

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
)

// customer is the aggregate root for a bank customer.
type customer struct {
	dogma.NoSnapshotBehavior

	Name string
}

func (c *customer) AggregateInstanceDescription() string {
	return c.Name
}

func (c *customer) Acquire(s dogma.AggregateCommandScope, m *commands.OpenAccountForNewCustomer) {
	if c.Name != "" {
		s.Log("customer has already been acquired")
		return
	}

	s.RecordEvent(&events.CustomerAcquired{
		CustomerID:   m.CustomerID,
		CustomerName: m.CustomerName,
		AccountID:    m.AccountID,
		AccountName:  m.AccountName,
	})
}

func (c *customer) ApplyEvent(m dogma.Event) {
	switch m := m.(type) {
	case *events.CustomerAcquired:
		c.Name = m.CustomerName
	}
}

// CustomerHandler implements the business logic for a bank customer.
type CustomerHandler struct{}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (CustomerHandler) Configure(c dogma.AggregateConfigurer) {
	c.Identity("customer", "f30111d5-f100-4495-90ad-b09746ba8477")

	c.Routes(
		dogma.HandlesCommand[*commands.OpenAccountForNewCustomer](),
		dogma.RecordsEvent[*events.CustomerAcquired](),
	)
}

// New returns a new customer instance.
func (CustomerHandler) New() dogma.AggregateRoot {
	return &customer{}
}

// RouteCommandToInstance returns the ID of the aggregate instance that is
// targetted by m.
func (CustomerHandler) RouteCommandToInstance(m dogma.Command) string {
	switch x := m.(type) {
	case *commands.OpenAccountForNewCustomer:
		return x.CustomerID
	default:
		panic(dogma.UnexpectedMessage)
	}
}

// HandleCommand handles a command message that has been routed to this handler.
func (CustomerHandler) HandleCommand(
	r dogma.AggregateRoot,
	s dogma.AggregateCommandScope,
	m dogma.Command,
) {
	c := r.(*customer)

	switch x := m.(type) {
	case *commands.OpenAccountForNewCustomer:
		c.Acquire(s, x)
	default:
		panic(dogma.UnexpectedMessage)
	}
}
