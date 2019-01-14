package app

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/examples/cmd/bank/internal/messages"
)

// TransferProcessHandler manages the process of transferring funds between accounts.
var TransferProcessHandler dogma.ProcessMessageHandler = transferProcessHandler{}

type transferProcessHandler struct {
	dogma.NoTimeoutBehavior
}

func (transferProcessHandler) New() dogma.ProcessRoot {
	return &transfer{}
}

func (transferProcessHandler) Configure(c dogma.ProcessConfigurer) {
	c.Name("transfer")
	c.RouteEventType(messages.TransferStarted{})
	c.RouteEventType(messages.AccountDebitedForTransfer{})
	c.RouteEventType(messages.AccountCreditedForTransfer{})
	c.RouteEventType(messages.TransferDeclined{})
}

func (transferProcessHandler) RouteEventToInstance(_ context.Context, m dogma.Message) (string, bool, error) {
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

func (transferProcessHandler) HandleEvent(
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

// transfer is the process root for a funds transfer.
type transfer struct {
	ToAccountID string
}

func (t *transfer) IsEqual(r dogma.ProcessRoot) bool {
	v, ok := r.(*transfer)
	return ok && *t == *v
}

func (t transfer) Clone() dogma.ProcessRoot {
	return &t
}