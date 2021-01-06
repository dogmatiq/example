package events

import (
	"errors"
	"fmt"
)

// CustomerAcquired is an event indicating that a new customer has been
// acquired.
type CustomerAcquired struct {
	CustomerID   string
	CustomerName string
	AccountID    string
	AccountName  string
}

// MessageDescription returns a human-readable description of the message.
func (m CustomerAcquired) MessageDescription() string {
	return fmt.Sprintf(
		"acquired customer %s %s",
		m.CustomerID,
		m.CustomerName,
	)
}

// Validate returns a non-nil error if the message is invalid.
func (m CustomerAcquired) Validate() error {
	if m.CustomerID == "" {
		return errors.New("CustomerAcquired needs a valid customer ID")
	}
	if m.CustomerName == "" {
		return errors.New("CustomerAcquired needs a valid name")
	}
	if m.AccountID == "" {
		return errors.New("CustomerAcquired needs a valid account ID")
	}
	if m.AccountName == "" {
		return errors.New("CustomerAcquired needs a valid account name")
	}

	return nil
}
