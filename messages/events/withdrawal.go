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
	dogma.RegisterEvent[*WithdrawalStarted]("ae327c68-3de9-4fe1-aa39-c4830200602c")
	dogma.RegisterEvent[*WithdrawalApproved]("8551621a-ef28-42aa-b815-d5a448299d38")
	dogma.RegisterEvent[*WithdrawalDeclined]("a87f5f6f-52ec-4eeb-954f-2af97a165181")
}

// WithdrawalStarted is an event indicating that the process of withdrawing
// funds from an account has begun.
type WithdrawalStarted struct {
	TransactionID string
	AccountID     string
	Amount        int64
	ScheduledTime time.Time
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
		"withdrawal %s: started withdrawal of %s from account %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.AccountID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m *WithdrawalApproved) MessageDescription() string {
	return fmt.Sprintf(
		"withdrawal %s: approved withdrawal of %s from account %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.AccountID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m *WithdrawalDeclined) MessageDescription() string {
	return fmt.Sprintf(
		"withdrawal %s: declined withdrawal of %s from account %s: %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.AccountID,
		m.Reason,
	)
}

// Validate returns a non-nil error if the message is invalid.
func (m *WithdrawalStarted) Validate(dogma.EventValidationScope) error {
	if m.TransactionID == "" {
		return errors.New("WithdrawalStarted must not have an empty transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("WithdrawalStarted must not have an empty account ID")
	}
	if m.Amount < 1 {
		return errors.New("WithdrawalStarted must have a positive amount")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m *WithdrawalApproved) Validate(dogma.EventValidationScope) error {
	if m.TransactionID == "" {
		return errors.New("WithdrawalApproved must not have an empty transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("WithdrawalApproved must not have an empty account ID")
	}
	if m.Amount < 1 {
		return errors.New("WithdrawalApproved must have a positive amount")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m *WithdrawalDeclined) Validate(dogma.EventValidationScope) error {
	if m.TransactionID == "" {
		return errors.New("WithdrawalDeclined must not have an empty transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("WithdrawalDeclined must not have an empty account ID")
	}
	if m.Amount < 1 {
		return errors.New("WithdrawalDeclined must have a positive amount")
	}
	if err := m.Reason.Validate(); err != nil {
		return fmt.Errorf("WithdrawalDeclined must have a valid reason: %w", err)
	}

	return nil
}

// MarshalBinary returns a binary representation of the message.
// For simplicity in this example we use JSON.
func (m *WithdrawalStarted) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary populates the message from its binary representation.
// For simplicity in this example we use JSON.
func (m *WithdrawalStarted) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

// MarshalBinary returns a binary representation of the message.
// For simplicity in this example we use JSON.
func (m *WithdrawalApproved) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary populates the message from its binary representation.
// For simplicity in this example we use JSON.
func (m *WithdrawalApproved) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

// MarshalBinary returns a binary representation of the message.
// For simplicity in this example we use JSON.
func (m *WithdrawalDeclined) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary populates the message from its binary representation.
// For simplicity in this example we use JSON.
func (m *WithdrawalDeclined) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}
