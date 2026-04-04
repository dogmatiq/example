package commands

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages"
)

func init() {
	dogma.RegisterCommand[*CreditThirdPartyAccount]("c4a2e6f1-8b3d-4e7a-9f5c-2d1b0e8a3c6f")
}

// CreditThirdPartyAccount is a command to credit an account held at a
// third-party bank.
type CreditThirdPartyAccount struct {
	TransactionID string
	AccountID     string
	Amount        int64
}

// MessageDescription returns a human-readable description of the message.
func (m *CreditThirdPartyAccount) MessageDescription() string {
	return fmt.Sprintf(
		"transfer %s: crediting %s to third-party account %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.AccountID,
	)
}

// MarshalBinary returns a binary representation of the message.
// For simplicity in this example we use JSON.
func (m *CreditThirdPartyAccount) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary populates the message from its binary representation.
// For simplicity in this example we use JSON.
func (m *CreditThirdPartyAccount) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

// Validate returns a non-nil error if the message is invalid.
func (m *CreditThirdPartyAccount) Validate(dogma.CommandValidationScope) error {
	if m.TransactionID == "" {
		return errors.New("CreditThirdPartyAccount must not have an empty transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("CreditThirdPartyAccount must not have an empty account ID")
	}
	if m.Amount < 1 {
		return errors.New("CreditThirdPartyAccount must have a positive amount")
	}

	return nil
}
