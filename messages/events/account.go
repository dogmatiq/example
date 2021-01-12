package events

import (
	"fmt"
	"time"

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
	ScheduledDate   time.Time
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
