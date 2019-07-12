package domain

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
)

// DepositProcessHandler is manages the process of depositing funds into an
// account.
type DepositProcessHandler struct {
	dogma.StatelessProcessBehavior
	dogma.NoTimeoutBehavior
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (DepositProcessHandler) Configure(c dogma.ProcessConfigurer) {
	c.Name("deposit")

	c.ConsumesEventType(events.DepositStarted{})
	c.ConsumesEventType(events.AccountCredited{})
	c.ConsumesEventType(events.DepositApproved{})

	c.ProducesCommandType(commands.CreditAccount{})
	c.ProducesCommandType(commands.ApproveDeposit{})
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
func (DepositProcessHandler) HandleEvent(_ context.Context, s dogma.ProcessEventScope, m dogma.Message) error {
	switch x := m.(type) {
	case events.DepositStarted:
		s.Begin()
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
