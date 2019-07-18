package domain_test

import (
	"testing"

	"github.com/dogmatiq/example/internal/testrunner"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
	. "github.com/dogmatiq/testkit/assert"
)

func Test_OpenAccountForNewCustomer(t *testing.T) {
	t.Run(
		"when the customer does not exist",
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
		"when the customer already exists",
		func(t *testing.T) {
			t.Run(
				"it does not reacquire the customer",
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
		},
	)
}

func Test_ChangeCustomerEmailAddress(t *testing.T) {
	t.Run(
		"when the email address is different",
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
		},
	)

	t.Run(
		"when the email address is the same",
		func(t *testing.T) {
			t.Run(
				"it does not change the email address",
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
