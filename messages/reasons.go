package messages

// DebitFailureReason defines reasons why a debits may fail.
type DebitFailureReason string

const (
	// ReasonInsufficientFunds means there was not enough funds available in the
	// account to perform the debit.
	ReasonInsufficientFunds DebitFailureReason = "insufficent-funds"

	// ReasonDailyDebitLimitExceeded means that the debit cannot be performed
	// because it will exceed the account daily debit limit.
	ReasonDailyDebitLimitExceeded DebitFailureReason = "daily-debit-limit-exceeded"
)
