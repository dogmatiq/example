package commands

import (
	"time"

	"github.com/dogmatiq/example/messages"
)

// Withdraw is a command requesting that funds be withdrawn from a bank account.
type Withdraw struct {
	TransactionID string
	AccountID     string
	Amount        int64
	ScheduledDate time.Time
}

// HoldFundsForWithdrawal is a command that requests bank account funds be held
// for a withdrawal.
type HoldFundsForWithdrawal struct {
	TransactionID string
	AccountID     string
	Amount        int64
	ScheduledDate time.Time
}

// SettleWithdrawal is a command that completes an account withdrawal.
type SettleWithdrawal struct {
	TransactionID string
	AccountID     string
	Amount        int64
}

// DeclineWithdrawal is a command that rejects an account withdrawal.
type DeclineWithdrawal struct {
	TransactionID string
	AccountID     string
	Amount        int64
	Reason        messages.DebitFailureReason
}
