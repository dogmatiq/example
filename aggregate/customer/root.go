package customer

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/message/event"
)

// root is the aggregate root for a bank customer.
type root struct {
	// Email is the customer email address.
	Email string
}

func (r *root) ApplyEvent(m dogma.Message) {
	switch x := m.(type) {
	case event.CustomerAcquired:
		r.Email = x.CustomerEmail
	case event.CustomerEmailAddressChanged:
		r.Email = x.CustomerEmail
	}
}
