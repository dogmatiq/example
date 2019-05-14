package transaction

import (
	"context"

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

	c.ConsumesEventType(events.WithdrawalStarted{})
	c.ConsumesEventType(events.AccountDebitedForWithdrawal{})
	c.ConsumesEventType(events.WithdrawalDeclinedDueToInsufficientFunds{})
	c.ConsumesEventType(events.DailyDebitAmountConsumed{})
	c.ConsumesEventType(events.DailyDebitAmountConsumtionRejected{})

	c.ProducesCommandType(commands.DebitAccountForWithdrawal{})
	c.ProducesCommandType(commands.ConsumeDailyDebitAmount{})
	c.ProducesCommandType(commands.RestoreDailyDebitAmount{})
	c.ProducesCommandType(commands.MarkWithdrawalDeclinedDueToDailyDebitLimit{})
}

// RouteEventToInstance returns the ID of the process instance that is targetted
// by m.
func (WithdrawalProcess) RouteEventToInstance(_ context.Context, m dogma.Message) (string, bool, error) {
	switch x := m.(type) {
	case events.WithdrawalStarted:
		return x.TransactionID, true, nil
	case events.AccountDebitedForWithdrawal:
		return x.TransactionID, true, nil
	case events.WithdrawalDeclinedDueToInsufficientFunds:
		return x.TransactionID, true, nil
	case events.DailyDebitAmountConsumed:
		return x.TransactionID, true, nil
	case events.DailyDebitAmountConsumtionRejected:
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
		s.Log("***NOT THIS WITHDRAWAL HANDLER! events.WithdrawalStarted***")
		s.Begin()
		s.ExecuteCommand(commands.ConsumeDailyDebitAmount{
			TransactionID:        x.TransactionID,
			AccountID:            x.AccountID,
			Amount:               x.Amount,
			TransactionTimestamp: x.TransactionTimestamp,
		})

	case events.DailyDebitAmountConsumtionRejected:
		s.Log("***NOT THIS WITHDRAWAL HANDLER! events.DailyDebitAmountConsumtionRejected***")
		s.ExecuteCommand(commands.MarkWithdrawalDeclinedDueToDailyDebitLimit{
			TransactionID: x.TransactionID,
			AccountID:     x.AccountID,
			Amount:        x.Amount,
		})
		s.End()

	case events.DailyDebitAmountConsumed:
		s.Log("***NOT THIS WITHDRAWAL HANDLER! events.DailyDebitAmountConsumed***")
		// TODO(KM): Disabled this for debugging the Transfer test.
		// s.ExecuteCommand(commands.DebitAccountForWithdrawal{
		// 	TransactionID:        x.TransactionID,
		// 	AccountID:            x.AccountID,
		// 	Amount:               x.Amount,
		// 	TransactionTimestamp: x.TransactionTimestamp,
		// })

	case events.WithdrawalDeclinedDueToInsufficientFunds:
		s.ExecuteCommand(commands.RestoreDailyDebitAmount{
			TransactionID:        x.TransactionID,
			AccountID:            x.AccountID,
			Amount:               x.Amount,
			TransactionTimestamp: x.TransactionTimestamp,
		})
		s.End()

	case events.AccountDebitedForWithdrawal:
		s.End()

	default:
		panic(dogma.UnexpectedMessage)
	}

	return nil
}
