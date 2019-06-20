package transfer

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/message/command"
	"github.com/dogmatiq/example/message/event"
)

// ProcessHandler manages the process of transferring funds between accounts.
type ProcessHandler struct {
	dogma.NoTimeoutBehavior
}

// New returns a new transfer instance.
func (ProcessHandler) New() dogma.ProcessRoot {
	return &root{}
}

// Configure configures the behavior of the engine as it relates to this handler.
func (ProcessHandler) Configure(c dogma.ProcessConfigurer) {
	c.Name("transfer")

	c.ConsumesEventType(event.TransferStarted{})
	c.ConsumesEventType(event.AccountDebitedForTransfer{})
	c.ConsumesEventType(event.AccountCreditedForTransfer{})
	c.ConsumesEventType(event.TransferDeclinedDueToInsufficientFunds{})

	c.ProducesCommandType(command.CreditAccountForTransfer{})
	c.ProducesCommandType(command.DebitAccountForTransfer{})
}

// RouteEventToInstance returns the ID of the process instance that is targetted
// by m.
func (ProcessHandler) RouteEventToInstance(_ context.Context, m dogma.Message) (string, bool, error) {
	switch x := m.(type) {
	case event.TransferStarted:
		return x.TransactionID, true, nil
	case event.AccountDebitedForTransfer:
		return x.TransactionID, true, nil
	case event.AccountCreditedForTransfer:
		return x.TransactionID, true, nil
	case event.TransferDeclinedDueToInsufficientFunds:
		return x.TransactionID, true, nil
	default:
		panic(dogma.UnexpectedMessage)
	}
}

// HandleEvent handles an event message that has been routed to this handler.
func (ProcessHandler) HandleEvent(
	_ context.Context,
	s dogma.ProcessEventScope,
	m dogma.Message,
) error {
	switch x := m.(type) {
	case event.TransferStarted:
		s.Begin()

		r := s.Root().(*root)
		r.ToAccountID = x.ToAccountID

		s.ExecuteCommand(command.DebitAccountForTransfer{
			TransactionID: x.TransactionID,
			AccountID:     x.FromAccountID,
			Amount:        x.Amount,
		})

	case event.AccountDebitedForTransfer:
		r := s.Root().(*root)

		s.ExecuteCommand(command.CreditAccountForTransfer{
			TransactionID: x.TransactionID,
			AccountID:     r.ToAccountID,
			Amount:        x.Amount,
		})

	case event.AccountCreditedForTransfer,
		event.TransferDeclinedDueToInsufficientFunds:
		s.End()

	default:
		panic(dogma.UnexpectedMessage)
	}

	return nil
}
