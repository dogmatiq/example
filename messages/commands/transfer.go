package commands

import "github.com/dogmatiq/example/messages"

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
