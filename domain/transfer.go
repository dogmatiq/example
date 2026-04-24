package domain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
)

func init() {
	dogma.RegisterTimeout[*TransferReadyToProceed]("b5474f89-e985-4661-b85b-645347a4a645")
}

// transfer is the process root for a funds transfer.
type transferProcess struct {
	FromAccountID    string
	ToAccountID      string
	ToThirdPartyBank bool
	Amount           int64
	DeclineReason    messages.DebitFailureReason
}

// ProcessInstanceDescription returns a human-readable description of the
// transfer's current state.
func (p *transferProcess) ProcessInstanceDescription(ended bool) string {
	if p.Amount == 0 {
		return ""
	}

	if !ended {
		return fmt.Sprintf(
			"transferring %s from %s to %s",
			messages.FormatAmount(p.Amount),
			p.FromAccountID,
			p.ToAccountID,
		)
	}

	if p.DeclineReason != "" {
		return fmt.Sprintf(
			"transfer of %s from %s to %s declined: %s",
			messages.FormatAmount(p.Amount),
			p.FromAccountID,
			p.ToAccountID,
			p.DeclineReason,
		)
	}

	return fmt.Sprintf(
		"transferred %s from %s to %s",
		messages.FormatAmount(p.Amount),
		p.FromAccountID,
		p.ToAccountID,
	)
}

// MarshalBinary returns the transferProcess encoded as binary data.
func (p *transferProcess) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

// UnmarshalBinary decodes binary data into the transferProcess.
func (p *transferProcess) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}

// TransferProcessHandler manages the process of transferring funds between
// accounts.
type TransferProcessHandler struct{}

// New returns a new transfer instance.
func (TransferProcessHandler) New() dogma.ProcessRoot {
	return &transferProcess{}
}

// Configure configures the behavior of the engine as it relates to this handler.
func (TransferProcessHandler) Configure(c dogma.ProcessConfigurer) {
	c.Identity("transfer", "35afbe82-24c1-4868-a689-c2ec96c2e953")

	c.Routes(
		dogma.HandlesEvent[*events.TransferStarted](),
		dogma.HandlesEvent[*events.AccountDebited](),
		dogma.HandlesEvent[*events.AccountDebitDeclined](),
		dogma.HandlesEvent[*events.DailyDebitLimitConsumed](),
		dogma.HandlesEvent[*events.DailyDebitLimitExceeded](),
		dogma.HandlesEvent[*events.AccountCredited](),
		dogma.HandlesEvent[*events.ThirdPartyAccountCredited](),
		dogma.HandlesEvent[*events.ThirdPartyAccountCreditFailed](),
		dogma.HandlesEvent[*events.TransferApproved](),
		dogma.HandlesEvent[*events.TransferDeclined](),
		dogma.HandlesEvent[*events.TransferFailed](),
		dogma.ExecutesCommand[*commands.DebitAccount](),
		dogma.ExecutesCommand[*commands.ConsumeDailyDebitLimit](),
		dogma.ExecutesCommand[*commands.CreditAccount](),
		dogma.ExecutesCommand[*commands.CreditThirdPartyAccount](),
		dogma.ExecutesCommand[*commands.ApproveTransfer](),
		dogma.ExecutesCommand[*commands.DeclineTransfer](),
		dogma.ExecutesCommand[*commands.MarkTransferAsFailed](),
		dogma.SchedulesTimeout[*TransferReadyToProceed](),
	)
}

// RouteEventToInstance returns the ID of the process instance that is targetted
// by m.
func (TransferProcessHandler) RouteEventToInstance(
	_ context.Context,
	m dogma.Event,
) (string, bool, error) {
	switch x := m.(type) {
	case *events.TransferStarted:
		return x.TransactionID, true, nil
	case *events.AccountDebited:
		return x.TransactionID, x.TransactionType == messages.Transfer, nil
	case *events.AccountDebitDeclined:
		return x.TransactionID, x.TransactionType == messages.Transfer, nil
	case *events.DailyDebitLimitConsumed:
		return x.TransactionID, x.DebitType == messages.Transfer, nil
	case *events.DailyDebitLimitExceeded:
		return x.TransactionID, x.DebitType == messages.Transfer, nil
	case *events.AccountCredited:
		return x.TransactionID, x.TransactionType == messages.Transfer, nil
	case *events.ThirdPartyAccountCredited:
		return x.TransactionID, true, nil
	case *events.ThirdPartyAccountCreditFailed:
		return x.TransactionID, true, nil
	case *events.TransferApproved:
		return x.TransactionID, true, nil
	case *events.TransferDeclined:
		return x.TransactionID, true, nil
	case *events.TransferFailed:
		return x.TransactionID, true, nil
	default:
		panic(dogma.UnexpectedMessage)
	}
}

// HandleEvent handles an event message that has been routed to this handler.
func (TransferProcessHandler) HandleEvent(
	_ context.Context,
	r dogma.ProcessRoot,
	s dogma.ProcessEventScope,
	m dogma.Event,
) error {
	t := r.(*transferProcess)

	switch x := m.(type) {
	case *events.TransferStarted:
		t.FromAccountID = x.FromAccountID
		t.ToAccountID = x.ToAccountID
		t.ToThirdPartyBank = x.ToThirdPartyBank
		t.Amount = x.Amount

		s.ScheduleTimeout(
			&TransferReadyToProceed{
				TransactionID: x.TransactionID,
			},
			x.ScheduledTime,
		)

	case *events.AccountDebited:
		s.ExecuteCommand(&commands.ConsumeDailyDebitLimit{
			TransactionID: x.TransactionID,
			AccountID:     x.AccountID,
			DebitType:     messages.Transfer,
			Amount:        x.Amount,
			Date:          messages.DailyDebitLimitDate(x.ScheduledTime),
		})

	case *events.AccountDebitDeclined:
		s.ExecuteCommand(&commands.DeclineTransfer{
			TransactionID: x.TransactionID,
			FromAccountID: t.FromAccountID,
			ToAccountID:   t.ToAccountID,
			Amount:        x.Amount,
			Reason:        x.Reason,
		})

	case *events.DailyDebitLimitConsumed:
		if t.ToThirdPartyBank {
			s.ExecuteCommand(&commands.CreditThirdPartyAccount{
				TransactionID: x.TransactionID,
				AccountID:     t.ToAccountID,
				Amount:        x.Amount,
			})
		} else {
			s.ExecuteCommand(&commands.CreditAccount{
				TransactionID:   x.TransactionID,
				AccountID:       t.ToAccountID,
				TransactionType: messages.Transfer,
				Amount:          x.Amount,
			})
		}

	case *events.DailyDebitLimitExceeded:
		t.DeclineReason = messages.DailyDebitLimitExceeded

		// compensate the initial debit
		s.ExecuteCommand(&commands.CreditAccount{
			TransactionID:   x.TransactionID,
			AccountID:       t.FromAccountID,
			TransactionType: messages.Transfer,
			Amount:          x.Amount,
		})

	case *events.AccountCredited:
		if t.ToAccountID == x.AccountID {
			// it was a credit to complete the transfer (success)
			s.ExecuteCommand(&commands.ApproveTransfer{
				TransactionID: x.TransactionID,
				FromAccountID: t.FromAccountID,
				ToAccountID:   t.ToAccountID,
				Amount:        x.Amount,
			})
		} else {
			// it was a compensating credit to undo the transfer (business rejection)
			s.ExecuteCommand(&commands.DeclineTransfer{
				TransactionID: x.TransactionID,
				FromAccountID: t.FromAccountID,
				ToAccountID:   t.ToAccountID,
				Amount:        x.Amount,
				Reason:        t.DeclineReason,
			})
		}

	case *events.ThirdPartyAccountCredited:
		s.ExecuteCommand(&commands.ApproveTransfer{
			TransactionID: x.TransactionID,
			FromAccountID: t.FromAccountID,
			ToAccountID:   t.ToAccountID,
			Amount:        t.Amount,
		})

	case *events.ThirdPartyAccountCreditFailed:
		s.ExecuteCommand(&commands.MarkTransferAsFailed{
			TransactionID: x.TransactionID,
			FromAccountID: t.FromAccountID,
			ToAccountID:   t.ToAccountID,
			Amount:        t.Amount,
		})

		s.ExecuteCommand(&commands.CreditAccount{
			TransactionID:   x.TransactionID,
			AccountID:       t.FromAccountID,
			TransactionType: messages.Transfer,
			Amount:          t.Amount,
		})

	case *events.TransferApproved, *events.TransferDeclined, *events.TransferFailed:
		s.End()

	default:
		panic(dogma.UnexpectedMessage)
	}

	return nil
}

// HandleTimeout handles a timeout message that has been routed to this handler.
func (TransferProcessHandler) HandleTimeout(
	_ context.Context,
	r dogma.ProcessRoot,
	s dogma.ProcessTimeoutScope,
	m dogma.Timeout,
) error {
	t := r.(*transferProcess)

	switch x := m.(type) {
	case *TransferReadyToProceed:
		s.ExecuteCommand(&commands.DebitAccount{
			TransactionID:   x.TransactionID,
			AccountID:       t.FromAccountID,
			TransactionType: messages.Transfer,
			Amount:          t.Amount,
			ScheduledTime:   s.ScheduledFor(),
		})

	default:
		panic(dogma.UnexpectedMessage)
	}

	return nil
}

// TransferReadyToProceed is a timeout message notifiying that the transfer is
// ready to proceed.
type TransferReadyToProceed struct {
	TransactionID string
}

// MessageDescription returns a human-readable description of the message.
func (m *TransferReadyToProceed) MessageDescription() string {
	return fmt.Sprintf("transfer %s is ready to proceed", m.TransactionID)
}

// Validate returns a non-nil error if the message is invalid.
func (m *TransferReadyToProceed) Validate(dogma.TimeoutValidationScope) error {
	if m.TransactionID == "" {
		return errors.New("TransferReadyToProceed must not have an empty transaction ID")
	}
	return nil
}

// MarshalBinary returns a binary representation of the message.
// For simplicity in this example we use JSON.
func (m *TransferReadyToProceed) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary populates the message from its binary representation.
// For simplicity in this example we use JSON.
func (m *TransferReadyToProceed) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}
