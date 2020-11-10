package events

import (
	"fmt"

	"github.com/dogmatiq/example/messages"
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
func (m *CustomerAcquired) MessageDescription() string {
	return fmt.Sprintf(
		"acquired customer %s %s with first account %s %s",
		messages.FormatID(m.CustomerID),
		m.CustomerName,
		messages.FormatID(m.AccountID),
		m.AccountName,
	)
}
