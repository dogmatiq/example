package commands

import (
	"fmt"

	"github.com/dogmatiq/example/messages"
)

// Transfer is a command requesting that funds be transferred from one bank
// account to another.
type Transfer struct {
	TransactionID string
	FromAccountID string
	ToAccountID   string
	Amount        int64
	ScheduledDate string
}

// ApproveTransfer is a command that approves an account transfer.
type ApproveTransfer struct {
	TransactionID string
	FromAccountID string
	ToAccountID   string
	Amount        int64
}

// DeclineTransfer is a command that rejects an account transfer.
type DeclineTransfer struct {
	TransactionID string
	FromAccountID string
	ToAccountID   string
	Amount        int64
	Reason        messages.DebitFailureReason
}

// MessageDescription returns a human-readable description of the message.
func (m *Transfer) MessageDescription() string {
	return fmt.Sprintf(
		"transfering %s from account %s to account %s",
		messages.FormatAmount(m.Amount),
		messages.FormatID(m.FromAccountID),
		messages.FormatID(m.ToAccountID),
	)
}

// MessageDescription returns a human-readable description of the message.
func (m *ApproveTransfer) MessageDescription() string {
	return fmt.Sprintf(
		"approving transfer of %s from account %s to account %s",
		messages.FormatAmount(m.Amount),
		messages.FormatID(m.FromAccountID),
		messages.FormatID(m.ToAccountID),
	)
}

// MessageDescription returns a human-readable description of the message.
func (m *DeclineTransfer) MessageDescription() string {
	return fmt.Sprintf(
		"declining transfer of %s from account %s to account %s for reason %s",
		messages.FormatAmount(m.Amount),
		messages.FormatID(m.FromAccountID),
		messages.FormatID(m.ToAccountID),
		m.Reason,
	)
}
