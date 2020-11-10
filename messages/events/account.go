package events

import (
	"fmt"

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
func (m *AccountOpened) MessageDescription() string {
	return fmt.Sprintf(
		"account %s %s opened for customer %s",
		messages.FormatID(m.AccountID),
		m.AccountName,
		messages.FormatID(m.CustomerID),
	)
}

// MessageDescription returns a human-readable description of the message.
func (m *AccountCredited) MessageDescription() string {
	return fmt.Sprintf(
		"credited %s to account %s",
		messages.FormatAmount(m.Amount),
		messages.FormatID(m.AccountID),
	)
}

// MessageDescription returns a human-readable description of the message.
func (m *AccountDebited) MessageDescription() string {
	return fmt.Sprintf(
		"debited %s from account %s",
		messages.FormatAmount(m.Amount),
		messages.FormatID(m.AccountID),
	)
}

// MessageDescription returns a human-readable description of the message.
func (m *AccountDebitDeclined) MessageDescription() string {
	return fmt.Sprintf(
		"declined debit of %s from account %s for reason %s",
		messages.FormatAmount(m.Amount),
		messages.FormatID(m.AccountID),
		m.Reason,
	)
}
