package commands

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages"
	"github.com/dogmatiq/example/messages/internal/validation"
)

func init() {
	dogma.RegisterCommand[*ConsumeDailyDebitLimit]("d86ca816-6333-4b78-a1d3-9368b3adcf65")
}

// ConsumeDailyDebitLimit is a command requesting that an amount of an account
// daily debit limit be consumed.
type ConsumeDailyDebitLimit struct {
	TransactionID string
	AccountID     string
	DebitType     messages.TransactionType
	Amount        int64
	Date          string
}

// MessageDescription returns a human-readable description of the message.
func (m *ConsumeDailyDebitLimit) MessageDescription() string {
	return fmt.Sprintf(
		"%s %s: consuming %s from %s daily debit limit of account %s",
		m.DebitType,
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.Date,
		m.AccountID,
	)
}

// Validate returns a non-nil error if the message is invalid.
func (m *ConsumeDailyDebitLimit) Validate(dogma.CommandValidationScope) error {
	if m.TransactionID == "" {
		return errors.New("ConsumeDailyDebitLimit must not have an empty transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("ConsumeDailyDebitLimit must not have an empty account ID")
	}
	if err := m.DebitType.Validate(); err != nil {
		return fmt.Errorf("ConsumeDailyDebitLimit must have a valid transaction type: %w", err)
	}
	if !m.DebitType.IsDebit() {
		return errors.New("ConsumeDailyDebitLimit must have a debit transaction type")
	}
	if m.Amount < 1 {
		return errors.New("ConsumeDailyDebitLimit must have a positive amount")
	}
	if !validation.IsValidDate(m.Date) {
		return errors.New("ConsumeDailyDebitLimit must have a valid date")
	}

	return nil
}

// MarshalBinary returns a binary representation of the message.
// For simplicity in this example we use JSON.
func (m *ConsumeDailyDebitLimit) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary populates the message from its binary representation.
// For simplicity in this example we use JSON.
func (m *ConsumeDailyDebitLimit) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}
