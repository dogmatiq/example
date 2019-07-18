package messages

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
	InsufficientFunds DebitFailureReason = "insufficent-funds"

	// DailyDebitLimitExceeded means that the debit cannot be performed
	// because it will exceed the account daily debit limit.
	DailyDebitLimitExceeded DebitFailureReason = "daily-debit-limit-exceeded"
)
