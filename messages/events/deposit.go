package events

import (
	"errors"
	"fmt"

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
func (m DepositStarted) Validate() error {
	if m.TransactionID == "" {
		return errors.New("DepositStarted needs a valid transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("DepositStarted needs a valid account ID")
	}
	if m.Amount < 1 {
		return errors.New("DepositStarted needs a valid amount")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m DepositApproved) Validate() error {
	if m.TransactionID == "" {
		return errors.New("DepositApproved needs a valid transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("DepositApproved needs a valid account ID")
	}
	if m.Amount < 1 {
		return errors.New("DepositApproved needs a valid amount")
	}

	return nil
}
