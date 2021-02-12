package commands

import (
	"errors"
	"fmt"
	"time"

	"github.com/dogmatiq/example/messages"
)

// OpenAccountForNewCustomer is a command requesting that a new bank account be
// opened for a new customer.
type OpenAccountForNewCustomer struct {
	CustomerID   string
	CustomerName string
	AccountID    string
	AccountName  string
}

// OpenAccount is a command requesting that a new bank account be opened for an
// existing customer.
type OpenAccount struct {
	CustomerID  string
	AccountID   string
	AccountName string
}

// CreditAccount is a command that requests a bank account be credited.
type CreditAccount struct {
	TransactionID   string
	AccountID       string
	TransactionType messages.TransactionType
	Amount          int64
}

// DebitAccount is a command that requests a bank account be debited.
type DebitAccount struct {
	TransactionID   string
	AccountID       string
	TransactionType messages.TransactionType
	Amount          int64
	ScheduledTime   time.Time
}

// MessageDescription returns a human-readable description of the message.
func (m OpenAccountForNewCustomer) MessageDescription() string {
	return fmt.Sprintf(
		"customer %s %s is opening their first account %s %s",
		m.CustomerID,
		m.CustomerName,
		m.AccountID,
		m.AccountName,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m OpenAccount) MessageDescription() string {
	return fmt.Sprintf(
		"opening account %s %s for customer %s",
		m.AccountID,
		m.AccountName,
		m.CustomerID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m CreditAccount) MessageDescription() string {
	return fmt.Sprintf(
		"%s %s: crediting %s to account %s",
		m.TransactionType,
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.AccountID,
	)
}

// MessageDescription returns a human-readable description of the message.
func (m DebitAccount) MessageDescription() string {
	return fmt.Sprintf(
		"%s %s: debiting %s from account %s",
		m.TransactionType,
		m.TransactionID,
		messages.FormatAmount(m.Amount),
		m.AccountID,
	)
}

// Validate returns a non-nil error if the message is invalid.
func (m OpenAccountForNewCustomer) Validate() error {
	if m.CustomerID == "" {
		return errors.New("OpenAccountForNewCustomer must not have an empty customer ID")
	}
	if m.CustomerName == "" {
		return errors.New("OpenAccountForNewCustomer must not have an empty customer name")
	}
	if m.AccountID == "" {
		return errors.New("OpenAccountForNewCustomer must not have an empty account ID")
	}
	if m.AccountName == "" {
		return errors.New("OpenAccountForNewCustomer must not have an empty account name")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m OpenAccount) Validate() error {
	if m.CustomerID == "" {
		return errors.New("OpenAccount must not have an empty customer ID")
	}
	if m.AccountID == "" {
		return errors.New("OpenAccount must not have an empty account ID")
	}
	if m.AccountName == "" {
		return errors.New("OpenAccount must not have an empty account name")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m CreditAccount) Validate() error {
	if m.TransactionID == "" {
		return errors.New("CreditAccount must not have an empty transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("CreditAccount must not have an empty account ID")
	}
	if err := m.TransactionType.Validate(); err != nil {
		return fmt.Errorf("CreditAccount must have a valid transaction type: %w", err)
	}
	if m.Amount < 1 {
		return errors.New("CreditAccount must have a positive amount")
	}

	return nil
}

// Validate returns a non-nil error if the message is invalid.
func (m DebitAccount) Validate() error {
	if m.TransactionID == "" {
		return errors.New("DebitAccount must not have an empty transaction ID")
	}
	if m.AccountID == "" {
		return errors.New("DebitAccount must not have an empty account ID")
	}
	if err := m.TransactionType.Validate(); err != nil {
		return fmt.Errorf("DebitAccount must have a valid transaction type: %w", err)
	}
	if m.Amount < 1 {
		return errors.New("DebitAccount must have a positive amount")
	}

	return nil
}
