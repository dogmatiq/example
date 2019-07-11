package domain

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
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

	c.ConsumesEventType(events.DepositStarted{})
	c.ConsumesEventType(events.AccountCreditedForDeposit{})

	c.ProducesCommandType(commands.CreditAccountForDeposit{})
}

// RouteEventToInstance returns the ID of the process instance that is targetted
// by m.
func (DepositProcess) RouteEventToInstance(_ context.Context, m dogma.Message) (string, bool, error) {
	switch x := m.(type) {
	case events.DepositStarted:
		return x.TransactionID, true, nil
	case events.AccountCreditedForDeposit:
		return x.TransactionID, true, nil
	default:
		panic(dogma.UnexpectedMessage)
	}
}

// HandleEvent handles an event message that has been routed to this handler.
func (DepositProcess) HandleEvent(_ context.Context, s dogma.ProcessEventScope, m dogma.Message) error {
	switch x := m.(type) {
	case events.DepositStarted:
		s.Begin()
		s.ExecuteCommand(commands.CreditAccountForDeposit{
			TransactionID: x.TransactionID,
			AccountID:     x.AccountID,
			Amount:        x.Amount,
		})

	case events.AccountCreditedForDeposit:
		s.End()

	default:
		panic(dogma.UnexpectedMessage)
	}

	return nil
}
