package messages

import "time"

// TransactionType defines types of debits.
type TransactionType string

const (
	// Deposit is a deposit transaction type.
	Deposit TransactionType = "deposit"

	// Withdrawal is a withdrawal transaction type.
	Withdrawal TransactionType = "withdrawal"

	// Transfer is a transfer transaction type.
	Transfer TransactionType = "transfer"
)

// DebitFailureReason defines reasons why a debits may fail.
type DebitFailureReason string

const (
	// InsufficientFunds means there was not enough funds available in the
	// account to perform the debit.
	InsufficientFunds DebitFailureReason = "insufficent funds"

	// DailyDebitLimitExceeded means that the debit cannot be performed
	// because it will exceed the account daily debit limit.
	DailyDebitLimitExceeded DebitFailureReason = "daily debit limit exceeded"
)

// DailyDebitLimitDate returns the date of a transaction for the purposes of
// checking daily debit limits.
//
// It normalizes the date to the UTC timezone so that regardless of which
// timezone is used to schedule the transaction two equivalent times always
// contribute towards the same day's limit.
func DailyDebitLimitDate(t time.Time) string {
	return t.In(time.UTC).Format("2006-01-02")
}
