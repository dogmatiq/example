package events

import (
	"fmt"

	"github.com/dogmatiq/example/messages"
)

// DailyDebitLimitConsumed is an event that indicates an amount of an account
// daily debit limit has been consumed.
type DailyDebitLimitConsumed struct {
	TransactionID string
	AccountID     string
	DebitType     messages.TransactionType
	Amount        int64
	Date          string
	LimitUsed     int64
	LimitMaximum  int64
}

// DailyDebitLimitExceeded is an event that indicates an attempt to consume from
// an account daily debit limit has been rejected due to reaching the limit.
type DailyDebitLimitExceeded struct {
	TransactionID string
	AccountID     string
	DebitType     messages.TransactionType
	Amount        int64
	Date          string
	LimitUsed     int64
	LimitMaximum  int64
}

// MessageDescription returns a human-readable description of the message.
func (m DailyDebitLimitConsumed) MessageDescription() string {
	return fmt.Sprintf(
		"%s %s: consumed %s from %s daily debit limit of account %s",
		m.DebitType,
		m.TransactionID,
		m.Date,
		messages.FormatAmount(m.Amount),
		m.AccountID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m DailyDebitLimitExceeded) MessageDescription() string {
	return fmt.Sprintf(
		"%s %s: exceeded %s daily debit limit of account %s by %s",
		m.DebitType,
		m.TransactionID,
		m.Date,
		m.AccountID,
		messages.FormatAmount((m.LimitUsed+m.Amount)-m.LimitMaximum),
	)
}
