package customer_test

import (
	"testing"

	"github.com/dogmatiq/example/internal/testrunner"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
	. "github.com/dogmatiq/testkit/assert"
)

func TestCustomer_OpenAccountForNewCustomer(t *testing.T) {
	t.Run(
		"it acquires the customer",
		func(t *testing.T) {
			testrunner.Runner.
				Begin(t).
				ExecuteCommand(
					commands.OpenAccountForNewCustomer{
						CustomerID:    "C001",
						CustomerName:  "Bob Jones",
						CustomerEmail: "bob@example.com",
						AccountID:     "A001",
						AccountName:   "Savings",
					},
					EventRecorded(
						events.CustomerAcquired{
							CustomerID:    "C001",
							CustomerName:  "Bob Jones",
							CustomerEmail: "bob@example.com",
							AccountID:     "A001",
							AccountName:   "Savings",
						},
					),
				)
		},
	)

	t.Run(
		"it does not acquire a customer that already exists",
		func(t *testing.T) {
			cmd := commands.OpenAccountForNewCustomer{
				CustomerID:    "C001",
				CustomerName:  "Bob Jones",
				CustomerEmail: "bob@example.com",
				AccountID:     "A001",
				AccountName:   "Savings",
			}

			testrunner.Runner.
				Begin(t).
				Prepare(cmd).
				ExecuteCommand(
					cmd,
					NoneOf(
						EventTypeRecorded(events.CustomerAcquired{}),
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
					commands.OpenAccountForNewCustomer{
						CustomerID:    "C001",
						CustomerName:  "Bob Jones",
						CustomerEmail: "bob@example.com",
						AccountID:     "A001",
						AccountName:   "Savings",
					},
				).
				ExecuteCommand(
					commands.ChangeCustomerEmailAddress{
						CustomerID:    "C001",
						CustomerEmail: "newbob@example.com",
					},
					EventRecorded(
						events.CustomerEmailAddressChanged{
							CustomerID:    "C001",
							CustomerEmail: "newbob@example.com",
						},
					),
				)
		},
	)

	t.Run(
		"it does not change the email address if it is the same",
		func(t *testing.T) {
			testrunner.Runner.
				Begin(t).
				Prepare(
					commands.OpenAccountForNewCustomer{
						CustomerID:    "C001",
						CustomerName:  "Bob Jones",
						CustomerEmail: "bob@example.com",
						AccountID:     "A001",
						AccountName:   "Savings",
					},
				).
				ExecuteCommand(
					commands.ChangeCustomerEmailAddress{
						CustomerID:    "C001",
						CustomerEmail: "bob@example.com",
					},
					NoneOf(
						EventTypeRecorded(events.CustomerEmailAddressChanged{}),
					),
				)
		},
	)
}
