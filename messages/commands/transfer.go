package commands

import (
	"errors"
	"fmt"

	"github.com/dogmatiq/example/messages"
	"github.com/dogmatiq/example/messages/internal/validation"
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
func (m Transfer) MessageDescription() string {
	return fmt.Sprintf(
		"transfer %s: transfering %s from account %s to account %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.FromAccountID,
		m.ToAccountID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m ApproveTransfer) MessageDescription() string {
	return fmt.Sprintf(
		"transfer %s: approving transfer of %s from account %s to account %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.FromAccountID,
		m.ToAccountID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m DeclineTransfer) MessageDescription() string {
	return fmt.Sprintf(
		"transfer %s: declining transfer of %s from account %s to account %s: %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.FromAccountID,
		m.ToAccountID,
		m.Reason,
	)
}

// Validate returns a non-nil error if the message is invalid.
func (m Transfer) Validate() error {
	if m.TransactionID == "" {
		return errors.New("Transfer needs a valid transaction ID")
	}
	if m.FromAccountID == "" {
		return errors.New("Transfer needs a valid from account ID")
	}
	if m.ToAccountID == "" {
		return errors.New("Transfer needs a valid to account ID")
	}
	if m.FromAccountID == m.ToAccountID {
		return errors.New("Transfer from account ID and to account ID must be different")
	}
	if m.Amount < 1 {
		return errors.New("Transfer needs a valid amount")
	}
	if !validation.IsValidBusinessDate(m.ScheduledDate) {
		return errors.New("Transfer needs a valid scheduled date")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m ApproveTransfer) Validate() error {
	if m.TransactionID == "" {
		return errors.New("ApproveTransfer needs a valid transaction ID")
	}
	if m.FromAccountID == "" {
		return errors.New("ApproveTransfer needs a valid from account ID")
	}
	if m.ToAccountID == "" {
		return errors.New("ApproveTransfer needs a valid to account ID")
	}
	if m.Amount < 1 {
		return errors.New("ApproveTransfer needs a valid amount")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m DeclineTransfer) Validate() error {
	if m.TransactionID == "" {
		return errors.New("DeclineTransfer needs a valid transaction ID")
	}
	if m.FromAccountID == "" {
		return errors.New("DeclineTransfer needs a valid from account ID")
	}
	if m.ToAccountID == "" {
		return errors.New("DeclineTransfer needs a valid to account ID")
	}
	if m.Amount < 1 {
		return errors.New("DeclineTransfer needs a valid amount")
	}
	if m.Reason == "" {
		return errors.New("DeclineTransfer needs a valid reason")
	}

	return nil
}
