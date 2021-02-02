package events

import (
	"fmt"
	"time"

	"github.com/dogmatiq/example/messages"
)

// TransferStarted is an event indicating that the process of transferring funds
// from one account to another has begun.
type TransferStarted struct {
	TransactionID string
	FromAccountID string
	ToAccountID   string
	Amount        int64
	ScheduledTime time.Time
}

// TransferApproved is an event that indicates a requested transfer has been
// approved.
type TransferApproved struct {
	TransactionID string
	FromAccountID string
	ToAccountID   string
	Amount        int64
}

// TransferDeclined is an event that indicates a requested transfer has been
// declined.
type TransferDeclined struct {
	TransactionID string
	FromAccountID string
	ToAccountID   string
	Amount        int64
	Reason        messages.DebitFailureReason
}

// MessageDescription returns a human-readable description of the message.
func (m TransferStarted) MessageDescription() string {
	return fmt.Sprintf(
		"transfer %s: started transfer of %s from account %s to account %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.FromAccountID,
		m.ToAccountID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m TransferApproved) MessageDescription() string {
	return fmt.Sprintf(
		"transfer %s: approved transfer of %s from account %s to account %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.FromAccountID,
		m.ToAccountID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m TransferDeclined) MessageDescription() string {
	return fmt.Sprintf(
		"transfer %s: declined transfer of %s from account %s to account %s: %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.FromAccountID,
		m.ToAccountID,
		m.Reason,
	)
}
