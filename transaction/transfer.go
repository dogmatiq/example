package transaction

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages"
)

// transfer is the process root for a funds transfer.
type transfer struct {
	ToAccountID string
}

// TransferProcess manages the process of transferring funds between accounts.
type TransferProcess struct {
	dogma.NoTimeoutBehavior
}

// New returns a new transfer instance.
func (TransferProcess) New() dogma.ProcessRoot {
	return &transfer{}
}

// Configure configures the behavior of the engine as it relates to this
// handler.
func (TransferProcess) Configure(c dogma.ProcessConfigurer) {
	c.Name("transfer")
	c.RouteEventType(messages.TransferStarted{})
	c.RouteEventType(messages.AccountDebitedForTransfer{})
	c.RouteEventType(messages.AccountCreditedForTransfer{})
	c.RouteEventType(messages.TransferDeclined{})
}

// RouteEventToInstance returns the ID of the process instance that is targetted
// by m.
func (TransferProcess) RouteEventToInstance(_ context.Context, m dogma.Message) (string, bool, error) {
	switch x := m.(type) {
	case messages.TransferStarted:
		return x.TransactionID, true, nil
	case messages.AccountDebitedForTransfer:
		return x.TransactionID, true, nil
	case messages.AccountCreditedForTransfer:
		return x.TransactionID, true, nil
	case messages.TransferDeclined:
		return x.TransactionID, true, nil
	default:
		panic(dogma.UnexpectedMessage)
	}
}

// HandleEvent handles an event message that has been routed to this handler.
func (TransferProcess) HandleEvent(
	_ context.Context,
	s dogma.ProcessEventScope,
	m dogma.Message,
) error {
	switch x := m.(type) {
	case messages.TransferStarted:
		s.Begin()

		xfer := s.Root().(*transfer)
		xfer.ToAccountID = x.ToAccountID

		s.ExecuteCommand(messages.DebitAccountForTransfer{
			TransactionID: x.TransactionID,
			AccountID:     x.FromAccountID,
			Amount:        x.Amount,
		})

	case messages.AccountDebitedForTransfer:
		xfer := s.Root().(*transfer)

		s.ExecuteCommand(messages.CreditAccountForTransfer{
			TransactionID: x.TransactionID,
			AccountID:     xfer.ToAccountID,
			Amount:        x.Amount,
		})

	case messages.AccountCreditedForTransfer, messages.TransferDeclined:
		s.End()

	default:
		panic(dogma.UnexpectedMessage)
	}

	return nil
}
