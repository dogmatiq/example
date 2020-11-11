package commands

import (
	"fmt"

	"github.com/dogmatiq/example/messages"
)

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

// MessageDescription returns a human-readable description of the message.
func (m *Deposit) MessageDescription() string {
	return fmt.Sprintf(
		"depositing %s into account %s",
		messages.FormatAmount(m.Amount),
		m.AccountID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m *ApproveDeposit) MessageDescription() string {
	return fmt.Sprintf(
		"approving deposit of %s into account %s",
		messages.FormatAmount(m.Amount),
		m.AccountID,
	)
}
