package domain

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
)

// WithdrawalProcessHandler manages the process of withdrawing funds from an
// account.
type WithdrawalProcessHandler struct {
	dogma.StatelessProcessBehavior
	dogma.NoTimeoutMessagesBehavior
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (WithdrawalProcessHandler) Configure(c dogma.ProcessConfigurer) {
	c.Identity("withdrawal", "23f70d2b-a289-4e0f-8a83-c0c6a69d11d9")

	c.Routes(
		dogma.HandlesEvent[events.WithdrawalStarted](),
		dogma.HandlesEvent[events.AccountDebited](),
		dogma.HandlesEvent[events.AccountDebitDeclined](),
		dogma.HandlesEvent[events.DailyDebitLimitConsumed](),
		dogma.HandlesEvent[events.DailyDebitLimitExceeded](),
		dogma.HandlesEvent[events.AccountCredited](),
		dogma.HandlesEvent[events.WithdrawalApproved](),
		dogma.HandlesEvent[events.WithdrawalDeclined](),
		dogma.ExecutesCommand[commands.DebitAccount](),
		dogma.ExecutesCommand[commands.ConsumeDailyDebitLimit](),
		dogma.ExecutesCommand[commands.CreditAccount](),
		dogma.ExecutesCommand[commands.ApproveWithdrawal](),
		dogma.ExecutesCommand[commands.DeclineWithdrawal](),
	)
}

// RouteEventToInstance returns the ID of the process instance that is targetted
// by m.
func (WithdrawalProcessHandler) RouteEventToInstance(
	_ context.Context,
	m dogma.Event,
) (string, bool, error) {
	switch x := m.(type) {
	case events.WithdrawalStarted:
		return x.TransactionID, true, nil
	case events.AccountDebited:
		return x.TransactionID, x.TransactionType == messages.Withdrawal, nil
	case events.AccountDebitDeclined:
		return x.TransactionID, x.TransactionType == messages.Withdrawal, nil
	case events.DailyDebitLimitConsumed:
		return x.TransactionID, x.DebitType == messages.Withdrawal, nil
	case events.DailyDebitLimitExceeded:
		return x.TransactionID, x.DebitType == messages.Withdrawal, nil
	case events.AccountCredited:
		return x.TransactionID, x.TransactionType == messages.Withdrawal, nil
	case events.WithdrawalApproved:
		return x.TransactionID, true, nil
	case events.WithdrawalDeclined:
		return x.TransactionID, true, nil
	default:
		panic(dogma.UnexpectedMessage)
	}
}

// HandleEvent handles an event message that has been routed to this handler.
func (WithdrawalProcessHandler) HandleEvent(
	_ context.Context,
	_ dogma.ProcessRoot,
	s dogma.ProcessEventScope,
	m dogma.Event,
) error {
	switch x := m.(type) {
	case events.WithdrawalStarted:
		s.ExecuteCommand(commands.DebitAccount{
			TransactionID:   x.TransactionID,
			AccountID:       x.AccountID,
			TransactionType: messages.Withdrawal,
			Amount:          x.Amount,
			ScheduledTime:   x.ScheduledTime,
		})

	case events.AccountDebited:
		s.ExecuteCommand(commands.ConsumeDailyDebitLimit{
			TransactionID: x.TransactionID,
			AccountID:     x.AccountID,
			DebitType:     messages.Withdrawal,
			Amount:        x.Amount,
			Date:          messages.DailyDebitLimitDate(x.ScheduledTime),
		})

	case events.AccountDebitDeclined:
		s.ExecuteCommand(commands.DeclineWithdrawal{
			TransactionID: x.TransactionID,
			AccountID:     x.AccountID,
			Amount:        x.Amount,
			Reason:        x.Reason,
		})

	case events.DailyDebitLimitConsumed:
		s.ExecuteCommand(commands.ApproveWithdrawal{
			TransactionID: x.TransactionID,
			AccountID:     x.AccountID,
			Amount:        x.Amount,
		})

	case events.DailyDebitLimitExceeded:
		s.ExecuteCommand(commands.CreditAccount{
			TransactionID:   x.TransactionID,
			AccountID:       x.AccountID,
			TransactionType: messages.Withdrawal,
			Amount:          x.Amount,
		})

	case events.AccountCredited:
		s.ExecuteCommand(commands.DeclineWithdrawal{
			TransactionID: x.TransactionID,
			AccountID:     x.AccountID,
			Amount:        x.Amount,
			Reason:        messages.DailyDebitLimitExceeded,
		})

	case events.WithdrawalApproved,
		events.WithdrawalDeclined:
		s.End()

	default:
		panic(dogma.UnexpectedMessage)
	}

	return nil
}
