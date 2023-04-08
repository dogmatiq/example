package events

import (
	"errors"
	"fmt"
	"time"

	"github.com/dogmatiq/example/messages"
)

// TransferStarted is an event indicating that the process of transferring funds
// from one account to another has begun.
type TransferStarted struct {
	TransactionID string
	FromAccountID string
	ToAccountID   string
	Amount        int64
	ScheduledTime time.Time
}

// TransferApproved is an event that indicates a requested transfer has been
// approved.
type TransferApproved struct {
	TransactionID string
	FromAccountID string
	ToAccountID   string
	Amount        int64
}

// TransferDeclined is an event that indicates a requested transfer has been
// declined.
type TransferDeclined struct {
	TransactionID string
	FromAccountID string
	ToAccountID   string
	Amount        int64
	Reason        messages.DebitFailureReason
}

// MessageDescription returns a human-readable description of the message.
func (m TransferStarted) MessageDescription() string {
	return fmt.Sprintf(
		"transfer %s: started transfer of %s from account %s to account %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.FromAccountID,
		m.ToAccountID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m TransferApproved) MessageDescription() string {
	return fmt.Sprintf(
		"transfer %s: approved transfer of %s from account %s to account %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.FromAccountID,
		m.ToAccountID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m TransferDeclined) MessageDescription() string {
	return fmt.Sprintf(
		"transfer %s: declined transfer of %s from account %s to account %s: %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.FromAccountID,
		m.ToAccountID,
		m.Reason,
	)
}

// Validate returns a non-nil error if the message is invalid.
func (m TransferStarted) Validate() error {
	if m.TransactionID == "" {
		return errors.New("TransferStarted must not have an empty transaction ID")
	}
	if m.FromAccountID == "" {
		return errors.New("TransferStarted must not have an empty 'from' account ID")
	}
	if m.ToAccountID == "" {
		return errors.New("TransferStarted must not have an empty 'to' account ID")
	}
	if m.FromAccountID == m.ToAccountID {
		return errors.New("TransferStarted from account ID and to account ID must be different")
	}
	if m.Amount < 1 {
		return errors.New("TransferStarted must have a positive amount")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m TransferApproved) Validate() error {
	if m.TransactionID == "" {
		return errors.New("TransferApproved must not have an empty transaction ID")
	}
	if m.FromAccountID == "" {
		return errors.New("TransferApproved must not have an empty 'from' account ID")
	}
	if m.ToAccountID == "" {
		return errors.New("TransferApproved must not have an empty 'to' account ID")
	}
	if m.FromAccountID == m.ToAccountID {
		return errors.New("TransferApproved from account ID and to account ID must be different")
	}
	if m.Amount < 1 {
		return errors.New("TransferApproved must have a positive amount")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m TransferDeclined) Validate() error {
	if m.TransactionID == "" {
		return errors.New("TransferDeclined must not have an empty transaction ID")
	}
	if m.FromAccountID == "" {
		return errors.New("TransferDeclined must not have an empty 'from' account ID")
	}
	if m.ToAccountID == "" {
		return errors.New("TransferDeclined must not have an empty 'to' account ID")
	}
	if m.FromAccountID == m.ToAccountID {
		return errors.New("TransferDeclined from account ID and to account ID must be different")
	}
	if m.Amount < 1 {
		return errors.New("TransferDeclined must have a positive amount")
	}
	if err := m.Reason.Validate(); err != nil {
		return fmt.Errorf("TransferDeclined must have a valid reason: %w", err)
	}

	return nil
}
