package events

import "github.com/dogmatiq/example/messages"

// DailyDebitLimitConsumed is an event that indicates an amount of an account
// daily debit limit has been consumed.
type DailyDebitLimitConsumed struct {
	TransactionID string
	AccountID     string
	DebitType     messages.TransactionType
	Amount        int64
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
	LimitUsed     int64
	LimitMaximum  int64
}
