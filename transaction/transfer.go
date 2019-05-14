package transaction

import (
	"context"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
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

// Configure configures the behavior of the engine as it relates to this handler.
func (TransferProcess) Configure(c dogma.ProcessConfigurer) {
	c.Name("transfer")

	c.ConsumesEventType(events.TransferStarted{})
	c.ConsumesEventType(events.AccountDebitedForTransfer{})
	c.ConsumesEventType(events.AccountCreditedForTransfer{})
	c.ConsumesEventType(events.TransferDeclinedDueToInsufficientFunds{})
	c.ConsumesEventType(events.DailyDebitAmountConsumed{})
	c.ConsumesEventType(events.DailyDebitAmountConsumtionRejected{})

	c.ProducesCommandType(commands.CreditAccountForTransfer{})
	c.ProducesCommandType(commands.DebitAccountForTransfer{})
	c.ProducesCommandType(commands.MarkTransferDeclinedDueToDailyDebitLimit{})
	c.ProducesCommandType(commands.RestoreDailyDebitAmount{})
}

// RouteEventToInstance returns the ID of the process instance that is targetted
// by m.
func (TransferProcess) RouteEventToInstance(_ context.Context, m dogma.Message) (string, bool, error) {
	switch x := m.(type) {
	case events.TransferStarted:
		return x.TransactionID, true, nil
	case events.AccountDebitedForTransfer:
		return x.TransactionID, true, nil
	case events.AccountCreditedForTransfer:
		return x.TransactionID, true, nil
	case events.TransferDeclinedDueToInsufficientFunds:
		return x.TransactionID, true, nil
	case events.DailyDebitAmountConsumed:
		return x.TransactionID, true, nil
	case events.DailyDebitAmountConsumtionRejected:
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
	case events.TransferStarted:
		s.Begin()

		xfer := s.Root().(*transfer)
		xfer.ToAccountID = x.ToAccountID

		s.ExecuteCommand(commands.ConsumeDailyDebitAmount{
			TransactionID:        x.TransactionID,
			AccountID:            x.FromAccountID,
			Amount:               x.Amount,
			TransactionTimestamp: x.TransactionTimestamp,
		})

	case events.DailyDebitAmountConsumtionRejected:
		s.ExecuteCommand(commands.MarkTransferDeclinedDueToDailyDebitLimit{
			TransactionID: x.TransactionID,
			AccountID:     x.AccountID,
			Amount:        x.Amount,
		})
		s.End()

	case events.DailyDebitAmountConsumed:
		s.ExecuteCommand(commands.DebitAccountForTransfer{
			TransactionID:        x.TransactionID,
			AccountID:            x.AccountID,
			Amount:               x.Amount,
			TransactionTimestamp: x.TransactionTimestamp,
		})

	case events.AccountDebitedForTransfer:
		xfer := s.Root().(*transfer)

		s.ExecuteCommand(commands.CreditAccountForTransfer{
			TransactionID: x.TransactionID,
			AccountID:     xfer.ToAccountID,
			Amount:        x.Amount,
		})

	case events.TransferDeclinedDueToInsufficientFunds:
		s.ExecuteCommand(commands.RestoreDailyDebitAmount{
			TransactionID:        x.TransactionID,
			AccountID:            x.AccountID,
			Amount:               x.Amount,
			TransactionTimestamp: x.TransactionTimestamp,
		})
		s.End()

	case events.AccountCreditedForTransfer:
		s.End()

	default:
		panic(dogma.UnexpectedMessage)
	}

	return nil
}
