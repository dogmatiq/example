package events

import "time"

// DailyDebitAmountConsumed is an event that indicates a requested daily debit
// amount has been consumed.
type DailyDebitAmountConsumed struct {
	TransactionID        string
	AccountID            string
	Amount               int64
	TransactionTimestamp time.Time
}

// DailyDebitAmountConsumtionRejected is an event that indicates a requested daily
// debit amount consumtion has been rejected.
type DailyDebitAmountConsumtionRejected struct {
	TransactionID string
	AccountID     string
	Amount        int64
}

// DailyDebitAmountRestored is an event that indicates a requested daily debit
// amount has been restored.
type DailyDebitAmountRestored struct {
	TransactionID string
	AccountID     string
	Amount        int64
}
