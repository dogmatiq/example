package domain_test

import (
	"testing"

	"github.com/dogmatiq/example/internal/testrunner"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
	. "github.com/dogmatiq/testkit/assert"
)

func Test_Customer(t *testing.T) {
	t.Run(
		"when a new customer opens an account",
		func(t *testing.T) {
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
								AccountName:   "Bob Jones",
							},
							EventRecorded(
								events.CustomerAcquired{
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
		},
	)

	t.Run(
		"when an existing customer opens an account",
		func(t *testing.T) {
			t.Run(
				"it does not reacquire a customer that has already been acquired",
				func(t *testing.T) {
					cmd := commands.OpenAccountForNewCustomer{
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
								EventTypeRecorded(events.CustomerAcquired{}),
							),
						)
				},
			)

			t.Run(
				"it does not reacquire a customer that has already been acquired",
				func(t *testing.T) {
					cmd := commands.OpenAccount{
						CustomerID:  "C001",
						AccountID:   "A002",
						AccountName: "Bob Jones and Co",
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
		},
	)

	t.Run(
		"when a customer changes email address",
		func(t *testing.T) {
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
								AccountName:   "Bob Jones",
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
				"it does not change the email address again if it is not different",
				func(t *testing.T) {
					testrunner.Runner.
						Begin(t).
						Prepare(
							commands.OpenAccountForNewCustomer{
								CustomerID:    "C001",
								CustomerName:  "Bob Jones",
								CustomerEmail: "bob@example.com",
								AccountID:     "A001",
								AccountName:   "Bob Jones",
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
		},
	)
}
