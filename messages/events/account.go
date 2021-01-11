package events

import (
	"errors"
	"fmt"

	"github.com/dogmatiq/example/messages"
	"github.com/dogmatiq/example/messages/internal/validation"
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
	ScheduledDate   string
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
func (m AccountOpened) Validate() error {
	if m.CustomerID == "" {
		return errors.New("AccountOpened needs a valid customer ID")
	}
	if m.AccountID == "" {
		return errors.New("AccountOpened needs a valid account ID")
	}
	if m.AccountName == "" {
		return errors.New("AccountOpened needs a valid account name")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m AccountCredited) Validate() error {
	if m.TransactionID == "" {
		return errors.New("AccountCredited needs a valid transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("AccountCredited needs a valid account ID")
	}
	if m.TransactionType == "" {
		return errors.New("AccountCredited needs a valid transaction type")
	}
	if m.Amount < 1 {
		return errors.New("AccountCredited needs a valid amount")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m AccountDebited) Validate() error {
	if m.TransactionID == "" {
		return errors.New("AccountDebited needs a valid transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("AccountDebited needs a valid account ID")
	}
	if m.TransactionType == "" {
		return errors.New("AccountDebited needs a valid transaction type")
	}
	if m.Amount < 1 {
		return errors.New("AccountDebited needs a valid amount")
	}
	if !validation.IsValidBusinessDate(m.ScheduledDate) {
		return errors.New("AccountDebited needs a valid scheduled date")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m AccountDebitDeclined) Validate() error {
	if m.TransactionID == "" {
		return errors.New("AccountDebitDeclined needs a valid transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("AccountDebitDeclined needs a valid account ID")
	}
	if m.TransactionType == "" {
		return errors.New("AccountDebitDeclined needs a valid transaction type")
	}
	if m.Amount < 1 {
		return errors.New("AccountDebitDeclineAccountDebitDeclined a valid amount")
	}
	if m.Reason == "" {
		return errors.New("AccountDebitDeclined needs a valid reason")
	}

	return nil
}
