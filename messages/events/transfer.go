package events

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages"
)

func init() {
	dogma.RegisterEvent[*TransferStarted]("e5a7db39-861a-4a98-b109-a6f4187ac407")
	dogma.RegisterEvent[*TransferApproved]("bcc989cc-4ec7-4175-84dc-24908ac82676")
	dogma.RegisterEvent[*TransferDeclined]("0e43679a-bf5b-4730-a4a0-543e17a67479")
}

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
func (m *TransferStarted) MessageDescription() string {
	return fmt.Sprintf(
		"transfer %s: started transfer of %s from account %s to account %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.FromAccountID,
		m.ToAccountID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m *TransferApproved) MessageDescription() string {
	return fmt.Sprintf(
		"transfer %s: approved transfer of %s from account %s to account %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.FromAccountID,
		m.ToAccountID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m *TransferDeclined) MessageDescription() string {
	return fmt.Sprintf(
		"transfer %s: declined transfer of %s from account %s to account %s: %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.FromAccountID,
		m.ToAccountID,
		m.Reason,
	)
}

// Validate returns a non-nil error if the message is invalid.
func (m *TransferStarted) Validate(dogma.EventValidationScope) error {
	if m.TransactionID == "" {
		return errors.New("TransferStarted must not have an empty transaction ID")
	}
	if m.FromAccountID == "" {
		return errors.New("TransferStarted must not have an empty 'from' account ID")
	}
	if m.ToAccountID == "" {
		return errors.New("TransferStarted must not have an empty 'to' account ID")
	}
	if m.FromAccountID == m.ToAccountID {
		return errors.New("TransferStarted from account ID and to account ID must be different")
	}
	if m.Amount < 1 {
		return errors.New("TransferStarted must have a positive amount")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m *TransferApproved) Validate(dogma.EventValidationScope) error {
	if m.TransactionID == "" {
		return errors.New("TransferApproved must not have an empty transaction ID")
	}
	if m.FromAccountID == "" {
		return errors.New("TransferApproved must not have an empty 'from' account ID")
	}
	if m.ToAccountID == "" {
		return errors.New("TransferApproved must not have an empty 'to' account ID")
	}
	if m.FromAccountID == m.ToAccountID {
		return errors.New("TransferApproved from account ID and to account ID must be different")
	}
	if m.Amount < 1 {
		return errors.New("TransferApproved must have a positive amount")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m *TransferDeclined) Validate(dogma.EventValidationScope) error {
	if m.TransactionID == "" {
		return errors.New("TransferDeclined must not have an empty transaction ID")
	}
	if m.FromAccountID == "" {
		return errors.New("TransferDeclined must not have an empty 'from' account ID")
	}
	if m.ToAccountID == "" {
		return errors.New("TransferDeclined must not have an empty 'to' account ID")
	}
	if m.FromAccountID == m.ToAccountID {
		return errors.New("TransferDeclined from account ID and to account ID must be different")
	}
	if m.Amount < 1 {
		return errors.New("TransferDeclined must have a positive amount")
	}
	if err := m.Reason.Validate(); err != nil {
		return fmt.Errorf("TransferDeclined must have a valid reason: %w", err)
	}

	return nil
}

// MarshalBinary returns a binary representation of the message.
// For simplicity in this example we use JSON.
func (m *TransferStarted) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary populates the message from its binary representation.
// For simplicity in this example we use JSON.
func (m *TransferStarted) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

// MarshalBinary returns a binary representation of the message.
// For simplicity in this example we use JSON.
func (m *TransferApproved) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary populates the message from its binary representation.
// For simplicity in this example we use JSON.
func (m *TransferApproved) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

// MarshalBinary returns a binary representation of the message.
// For simplicity in this example we use JSON.
func (m *TransferDeclined) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary populates the message from its binary representation.
// For simplicity in this example we use JSON.
func (m *TransferDeclined) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}
