package events

// TransferStarted is an event indicating that the process of transferring funds
// from one account to another has begun.
type TransferStarted struct {
	TransactionID string
	FromAccountID string
	ToAccountID   string
	Amount        int64
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
	TransactionID string
	AccountID     string
	Amount        int64
}

// TransferDeclined is an event that indicates a requested transfer has been
// declined due to insufficient funds.
type TransferDeclined struct {
	TransactionID string
	AccountID     string
	Amount        int64
}
