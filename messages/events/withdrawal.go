package events

import (
	"github.com/dogmatiq/example/messages"
)

// WithdrawalStarted is an event indicating that the process of withdrawing
// funds from an account has begun.
type WithdrawalStarted struct {
	TransactionID string
	AccountID     string
	Amount        int64
	ScheduledDate string
}

// FundsHeldForWithdrawal is an event that indicates account funds have been
// held for a withdrawal.
type FundsHeldForWithdrawal struct {
	TransactionID string
	AccountID     string
	Amount        int64
	ScheduledDate string
}

// AccountDebitedForWithdrawal is an event that indicates an account has been
// debited funds due to a withdrawal.
type AccountDebitedForWithdrawal struct {
	TransactionID string
	AccountID     string
	Amount        int64
}

// WithdrawalDeclined is an event that indicates a requested withdrawal has been
// declined due to insufficient funds.
type WithdrawalDeclined struct {
	TransactionID string
	AccountID     string
	Amount        int64
	Reason        messages.DebitFailureReason
}
