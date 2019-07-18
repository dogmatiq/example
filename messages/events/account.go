package events

import "github.com/dogmatiq/example/messages"

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
