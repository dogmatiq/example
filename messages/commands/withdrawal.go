package commands

// Withdraw is a command requesting that funds be withdrawn from a bank account.
type Withdraw struct {
	TransactionID string
	AccountID     string
	Amount        int64
}

// DebitAccountForWithdrawal is a command that requests a bank account be
// debited for a withdrawal.
type DebitAccountForWithdrawal struct {
	TransactionID string
	AccountID     string
	Amount        int64
}
