package commands

// Deposit is a command requesting that funds be deposited into a bank account.
type Deposit struct {
	TransactionID string
	AccountID     string
	Amount        int64
}

// ApproveDeposit is a command that approves an account deposit.
type ApproveDeposit struct {
	TransactionID string
	AccountID     string
	Amount        int64
}
