package withdrawal

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/message/command"
	"github.com/dogmatiq/example/message/event"
)

// ProcessHandler manages the process of withdrawing funds from an account.
type ProcessHandler struct {
	dogma.StatelessProcessBehavior
	dogma.NoTimeoutBehavior
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (ProcessHandler) Configure(c dogma.ProcessConfigurer) {
	c.Name("withdrawal")

	c.ConsumesEventType(event.WithdrawalStarted{})
	c.ConsumesEventType(event.AccountDebitedForWithdrawal{})
	c.ConsumesEventType(event.WithdrawalDeclinedDueToInsufficientFunds{})

	c.ProducesCommandType(command.DebitAccountForWithdrawal{})
}

// RouteEventToInstance returns the ID of the process instance that is targetted
// by m.
func (ProcessHandler) RouteEventToInstance(_ context.Context, m dogma.Message) (string, bool, error) {
	switch x := m.(type) {
	case event.WithdrawalStarted:
		return x.TransactionID, true, nil
	case event.AccountDebitedForWithdrawal:
		return x.TransactionID, true, nil
	case event.WithdrawalDeclinedDueToInsufficientFunds:
		return x.TransactionID, true, nil
	default:
		return "", false, nil
	}
}

// HandleEvent handles an event message that has been routed to this handler.
func (ProcessHandler) HandleEvent(
	_ context.Context,
	s dogma.ProcessEventScope,
	m dogma.Message,
) error {
	switch x := m.(type) {
	case event.WithdrawalStarted:
		s.Begin()
		s.ExecuteCommand(command.DebitAccountForWithdrawal{
			TransactionID: x.TransactionID,
			AccountID:     x.AccountID,
			Amount:        x.Amount,
		})

	case event.AccountDebitedForWithdrawal,
		event.WithdrawalDeclinedDueToInsufficientFunds:
		s.End()

	default:
		panic(dogma.UnexpectedMessage)
	}

	return nil
}
