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
	dogma.RegisterEvent[*AccountOpened]("75ef425c-3c42-4ead-8925-cd26cbea3139")
	dogma.RegisterEvent[*AccountCredited]("f3091f46-7d36-4e5f-b0ab-fb96029e5d7a")
	dogma.RegisterEvent[*AccountDebited]("76552351-0095-442e-85dd-68f8f7fae286")
	dogma.RegisterEvent[*AccountDebitDeclined]("34b0c426-5467-44d8-a5db-7d09c789c930")
}

// AccountOpened is an event indicating that a new bank account has been opened.
type AccountOpened struct {
	CustomerID  string
	AccountID   string
	AccountName string
}

// AccountCredited is an event indicating that a bank account was credited.
type AccountCredited struct {
	TransactionID   string
	AccountID       string
	TransactionType messages.TransactionType
	Amount          int64
}

// AccountDebited is an event indicating that a bank account was debited.
type AccountDebited struct {
	TransactionID   string
	AccountID       string
	TransactionType messages.TransactionType
	Amount          int64
	ScheduledTime   time.Time
}

// AccountDebitDeclined is an event indicating that a bank account debit was
// declined.
type AccountDebitDeclined struct {
	TransactionID   string
	AccountID       string
	TransactionType messages.TransactionType
	Amount          int64
	Reason          messages.DebitFailureReason
}

// MessageDescription returns a human-readable description of the message.
func (m *AccountOpened) MessageDescription() string {
	return fmt.Sprintf(
		"opened account %s %s for customer %s",
		m.AccountID,
		m.AccountName,
		m.CustomerID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m *AccountCredited) MessageDescription() string {
	return fmt.Sprintf(
		"%s %s: credited %s to account %s",
		m.TransactionType,
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.AccountID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m *AccountDebited) MessageDescription() string {
	return fmt.Sprintf(
		"%s %s: debited %s from account %s",
		m.TransactionType,
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.AccountID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m *AccountDebitDeclined) MessageDescription() string {
	return fmt.Sprintf(
		"%s %s: declined debit of %s from account %s: %s",
		m.TransactionType,
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.AccountID,
		m.Reason,
	)
}

// Validate returns a non-nil error if the message is invalid.
func (m *AccountOpened) Validate(dogma.EventValidationScope) error {
	if m.CustomerID == "" {
		return errors.New("AccountOpened must not have an empty customer ID")
	}
	if m.AccountID == "" {
		return errors.New("AccountOpened must not have an empty account ID")
	}
	if m.AccountName == "" {
		return errors.New("AccountOpened must not have an empty account name")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m *AccountCredited) Validate(dogma.EventValidationScope) error {
	if m.TransactionID == "" {
		return errors.New("AccountCredited must not have an empty transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("AccountCredited must not have an empty account ID")
	}
	if err := m.TransactionType.Validate(); err != nil {
		return fmt.Errorf("AccountCredited must have a valid transaction type: %w", err)
	}
	if m.Amount < 1 {
		return errors.New("AccountCredited must have a positive amount")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m *AccountDebited) Validate(dogma.EventValidationScope) error {
	if m.TransactionID == "" {
		return errors.New("AccountDebited must not have an empty transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("AccountDebited must not have an empty account ID")
	}
	if err := m.TransactionType.Validate(); err != nil {
		return fmt.Errorf("AccountDebited must have a valid transaction type: %w", err)
	}
	if m.Amount < 1 {
		return errors.New("AccountDebited must have a positive amount")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m *AccountDebitDeclined) Validate(dogma.EventValidationScope) error {
	if m.TransactionID == "" {
		return errors.New("AccountDebitDeclined must not have an empty transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("AccountDebitDeclined must not have an empty account ID")
	}
	if err := m.TransactionType.Validate(); err != nil {
		return fmt.Errorf("AccountDebitDeclined must have a valid transaction type: %w", err)
	}
	if m.Amount < 1 {
		return errors.New("AccountDebitDeclined must have a positive amount")
	}
	if err := m.Reason.Validate(); err != nil {
		return fmt.Errorf("AccountDebitDeclined must have a valid reason: %w", err)
	}

	return nil
}

// MarshalBinary returns a binary representation of the message.
// For simplicity in this example we use JSON.
func (m *AccountOpened) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary populates the message from its binary representation.
// For simplicity in this example we use JSON.
func (m *AccountOpened) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

// MarshalBinary returns a binary representation of the message.
// For simplicity in this example we use JSON.
func (m *AccountCredited) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary populates the message from its binary representation.
// For simplicity in this example we use JSON.
func (m *AccountCredited) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

// MarshalBinary returns a binary representation of the message.
// For simplicity in this example we use JSON.
func (m *AccountDebited) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary populates the message from its binary representation.
// For simplicity in this example we use JSON.
func (m *AccountDebited) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

// MarshalBinary returns a binary representation of the message.
// For simplicity in this example we use JSON.
func (m *AccountDebitDeclined) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary populates the message from its binary representation.
// For simplicity in this example we use JSON.
func (m *AccountDebitDeclined) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}
