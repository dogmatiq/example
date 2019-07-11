package messages

// DebitType defines reasons why a debits may fail.
type DebitType string

const (
	// Withdrawal is a withdrawal transaction type.
	Withdrawal DebitType = "withdrawal"

	// Transfer is a transfer transaction type.
	Transfer DebitType = "transfer"
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
