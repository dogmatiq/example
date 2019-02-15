package events

import "time"

// TransferStarted is an event indicating that the process of transferring funds
// from one account to another has begun.
type TransferStarted struct {
	TransactionID        string
	FromAccountID        string
	ToAccountID          string
	Amount               int64
	TransactionTimestamp time.Time
}

// AccountCreditedForTransfer is an event that indicates an account has been
// credited with funds from a transfer.
type AccountCreditedForTransfer struct {
	TransactionID string
	AccountID     string
	Amount        int64
}

// AccountDebitedForTransfer is an event that indicates an account has been
// debited funds for a transfer.
type AccountDebitedForTransfer struct {
	TransactionID        string
	AccountID            string
	Amount               int64
	TransactionTimestamp time.Time
}

// TransferDeclinedDueToInsufficientFunds is an event that indicates a requested
// transfer has been declined due to insufficient funds.
type TransferDeclinedDueToInsufficientFunds struct {
	TransactionID        string
	AccountID            string
	Amount               int64
	TransactionTimestamp time.Time
}

// TransferDeclinedDueToDailyDebitLimit is an event that indicates a requested
// transfer has been declined due to daily debit limits.
type TransferDeclinedDueToDailyDebitLimit struct {
	TransactionID string
	AccountID     string
	Amount        int64
}
