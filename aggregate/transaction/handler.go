package transaction

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/message/command"
	"github.com/dogmatiq/example/message/event"
)

// AggregateHandler implements the business logic for a transaction of any
// kind against an account.
//
// It's sole purpose is to ensure the global uniqueness of transaction IDs.
type AggregateHandler struct {
	dogma.StatelessAggregateBehavior
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (AggregateHandler) Configure(c dogma.AggregateConfigurer) {
	c.Name("transaction")

	c.ConsumesCommandType(command.Deposit{})
	c.ConsumesCommandType(command.Withdraw{})
	c.ConsumesCommandType(command.Transfer{})

	c.ProducesEventType(event.DepositStarted{})
	c.ProducesEventType(event.WithdrawalStarted{})
	c.ProducesEventType(event.TransferStarted{})
}

// RouteCommandToInstance returns the ID of the aggregate instance that is
// targetted by m.
func (AggregateHandler) RouteCommandToInstance(m dogma.Message) string {
	switch x := m.(type) {
	case command.Deposit:
		return x.TransactionID
	case command.Withdraw:
		return x.TransactionID
	case command.Transfer:
		return x.TransactionID
	default:
		panic(dogma.UnexpectedMessage)
	}
}

// HandleCommand handles a command message that has been routed to this
// handler.
func (AggregateHandler) HandleCommand(s dogma.AggregateCommandScope, m dogma.Message) {
	if !s.Create() {
		s.Log("transaction already exists")
		return
	}

	switch x := m.(type) {
	case command.Deposit:
		s.RecordEvent(event.DepositStarted{
			TransactionID: x.TransactionID,
			AccountID:     x.AccountID,
			Amount:        x.Amount,
		})

	case command.Withdraw:
		s.RecordEvent(event.WithdrawalStarted{
			TransactionID: x.TransactionID,
			AccountID:     x.AccountID,
			Amount:        x.Amount,
		})

	case command.Transfer:
		s.RecordEvent(event.TransferStarted{
			TransactionID: x.TransactionID,
			FromAccountID: x.FromAccountID,
			ToAccountID:   x.ToAccountID,
			Amount:        x.Amount,
		})

	default:
		panic(dogma.UnexpectedMessage)
	}
}
