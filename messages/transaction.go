package messages

import (
	"fmt"
	"time"
)

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

// IsDebit returns true if the transaction type is a debit type.
func (t TransactionType) IsDebit() bool {
	switch t {
	case Withdrawal:
		return true
	case Transfer:
		return true
	}

	return false
}

// Validate return an error if t is not a valid transaction type.
func (t TransactionType) Validate() error {
	switch t {
	case Deposit,
		Withdrawal,
		Transfer:
		return nil
	default:
		return fmt.Errorf("invalid transaction type: %s", string(t))
	}
}

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

// Validate return an error if r is not a valid reason.
func (r DebitFailureReason) Validate() error {
	switch r {
	case InsufficientFunds,
		DailyDebitLimitExceeded:
		return nil
	default:
		return fmt.Errorf("invalid transaction type: %s", string(r))
	}
}

// DailyDebitLimitDate returns the date of a transaction for the purposes of
// checking daily debit limits.
//
// It normalizes the date to the UTC timezone so that regardless of which
// timezone is used to schedule the transaction two equivalent times always
// contribute towards the same day's limit.
func DailyDebitLimitDate(t time.Time) string {
	return t.In(time.UTC).Format("2006-01-02")
}
