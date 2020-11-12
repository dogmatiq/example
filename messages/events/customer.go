package events

import (
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
		"acquired customer %s %s,
		m.CustomerID,
		m.CustomerName,
	)
}
