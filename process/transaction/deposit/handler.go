package deposit

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/message/command"
	"github.com/dogmatiq/example/message/event"
)

// ProcessHandler is manages the process of depositing funds into an account.
type ProcessHandler struct {
	dogma.StatelessProcessBehavior
	dogma.NoTimeoutBehavior
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (ProcessHandler) Configure(c dogma.ProcessConfigurer) {
	c.Name("deposit")

	c.ConsumesEventType(event.DepositStarted{})
	c.ConsumesEventType(event.AccountCreditedForDeposit{})

	c.ProducesCommandType(command.CreditAccountForDeposit{})
}

// RouteEventToInstance returns the ID of the process instance that is targetted
// by m.
func (ProcessHandler) RouteEventToInstance(_ context.Context, m dogma.Message) (string, bool, error) {
	switch x := m.(type) {
	case event.DepositStarted:
		return x.TransactionID, true, nil
	case event.AccountCreditedForDeposit:
		return x.TransactionID, true, nil
	default:
		panic(dogma.UnexpectedMessage)
	}
}

// HandleEvent handles an event message that has been routed to this handler.
func (ProcessHandler) HandleEvent(_ context.Context, s dogma.ProcessEventScope, m dogma.Message) error {
	switch x := m.(type) {
	case event.DepositStarted:
		s.Begin()
		s.ExecuteCommand(command.CreditAccountForDeposit{
			TransactionID: x.TransactionID,
			AccountID:     x.AccountID,
			Amount:        x.Amount,
		})

	case event.AccountCreditedForDeposit:
		s.End()

	default:
		panic(dogma.UnexpectedMessage)
	}

	return nil
}
