package domain

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
)

// DepositProcessHandler manages the process of depositing funds into an
// account.
type DepositProcessHandler struct {
	dogma.StatelessProcessBehavior
	dogma.NoTimeoutMessagesBehavior
	dogma.NoTimeoutHintBehavior
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (DepositProcessHandler) Configure(c dogma.ProcessConfigurer) {
	c.Identity("deposit", "4e50b66b-b2c4-4522-bf15-39756186caee")

	c.Routes(
		dogma.HandlesEvent[events.DepositStarted](),
		dogma.HandlesEvent[events.AccountCredited](),
		dogma.HandlesEvent[events.DepositApproved](),
		dogma.ExecutesCommand[commands.CreditAccount](),
		dogma.ExecutesCommand[commands.ApproveDeposit](),
	)
}

// RouteEventToInstance returns the ID of the process instance that is targetted
// by m.
func (DepositProcessHandler) RouteEventToInstance(_ context.Context, m dogma.Message) (string, bool, error) {
	switch x := m.(type) {
	case events.DepositStarted:
		return x.TransactionID, true, nil
	case events.AccountCredited:
		return x.TransactionID, x.TransactionType == messages.Deposit, nil
	case events.DepositApproved:
		return x.TransactionID, true, nil
	default:
		panic(dogma.UnexpectedMessage)
	}
}

// HandleEvent handles an event message that has been routed to this handler.
func (DepositProcessHandler) HandleEvent(
	_ context.Context,
	r dogma.ProcessRoot,
	s dogma.ProcessEventScope,
	m dogma.Message,
) error {
	switch x := m.(type) {
	case events.DepositStarted:
		s.ExecuteCommand(commands.CreditAccount{
			TransactionID:   x.TransactionID,
			AccountID:       x.AccountID,
			TransactionType: messages.Deposit,
			Amount:          x.Amount,
		})

	case events.AccountCredited:
		s.ExecuteCommand(commands.ApproveDeposit{
			TransactionID: x.TransactionID,
			AccountID:     x.AccountID,
			Amount:        x.Amount,
		})

	case events.DepositApproved:
		s.End()

	default:
		panic(dogma.UnexpectedMessage)
	}

	return nil
}
