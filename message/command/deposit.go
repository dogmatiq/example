package command

// Deposit is a command requesting that funds be deposited into a bank account.
type Deposit struct {
	TransactionID string
	AccountID     string
	Amount        int64
}

// CreditAccountForDeposit is a command that credits a bank account with
// deposited funds.
type CreditAccountForDeposit struct {
	TransactionID string
	AccountID     string
	Amount        int64
}
