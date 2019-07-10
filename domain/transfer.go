package domain

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
)

// transfer is the process root for a funds transfer.
type transferProcess struct {
	ToAccountID string
}

// TransferProcessHandler manages the process of transferring funds between
// accounts.
type TransferProcessHandler struct {
	dogma.NoTimeoutBehavior
}

// New returns a new transfer instance.
func (TransferProcessHandler) New() dogma.ProcessRoot {
	return &transferProcess{}
}

// Configure configures the behavior of the engine as it relates to this handler.
func (TransferProcessHandler) Configure(c dogma.ProcessConfigurer) {
	c.Name("transfer-process")

	c.ConsumesEventType(events.TransferStarted{})
	c.ConsumesEventType(events.AccountDebitedForTransfer{})
	c.ConsumesEventType(events.TransferDeclinedDueToInsufficientFunds{})
	c.ConsumesEventType(events.AccountCreditedForTransfer{})

	c.ProducesCommandType(commands.DebitAccountForTransfer{})
	c.ProducesCommandType(commands.CreditAccountForTransfer{})
}

// RouteEventToInstance returns the ID of the process instance that is targetted
// by m.
func (TransferProcessHandler) RouteEventToInstance(_ context.Context, m dogma.Message) (string, bool, error) {
	switch x := m.(type) {
	case events.TransferStarted:
		return x.TransactionID, true, nil
	case events.AccountDebitedForTransfer:
		return x.TransactionID, true, nil
	case events.TransferDeclinedDueToInsufficientFunds:
		return x.TransactionID, true, nil
	case events.AccountCreditedForTransfer:
		return x.TransactionID, true, nil
	default:
		panic(dogma.UnexpectedMessage)
	}
}

// HandleEvent handles an event message that has been routed to this handler.
func (TransferProcessHandler) HandleEvent(
	_ context.Context,
	s dogma.ProcessEventScope,
	m dogma.Message,
) error {
	switch x := m.(type) {
	case events.TransferStarted:
		s.Begin()

		r := s.Root().(*transferProcess)
		r.ToAccountID = x.ToAccountID

		s.ExecuteCommand(commands.DebitAccountForTransfer{
			TransactionID: x.TransactionID,
			AccountID:     x.FromAccountID,
			Amount:        x.Amount,
		})

	case events.AccountDebitedForTransfer:
		r := s.Root().(*transferProcess)

		s.ExecuteCommand(commands.CreditAccountForTransfer{
			TransactionID: x.TransactionID,
			AccountID:     r.ToAccountID,
			Amount:        x.Amount,
		})

	case events.AccountCreditedForTransfer,
		events.TransferDeclinedDueToInsufficientFunds:
		s.End()

	default:
		panic(dogma.UnexpectedMessage)
	}

	return nil
}
