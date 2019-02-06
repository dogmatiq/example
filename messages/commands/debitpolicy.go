package commands

import "time"

// CheckWithdrawalAllowedByDebitPolicy is a command requesting approval from the
// debit policy for funds to be withdrawn from a bank account for a withdrawal.
type CheckWithdrawalAllowedByDebitPolicy struct {
	Timestamp     time.Time
	TransactionID string
	AccountID     string
	Amount        int64
}

// CheckTransferAllowedByDebitPolicy is a command requesting approval from the
// debit policy for funds to be withdrawn froma bank about for a transfer.
type CheckTransferAllowedByDebitPolicy struct {
	Timestamp     time.Time
	TransactionID string
	AccountID     string
	Amount        int64
}
