package commands

import (
	"fmt"
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

// ApproveWithdrawal is a command that approves an account withdrawal.
type ApproveWithdrawal struct {
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

// MessageDescription returns a human-readable description of the message.
func (m Withdraw) MessageDescription() string {
	return fmt.Sprintf(
		"withdrawal %s: withdrawing %s from account %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.AccountID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m ApproveWithdrawal) MessageDescription() string {
	return fmt.Sprintf(
		"withdrawal %s: approving withdrawal of %s from account %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.AccountID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m DeclineWithdrawal) MessageDescription() string {
	return fmt.Sprintf(
		"withdrawal %s: declining withdrawal of %s from account %s: %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.AccountID,
		m.Reason,
	)
}
