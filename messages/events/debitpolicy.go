package events

// WithdrawalApprovedByDebitPolicy is an event that indicates a requested
// withdrawal has been approved by the debit policy.
type WithdrawalApprovedByDebitPolicy struct {
	TransactionID string
	AccountID     string
	Amount        int64
}

// WithdrawalDeclinedByDebitPolicy is an event that indicates a requested
// withdrawal has been declined due to the debit policy.
type WithdrawalDeclinedByDebitPolicy struct {
	TransactionID string
	AccountID     string
	Amount        int64
}

// TransferApprovedByDebitPolicy is an event that indicates a requested
// transfer has been approved by the debit policy.
type TransferApprovedByDebitPolicy struct {
	TransactionID string
	AccountID     string
	Amount        int64
}

// TransferDeclinedByDebitPolicy is an event that indicates a requested
// transfer has been declined due to the debit policy.
type TransferDeclinedByDebitPolicy struct {
	TransactionID string
	AccountID     string
	Amount        int64
}
