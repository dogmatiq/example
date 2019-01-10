package app

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/examples/cmd/bank/internal/messages"
)

// TransactionHandler implements the domain logic for a transaction of any
// kind against an account.
//
// It's sole purpose is to ensure the global uniqueness of transaction IDs.
var TransactionHandler dogma.AggregateMessageHandler = transactionHandler{}

type transactionHandler struct {
	dogma.StatelessAggregateBehavior
}

func (transactionHandler) Configure(c dogma.AggregateConfigurer) {
	c.Name("transaction")
	c.RouteCommandType(messages.Deposit{})
	c.RouteCommandType(messages.Withdraw{})
	c.RouteCommandType(messages.Transfer{})
}

func (transactionHandler) RouteCommandToInstance(m dogma.Message) string {
	switch x := m.(type) {
	case messages.Deposit:
		return x.TransactionID
	case messages.Withdraw:
		return x.TransactionID
	case messages.Transfer:
		return x.TransactionID
	default:
		panic(dogma.UnexpectedMessage)
	}
}

func (transactionHandler) HandleCommand(s dogma.AggregateCommandScope, m dogma.Message) {
	if !s.Create() {
		s.Log("transaction already exists")
		return
	}

	switch x := m.(type) {
	case messages.Deposit:
		s.RecordEvent(messages.DepositStarted{
			TransactionID: x.TransactionID,
			AccountID:     x.AccountID,
			Amount:        x.Amount,
		})

	case messages.Withdraw:
		s.RecordEvent(messages.WithdrawalStarted{
			TransactionID: x.TransactionID,
			AccountID:     x.AccountID,
			Amount:        x.Amount,
		})

	case messages.Transfer:
		s.RecordEvent(messages.TransferStarted{
			TransactionID: x.TransactionID,
			FromAccountID: x.FromAccountID,
			ToAccountID:   x.ToAccountID,
			Amount:        x.Amount,
		})

	default:
		panic(dogma.UnexpectedMessage)
	}
}
