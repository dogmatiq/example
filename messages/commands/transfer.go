package commands

// Transfer is a command requesting that funds be transferred from one bank
// account to another.
type Transfer struct {
	TransactionID string
	FromAccountID string
	ToAccountID   string
	Amount        int64
}

// CreditAccountForTransfer is a command that credits a bank account with
// transferred funds.
type CreditAccountForTransfer struct {
	TransactionID string
	AccountID     string
	Amount        int64
}

// DebitAccountForTransfer is a command that requests a bank account be debited
// for a transfer.
type DebitAccountForTransfer struct {
	TransactionID string
	AccountID     string
	Amount        int64
}
