package commands

// MarkWithdrawalDeclinedDueToDailyDebitLimit is a command requesting that the
// withdrawal be marked as declined due to the daily debit limit being reached.
type MarkWithdrawalDeclinedDueToDailyDebitLimit struct {
	TransactionID string
	AccountID     string
	Amount        int64
}

// MarkTransferDeclinedDueToDailyDebitLimit is a command requesting that the
// transfer be marked as declined due to the daily debit limit being reached.
type MarkTransferDeclinedDueToDailyDebitLimit struct {
	TransactionID string
	AccountID     string
	Amount        int64
}
