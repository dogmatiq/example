package event

// DepositStarted is an event indicating that the process of depositing funds
// into an account has begun.
type DepositStarted struct {
	TransactionID string
	AccountID     string
	Amount        int64
}

// AccountCreditedForDeposit is an event that indicates an account has been
// credited with funds from a deposit.
type AccountCreditedForDeposit struct {
	TransactionID string
	AccountID     string
	Amount        int64
}
