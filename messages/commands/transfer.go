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
	dogma.RegisterCommand[*Transfer]("5ee87c7b-bde3-4b39-9f12-44968cdb9889")
	dogma.RegisterCommand[*ApproveTransfer]("0d22aaa5-4449-459a-b9b1-c5fb0ce4a990")
	dogma.RegisterCommand[*DeclineTransfer]("d7d069a2-41fc-415e-91dd-7db3affa9f6d")
}

// Transfer is a command requesting that funds be transferred from one bank
// account to another.
type Transfer struct {
	TransactionID string
	FromAccountID string
	ToAccountID   string
	Amount        int64
	ScheduledTime time.Time
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
		"transfer %s: transfering %s from account %s to account %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.FromAccountID,
		m.ToAccountID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m *ApproveTransfer) MessageDescription() string {
	return fmt.Sprintf(
		"transfer %s: approving transfer of %s from account %s to account %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.FromAccountID,
		m.ToAccountID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m *DeclineTransfer) MessageDescription() string {
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
func (m *Transfer) Validate(dogma.CommandValidationScope) error {
	if m.TransactionID == "" {
		return errors.New("Transfer must not have an empty transaction ID")
	}
	if m.FromAccountID == "" {
		return errors.New("Transfer must not have an empty 'from' account ID")
	}
	if m.ToAccountID == "" {
		return errors.New("Transfer must not have an empty 'to' account ID")
	}
	if m.FromAccountID == m.ToAccountID {
		return errors.New("Transfer from account ID and to account ID must be different")
	}
	if m.Amount < 1 {
		return errors.New("Transfer must have a positive amount")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m *ApproveTransfer) Validate(dogma.CommandValidationScope) error {
	if m.TransactionID == "" {
		return errors.New("ApproveTransfer must not have an empty transaction ID")
	}
	if m.FromAccountID == "" {
		return errors.New("ApproveTransfer must not have an empty 'from' account ID")
	}
	if m.ToAccountID == "" {
		return errors.New("ApproveTransfer must not have an empty 'to' account ID")
	}
	if m.Amount < 1 {
		return errors.New("ApproveTransfer must have a positive amount")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m *DeclineTransfer) Validate(dogma.CommandValidationScope) error {
	if m.TransactionID == "" {
		return errors.New("DeclineTransfer must not have an empty transaction ID")
	}
	if m.FromAccountID == "" {
		return errors.New("DeclineTransfer must not have an empty 'from' account ID")
	}
	if m.ToAccountID == "" {
		return errors.New("DeclineTransfer must not have an empty 'to' account ID")
	}
	if m.Amount < 1 {
		return errors.New("DeclineTransfer must have a positive amount")
	}
	if err := m.Reason.Validate(); err != nil {
		return fmt.Errorf("DeclineTransfer must have a valid reason: %w", err)
	}

	return nil
}

// MarshalBinary returns a binary representation of the message.
// For simplicity in this example we use JSON.
func (m *Transfer) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary populates the message from its binary representation.
// For simplicity in this example we use JSON.
func (m *Transfer) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

// MarshalBinary returns a binary representation of the message.
// For simplicity in this example we use JSON.
func (m *ApproveTransfer) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary populates the message from its binary representation.
// For simplicity in this example we use JSON.
func (m *ApproveTransfer) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

// MarshalBinary returns a binary representation of the message.
// For simplicity in this example we use JSON.
func (m *DeclineTransfer) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary populates the message from its binary representation.
// For simplicity in this example we use JSON.
func (m *DeclineTransfer) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}
