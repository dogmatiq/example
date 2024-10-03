package events

import (
	"errors"
	"fmt"

	"github.com/dogmatiq/dogma"
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
func (m CustomerAcquired) Validate(dogma.EventValidationScope) error {
	if m.CustomerID == "" {
		return errors.New("CustomerAcquired must not have an empty customer ID")
	}
	if m.CustomerName == "" {
		return errors.New("CustomerAcquired must not have an empty name")
	}
	if m.AccountID == "" {
		return errors.New("CustomerAcquired must not have an empty account ID")
	}
	if m.AccountName == "" {
		return errors.New("CustomerAcquired must not have an empty account name")
	}

	return nil
}
