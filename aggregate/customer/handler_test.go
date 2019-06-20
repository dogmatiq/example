package customer_test

import (
	"testing"

	"github.com/dogmatiq/example/internal/testrunner"
	"github.com/dogmatiq/example/message/command"
	"github.com/dogmatiq/example/message/event"
	. "github.com/dogmatiq/testkit/assert"
)

func TestCustomer_OpenAccountForNewCustomer(t *testing.T) {
	t.Run(
		"it acquires the customer",
		func(t *testing.T) {
			testrunner.Runner.
				Begin(t).
				ExecuteCommand(
					command.OpenAccountForNewCustomer{
						CustomerID:    "C001",
						CustomerName:  "Bob Jones",
						CustomerEmail: "bob@example.com",
						AccountID:     "A001",
						AccountName:   "Bob Jones",
					},
					EventRecorded(
						event.CustomerAcquired{
							CustomerID:    "C001",
							CustomerName:  "Bob Jones",
							CustomerEmail: "bob@example.com",
							AccountID:     "A001",
							AccountName:   "Bob Jones",
						},
					),
				)
		},
	)

	t.Run(
		"it does not reacquire a customer that has already been acquired",
		func(t *testing.T) {
			cmd := command.OpenAccountForNewCustomer{
				CustomerID:    "C001",
				CustomerName:  "Bob Jones",
				CustomerEmail: "bob@example.com",
				AccountID:     "A001",
				AccountName:   "Bob Jones",
			}

			testrunner.Runner.
				Begin(t).
				Prepare(cmd).
				ExecuteCommand(
					cmd,
					NoneOf(
						EventTypeRecorded(event.CustomerAcquired{}),
					),
				)
		},
	)
}

func TestCustomer_ChangeCustomerEmailAddress(t *testing.T) {
	t.Run(
		"it changes the email address of the customer",
		func(t *testing.T) {
			testrunner.Runner.
				Begin(t).
				Prepare(
					command.OpenAccountForNewCustomer{
						CustomerID:    "C001",
						CustomerName:  "Bob Jones",
						CustomerEmail: "bob@example.com",
						AccountID:     "A001",
						AccountName:   "Bob Jones",
					},
				).
				ExecuteCommand(
					command.ChangeCustomerEmailAddress{
						CustomerID:    "C001",
						CustomerEmail: "newbob@example.com",
					},
					EventRecorded(
						event.CustomerEmailAddressChanged{
							CustomerID:    "C001",
							CustomerEmail: "newbob@example.com",
						},
					),
				)
		},
	)

	t.Run(
		"it does not change the email address again if it has not changed",
		func(t *testing.T) {
			testrunner.Runner.
				Begin(t).
				Prepare(
					command.OpenAccountForNewCustomer{
						CustomerID:    "C001",
						CustomerName:  "Bob Jones",
						CustomerEmail: "bob@example.com",
						AccountID:     "A001",
						AccountName:   "Bob Jones",
					},
				).
				ExecuteCommand(
					command.ChangeCustomerEmailAddress{
						CustomerID:    "C001",
						CustomerEmail: "bob@example.com",
					},
					NoneOf(
						EventTypeRecorded(event.CustomerEmailAddressChanged{}),
					),
				)
		},
	)
}
