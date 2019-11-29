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
					testrunner.New(nil).
						Begin(t).
						ExecuteCommand(
							commands.OpenAccountForNewCustomer{
								CustomerID:   "C001",
								CustomerName: "Bob Jones",
								AccountID:    "A001",
								AccountName:  "Bob Jones",
							},
							EventRecorded(
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

					testrunner.New(nil).
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
