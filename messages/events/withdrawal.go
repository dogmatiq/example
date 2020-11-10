package events

import (
	"fmt"

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

// WithdrawalApproved is an event that indicates a requested withdrawal has been
// approved.
type WithdrawalApproved struct {
	TransactionID string
	AccountID     string
	Amount        int64
}

// WithdrawalDeclined is an event that indicates a requested withdrawal has been
// declined.
type WithdrawalDeclined struct {
	TransactionID string
	AccountID     string
	Amount        int64
	Reason        messages.DebitFailureReason
}

// MessageDescription returns a human-readable description of the message.
func (m *WithdrawalStarted) MessageDescription() string {
	return fmt.Sprintf(
		"started withdrawal of %s from account %s",
		messages.FormatAmount(m.Amount),
		messages.FormatID(m.AccountID),
	)
}

// MessageDescription returns a human-readable description of the message.
func (m *WithdrawalApproved) MessageDescription() string {
	return fmt.Sprintf(
		"approved withdrawal of %s from account %s",
		messages.FormatAmount(m.Amount),
		messages.FormatID(m.AccountID),
	)
}

// MessageDescription returns a human-readable description of the message.
func (m *WithdrawalDeclined) MessageDescription() string {
	return fmt.Sprintf(
		"declined withdrawal of %s from account %s for reason %s",
		messages.FormatAmount(m.Amount),
		messages.FormatID(m.AccountID),
		m.Reason,
	)
}
