package events

import "time"

// WithdrawalStarted is an event indicating that the process of withdrawing
// funds from an account has begun.
type WithdrawalStarted struct {
	TransactionID        string
	AccountID            string
	Amount               int64
	TransactionTimestamp time.Time
}

// AccountDebitedForWithdrawal is an event that indicates an account has been
// debited funds for a withdrawal.
type AccountDebitedForWithdrawal struct {
	TransactionID        string
	AccountID            string
	Amount               int64
	TransactionTimestamp time.Time
}

// WithdrawalDeclinedDueToInsufficientFunds is an event that indicates a
// requested withdrawal has been declined due to insufficient funds.
type WithdrawalDeclinedDueToInsufficientFunds struct {
	TransactionID        string
	AccountID            string
	Amount               int64
	TransactionTimestamp time.Time
}

// WithdrawalDeclinedDueToDailyDebitLimit is an event that indicates a requested
// withdrawal has been declined due to daily debit limits.
type WithdrawalDeclinedDueToDailyDebitLimit struct {
	TransactionID string
	AccountID     string
	Amount        int64
}
