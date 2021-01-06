package commands

import (
	"errors"
	"fmt"

	"github.com/dogmatiq/example/messages"
)

// Withdraw is a command requesting that funds be withdrawn from a bank account.
type Withdraw struct {
	TransactionID string
	AccountID     string
	Amount        int64
	ScheduledDate string
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

// Validate returns a non-nil error if the message is invalid.
func (m Withdraw) Validate() error {
	if m.TransactionID == "" {
		return errors.New("Withdraw needs a valid transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("Withdraw needs a valid account ID")
	}
	if m.Amount < 1 {
		return errors.New("Withdraw needs a valid amount")
	}
	if !messages.IsValidBusinessDate(m.ScheduledDate) {
		return errors.New("Withdraw needs a valid scheduled date")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m ApproveWithdrawal) Validate() error {
	if m.TransactionID == "" {
		return errors.New("ApproveWithdrawal needs a valid transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("ApproveWithdrawal needs a valid account ID")
	}
	if m.Amount < 1 {
		return errors.New("ApproveWithdrawal needs a valid amount")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m DeclineWithdrawal) Validate() error {
	if m.TransactionID == "" {
		return errors.New("DeclineWithdrawal needs a valid transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("DeclineWithdrawal needs a valid account ID")
	}
	if m.Amount < 1 {
		return errors.New("DeclineWithdrawal needs a valid amount")
	}
	if m.Reason == "" {
		return errors.New("DeclineWithdrawal needs a valid reason")
	}

	return nil
}
