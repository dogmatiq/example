package events

import (
	"errors"
	"fmt"
	"time"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages"
)

// AccountOpened is an event indicating that a new bank account has been opened.
type AccountOpened struct {
	CustomerID  string
	AccountID   string
	AccountName string
}

// AccountCredited is an event indicating that a bank account was credited.
type AccountCredited struct {
	TransactionID   string
	AccountID       string
	TransactionType messages.TransactionType
	Amount          int64
}

// AccountDebited is an event indicating that a bank account was debited.
type AccountDebited struct {
	TransactionID   string
	AccountID       string
	TransactionType messages.TransactionType
	Amount          int64
	ScheduledTime   time.Time
}

// AccountDebitDeclined is an event indicating that a bank account debit was
// declined.
type AccountDebitDeclined struct {
	TransactionID   string
	AccountID       string
	TransactionType messages.TransactionType
	Amount          int64
	Reason          messages.DebitFailureReason
}

// MessageDescription returns a human-readable description of the message.
func (m AccountOpened) MessageDescription() string {
	return fmt.Sprintf(
		"opened account %s %s for customer %s",
		m.AccountID,
		m.AccountName,
		m.CustomerID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m AccountCredited) MessageDescription() string {
	return fmt.Sprintf(
		"%s %s: credited %s to account %s",
		m.TransactionType,
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.AccountID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m AccountDebited) MessageDescription() string {
	return fmt.Sprintf(
		"%s %s: debited %s from account %s",
		m.TransactionType,
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.AccountID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m AccountDebitDeclined) MessageDescription() string {
	return fmt.Sprintf(
		"%s %s: declined debit of %s from account %s: %s",
		m.TransactionType,
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.AccountID,
		m.Reason,
	)
}

// Validate returns a non-nil error if the message is invalid.
func (m AccountOpened) Validate(dogma.EventValidationScope) error {
	if m.CustomerID == "" {
		return errors.New("AccountOpened must not have an empty customer ID")
	}
	if m.AccountID == "" {
		return errors.New("AccountOpened must not have an empty account ID")
	}
	if m.AccountName == "" {
		return errors.New("AccountOpened must not have an empty account name")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m AccountCredited) Validate(dogma.EventValidationScope) error {
	if m.TransactionID == "" {
		return errors.New("AccountCredited must not have an empty transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("AccountCredited must not have an empty account ID")
	}
	if err := m.TransactionType.Validate(); err != nil {
		return fmt.Errorf("AccountCredited must have a valid transaction type: %w", err)
	}
	if m.Amount < 1 {
		return errors.New("AccountCredited must have a positive amount")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m AccountDebited) Validate(dogma.EventValidationScope) error {
	if m.TransactionID == "" {
		return errors.New("AccountDebited must not have an empty transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("AccountDebited must not have an empty account ID")
	}
	if err := m.TransactionType.Validate(); err != nil {
		return fmt.Errorf("AccountDebited must have a valid transaction type: %w", err)
	}
	if m.Amount < 1 {
		return errors.New("AccountDebited must have a positive amount")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m AccountDebitDeclined) Validate(dogma.EventValidationScope) error {
	if m.TransactionID == "" {
		return errors.New("AccountDebitDeclined must not have an empty transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("AccountDebitDeclined must not have an empty account ID")
	}
	if err := m.TransactionType.Validate(); err != nil {
		return fmt.Errorf("AccountDebitDeclined must have a valid transaction type: %w", err)
	}
	if m.Amount < 1 {
		return errors.New("AccountDebitDeclined must have a positive amount")
	}
	if err := m.Reason.Validate(); err != nil {
		return fmt.Errorf("AccountDebitDeclined must have a valid reason: %w", err)
	}

	return nil
}
