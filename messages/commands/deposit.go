package commands

import (
	"errors"
	"fmt"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/messages"
)

// Deposit is a command requesting that funds be deposited into a bank account.
type Deposit struct {
	TransactionID string
	AccountID     string
	Amount        int64
}

// ApproveDeposit is a command that approves an account deposit.
type ApproveDeposit struct {
	TransactionID string
	AccountID     string
	Amount        int64
}

// MessageDescription returns a human-readable description of the message.
func (m Deposit) MessageDescription() string {
	return fmt.Sprintf(
		"deposit %s: depositing %s into account %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.AccountID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m ApproveDeposit) MessageDescription() string {
	return fmt.Sprintf(
		"deposit %s: approving deposit of %s into account %s",
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.AccountID,
	)
}

// Validate returns a non-nil error if the message is invalid.
func (m Deposit) Validate(dogma.CommandValidationScope) error {
	if m.TransactionID == "" {
		return errors.New("Deposit must not have an empty transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("Deposit must not have an empty account ID")
	}
	if m.Amount < 1 {
		return errors.New("Deposit must have a positive amount")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m ApproveDeposit) Validate(dogma.CommandValidationScope) error {
	if m.TransactionID == "" {
		return errors.New("ApproveDeposit must not have an empty transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("ApproveDeposit must not have an empty account ID")
	}
	if m.Amount < 1 {
		return errors.New("ApproveDeposit must have a positive amount")
	}

	return nil
}
