package events

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages"
)

func init() {
	dogma.RegisterEvent[*ThirdPartyAccountCredited]("bfecec06-c8ce-45a3-ba10-4d4e3649dc2e")
	dogma.RegisterEvent[*ThirdPartyAccountCreditFailed]("7970983c-2ba2-440b-8c25-88c17e9705e0")
}

// ThirdPartyAccountCredited is an event indicating that the credit to a
// third-party bank account was completed successfully.
type ThirdPartyAccountCredited struct {
	TransactionID string
	AccountID     string
	Amount        int64
}

// ThirdPartyAccountCreditFailed is an event indicating that the credit to a
// third-party bank account could not be completed.
type ThirdPartyAccountCreditFailed struct {
	TransactionID string
	AccountID     string
	Amount        int64
}

// MessageDescription returns a human-readable description of the message.
func (m *ThirdPartyAccountCredited) MessageDescription() string {
	return fmt.Sprintf(
		"transfer %s: credited %s to third-party account %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.AccountID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m *ThirdPartyAccountCreditFailed) MessageDescription() string {
	return fmt.Sprintf(
		"transfer %s: failed to credit %s to third-party account %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.AccountID,
	)
}

// MarshalBinary returns a binary representation of the message.
// For simplicity in this example we use JSON.
func (m *ThirdPartyAccountCredited) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary populates the message from its binary representation.
// For simplicity in this example we use JSON.
func (m *ThirdPartyAccountCredited) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

// MarshalBinary returns a binary representation of the message.
// For simplicity in this example we use JSON.
func (m *ThirdPartyAccountCreditFailed) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalBinary populates the message from its binary representation.
// For simplicity in this example we use JSON.
func (m *ThirdPartyAccountCreditFailed) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

// Validate returns a non-nil error if the message is invalid.
func (m *ThirdPartyAccountCredited) Validate(dogma.EventValidationScope) error {
	if m.TransactionID == "" {
		return errors.New("ThirdPartyAccountCredited must not have an empty transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("ThirdPartyAccountCredited must not have an empty account ID")
	}
	if m.Amount < 1 {
		return errors.New("ThirdPartyAccountCredited must have a positive amount")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m *ThirdPartyAccountCreditFailed) Validate(dogma.EventValidationScope) error {
	if m.TransactionID == "" {
		return errors.New("ThirdPartyAccountCreditFailed must not have an empty transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("ThirdPartyAccountCreditFailed must not have an empty account ID")
	}
	if m.Amount < 1 {
		return errors.New("ThirdPartyAccountCreditFailed must have a positive amount")
	}

	return nil
}
