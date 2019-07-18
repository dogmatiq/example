package events

// DepositStarted is an event indicating that the process of depositing funds
// into an account has begun.
type DepositStarted struct {
	TransactionID string
	AccountID     string
	Amount        int64
}

// DepositApproved is an event that indicates a requested deposit has been
// approved.
type DepositApproved struct {
	TransactionID string
	AccountID     string
	Amount        int64
}
