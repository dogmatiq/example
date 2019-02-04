package transaction

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
)

// Aggregate implements the business logic for a transaction of any
// kind against an account.
//
// It's sole purpose is to ensure the global uniqueness of transaction IDs.
type Aggregate struct {
	dogma.StatelessAggregateBehavior
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (Aggregate) Configure(c dogma.AggregateConfigurer) {
	c.Name("transaction")
	c.RouteCommandType(commands.Deposit{})
	c.RouteCommandType(commands.Withdraw{})
	c.RouteCommandType(commands.Transfer{})
}

// RouteCommandToInstance returns the ID of the aggregate instance that is
// targetted by m.
func (Aggregate) RouteCommandToInstance(m dogma.Message) string {
	switch x := m.(type) {
	case commands.Deposit:
		return x.TransactionID
	case commands.Withdraw:
		return x.TransactionID
	case commands.Transfer:
		return x.TransactionID
	default:
		panic(dogma.UnexpectedMessage)
	}
}

// HandleCommand handles a command message that has been routed to this
// handler.
func (Aggregate) HandleCommand(s dogma.AggregateCommandScope, m dogma.Message) {
	if !s.Create() {
		s.Log("transaction already exists")
		return
	}

	switch x := m.(type) {
	case commands.Deposit:
		s.RecordEvent(events.DepositStarted{
			TransactionID: x.TransactionID,
			AccountID:     x.AccountID,
			Amount:        x.Amount,
		})

	case commands.Withdraw:
		s.RecordEvent(events.WithdrawalStarted{
			TransactionID: x.TransactionID,
			AccountID:     x.AccountID,
			Amount:        x.Amount,
		})

	case commands.Transfer:
		s.RecordEvent(events.TransferStarted{
			TransactionID: x.TransactionID,
			FromAccountID: x.FromAccountID,
			ToAccountID:   x.ToAccountID,
			Amount:        x.Amount,
		})

	default:
		panic(dogma.UnexpectedMessage)
	}
}
