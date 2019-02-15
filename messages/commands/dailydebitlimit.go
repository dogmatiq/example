package commands

import "time"

// ConsumeDailyDebitAmount is a command requesting that an amount of the daily
// debit limit be consumed. This is usually done when performing a debit
// transaction on a bank account.
type ConsumeDailyDebitAmount struct {
	TransactionID        string
	AccountID            string
	Amount               int64
	TransactionTimestamp time.Time
}

// RestoreDailyDebitAmount is a command requesting that an amount of the daily
// debit limit be restored. This is usually done when reversing an existing
// debit transaction on a bank account.
type RestoreDailyDebitAmount struct {
	TransactionID        string
	AccountID            string
	Amount               int64
	TransactionTimestamp time.Time
}
