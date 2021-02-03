package events

import (
	"errors"
	"fmt"
	"time"

	"github.com/dogmatiq/example/messages"
)

// WithdrawalStarted is an event indicating that the process of withdrawing
// funds from an account has begun.
type WithdrawalStarted struct {
	TransactionID string
	AccountID     string
	Amount        int64
	ScheduledTime time.Time
}

// WithdrawalApproved is an event that indicates a requested withdrawal has been
// approved.
type WithdrawalApproved struct {
	TransactionID string
	AccountID     string
	Amount        int64
}

// WithdrawalDeclined is an event that indicates a requested withdrawal has been
// declined.
type WithdrawalDeclined struct {
	TransactionID string
	AccountID     string
	Amount        int64
	Reason        messages.DebitFailureReason
}

// MessageDescription returns a human-readable description of the message.
func (m WithdrawalStarted) MessageDescription() string {
	return fmt.Sprintf(
		"withdrawal %s: started withdrawal of %s from account %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.AccountID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m WithdrawalApproved) MessageDescription() string {
	return fmt.Sprintf(
		"withdrawal %s: approved withdrawal of %s from account %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.AccountID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m WithdrawalDeclined) MessageDescription() string {
	return fmt.Sprintf(
		"withdrawal %s: declined withdrawal of %s from account %s: %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.AccountID,
		m.Reason,
	)
}

// Validate returns a non-nil error if the message is invalid.
func (m WithdrawalStarted) Validate() error {
	if m.TransactionID == "" {
		return errors.New("WithdrawalStarted needs a valid transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("WithdrawalStarted needs a valid account ID")
	}
	if m.Amount < 1 {
		return errors.New("WithdrawalStarted needs a valid amount")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m WithdrawalApproved) Validate() error {
	if m.TransactionID == "" {
		return errors.New("WithdrawalApproved needs a valid transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("WithdrawalApproved needs a valid account ID")
	}
	if m.Amount < 1 {
		return errors.New("WithdrawalApproved needs a valid amount")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m WithdrawalDeclined) Validate() error {
	if m.TransactionID == "" {
		return errors.New("WithdrawalDeclined needs a valid transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("WithdrawalDeclined needs a valid account ID")
	}
	if m.Amount < 1 {
		return errors.New("WithdrawalDeclined needs a valid amount")
	}
	if m.Reason == "" {
		return errors.New("WithdrawalDeclined needs a valid reason")
	}

	return nil
}
