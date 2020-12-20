package domain_test

import (
	"testing"

	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
	. "github.com/dogmatiq/testkit"
)

func Test_OpenAccountForNewCustomer(t *testing.T) {
	t.Run(
		"when the customer does not exist",
		func(t *testing.T) {
			t.Run(
				"it acquires the customer",
				func(t *testing.T) {
					Begin(t, &example.App{}).
						Expect(
							ExecuteCommand(
								commands.OpenAccountForNewCustomer{
									CustomerID:   "C001",
									CustomerName: "Bob Jones",
									AccountID:    "A001",
									AccountName:  "Bob Jones",
								},
							),
							ToRecordEvent(
								events.CustomerAcquired{
									CustomerID:   "C001",
									CustomerName: "Bob Jones",
									AccountID:    "A001",
									AccountName:  "Bob Jones",
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
						CustomerID:   "C001",
						CustomerName: "Bob Jones",
						AccountID:    "A001",
						AccountName:  "Bob Jones",
					}

					Begin(t, &example.App{}).
						Prepare(
							ExecuteCommand(cmd),
						).
						Expect(
							ExecuteCommand(cmd),
							NoneOf(
								ToRecordEventOfType(events.CustomerAcquired{}),
							),
						)
				},
			)
		},
	)
}
