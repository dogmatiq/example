package events

import (
	"errors"
	"fmt"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages"
)

// DepositStarted is an event indicating that the process of depositing funds
// into an account has begun.
type DepositStarted struct {
	TransactionID string
	AccountID     string
	Amount        int64
}

// DepositApproved is an event that indicates a requested deposit has been
// approved.
type DepositApproved struct {
	TransactionID string
	AccountID     string
	Amount        int64
}

// MessageDescription returns a human-readable description of the message.
func (m DepositStarted) MessageDescription() string {
	return fmt.Sprintf(
		"deposit %s: started deposit of %s into account %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.AccountID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m DepositApproved) MessageDescription() string {
	return fmt.Sprintf(
		"deposit %s: approved deposit of %s into account %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.AccountID,
	)
}

// Validate returns a non-nil error if the message is invalid.
func (m DepositStarted) Validate(dogma.EventValidationScope) error {
	if m.TransactionID == "" {
		return errors.New("DepositStarted must not have an empty transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("DepositStarted must not have an empty account ID")
	}
	if m.Amount < 1 {
		return errors.New("DepositStarted must have a positive amount")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m DepositApproved) Validate(dogma.EventValidationScope) error {
	if m.TransactionID == "" {
		return errors.New("DepositApproved must not have an empty transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("DepositApproved must not have an empty account ID")
	}
	if m.Amount < 1 {
		return errors.New("DepositApproved must have a positive amount")
	}

	return nil
}
