package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages"
)

func init() {
	dogma.RegisterCommand[*Withdraw]("5528ea57-2fa4-4ddd-ac96-3ec56ec89fce")
	dogma.RegisterCommand[*ApproveWithdrawal]("1023356f-f49c-4ea1-9244-83fdb3225da9")
	dogma.RegisterCommand[*DeclineWithdrawal]("ee1e59ff-ed96-4d66-a1aa-f0a11ca8c47d")
}

// Withdraw is a command requesting that funds be withdrawn from a bank account.
type Withdraw struct {
	TransactionID string
	AccountID     string
	Amount        int64
	ScheduledTime time.Time
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
func (m *Withdraw) MessageDescription() string {
	return fmt.Sprintf(
		"withdrawal %s: withdrawing %s from account %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.AccountID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m *ApproveWithdrawal) MessageDescription() string {
	return fmt.Sprintf(
		"withdrawal %s: approving withdrawal of %s from account %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.AccountID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m *DeclineWithdrawal) MessageDescription() string {
	return fmt.Sprintf(
		"withdrawal %s: declining withdrawal of %s from account %s: %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.AccountID,
		m.Reason,
	)
}

// Validate returns a non-nil error if the message is invalid.
func (m *Withdraw) Validate(dogma.CommandValidationScope) error {
	if m.TransactionID == "" {
		return errors.New("Withdraw must not have an empty transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("Withdraw must not have an empty account ID")
	}
	if m.Amount < 1 {
		return errors.New("Withdraw must have a positive amount")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m *ApproveWithdrawal) Validate(dogma.CommandValidationScope) error {
	if m.TransactionID == "" {
		return errors.New("ApproveWithdrawal must not have an empty transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("ApproveWithdrawal must not have an empty account ID")
	}
	if m.Amount < 1 {
		return errors.New("ApproveWithdrawal must have a positive amount")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m *DeclineWithdrawal) Validate(dogma.CommandValidationScope) error {
	if m.TransactionID == "" {
		return errors.New("DeclineWithdrawal must not have an empty transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("DeclineWithdrawal must not have an empty account ID")
	}
	if m.Amount < 1 {
		return errors.New("DeclineWithdrawal must have a positive amount")
	}
	if err := m.Reason.Validate(); err != nil {
		return fmt.Errorf("DeclineWithdrawal must have a valid reason: %w", err)
	}

	return nil
}

// MarshalBinary returns a binary representation of the message.
// For simplicity in this example we use JSON.
func (m *Withdraw) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary populates the message from its binary representation.
// For simplicity in this example we use JSON.
func (m *Withdraw) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

// MarshalBinary returns a binary representation of the message.
// For simplicity in this example we use JSON.
func (m *ApproveWithdrawal) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary populates the message from its binary representation.
// For simplicity in this example we use JSON.
func (m *ApproveWithdrawal) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

// MarshalBinary returns a binary representation of the message.
// For simplicity in this example we use JSON.
func (m *DeclineWithdrawal) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary populates the message from its binary representation.
// For simplicity in this example we use JSON.
func (m *DeclineWithdrawal) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}
