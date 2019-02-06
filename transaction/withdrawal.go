package transaction

import (
	"context"
	"time"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
)

// WithdrawalProcess manages the process of withdrawing funds from an account.
type WithdrawalProcess struct {
	dogma.StatelessProcessBehavior
	dogma.NoTimeoutBehavior
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (WithdrawalProcess) Configure(c dogma.ProcessConfigurer) {
	c.Name("withdrawal")
	c.RouteEventType(events.WithdrawalStarted{})
	c.RouteEventType(events.WithdrawalApprovedByDebitPolicy{})
	c.RouteEventType(events.AccountDebitedForWithdrawal{})
	c.RouteEventType(events.WithdrawalDeclinedDueToInsufficientFunds{})
	c.RouteEventType(events.WithdrawalDeclinedByDebitPolicy{})
}

// RouteEventToInstance returns the ID of the process instance that is targetted
// by m.
func (WithdrawalProcess) RouteEventToInstance(_ context.Context, m dogma.Message) (string, bool, error) {
	switch x := m.(type) {
	case events.WithdrawalStarted:
		return x.TransactionID, true, nil
	case events.WithdrawalApprovedByDebitPolicy:
		return x.TransactionID, true, nil
	case events.AccountDebitedForWithdrawal:
		return x.TransactionID, true, nil
	case events.WithdrawalDeclinedDueToInsufficientFunds:
		return x.TransactionID, true, nil
	case events.WithdrawalDeclinedByDebitPolicy:
		return x.TransactionID, true, nil
	default:
		return "", false, nil
	}
}

// HandleEvent handles an event message that has been routed to this handler.
func (WithdrawalProcess) HandleEvent(
	_ context.Context,
	s dogma.ProcessEventScope,
	m dogma.Message,
) error {
	switch x := m.(type) {
	case events.WithdrawalStarted:
		s.Begin()
		s.ExecuteCommand(commands.CheckWithdrawalAllowedByDebitPolicy{
			Timestamp:     time.Now(),
			TransactionID: x.TransactionID,
			AccountID:     x.AccountID,
			Amount:        x.Amount,
		})

	case events.WithdrawalApprovedByDebitPolicy:
		s.ExecuteCommand(commands.DebitAccountForWithdrawal{
			TransactionID: x.TransactionID,
			AccountID:     x.AccountID,
			Amount:        x.Amount,
		})

	case events.AccountDebitedForWithdrawal,
		events.WithdrawalDeclinedDueToInsufficientFunds,
		events.WithdrawalDeclinedByDebitPolicy:
		s.End()

	default:
		panic(dogma.UnexpectedMessage)
	}

	return nil
}
