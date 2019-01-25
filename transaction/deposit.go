package transaction

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages"
)

// DepositProcess is manages the process of depositing funds into an account.
type DepositProcess struct {
	dogma.StatelessProcessBehavior
	dogma.NoTimeoutBehavior
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (DepositProcess) Configure(c dogma.ProcessConfigurer) {
	c.Name("deposit")
	c.RouteEventType(messages.DepositStarted{})
	c.RouteEventType(messages.AccountCreditedForDeposit{})
}

// RouteEventToInstance returns the ID of the process instance that is targetted
// by m.
func (DepositProcess) RouteEventToInstance(_ context.Context, m dogma.Message) (string, bool, error) {
	switch x := m.(type) {
	case messages.DepositStarted:
		return x.TransactionID, true, nil
	case messages.AccountCreditedForDeposit:
		return x.TransactionID, true, nil
	default:
		panic(dogma.UnexpectedMessage)
	}
}

// HandleEvent handles an event message that has been routed to this handler.
func (DepositProcess) HandleEvent(_ context.Context, s dogma.ProcessEventScope, m dogma.Message) error {
	switch x := m.(type) {
	case messages.DepositStarted:
		s.Begin()
		s.ExecuteCommand(messages.CreditAccountForDeposit{
			TransactionID: x.TransactionID,
			AccountID:     x.AccountID,
			Amount:        x.Amount,
		})

	case messages.AccountCreditedForDeposit:
		s.End()

	default:
		panic(dogma.UnexpectedMessage)
	}

	return nil
}
