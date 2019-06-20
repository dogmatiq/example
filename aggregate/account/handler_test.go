package account_test

import (
	"testing"

	"github.com/dogmatiq/example/internal/testrunner"
	"github.com/dogmatiq/example/message/command"
	"github.com/dogmatiq/example/message/event"
	. "github.com/dogmatiq/testkit/assert"
)

func TestAccount_OpenAccount(t *testing.T) {
	t.Run(
		"it opens the account",
		func(t *testing.T) {
			testrunner.Runner.
				Begin(t).
				ExecuteCommand(
					command.OpenAccount{
						CustomerID:  "C001",
						AccountID:   "A001",
						AccountName: "Anna Smith",
					},
					EventRecorded(
						event.AccountOpened{
							CustomerID:  "C001",
							AccountID:   "A001",
							AccountName: "Anna Smith",
						},
					),
				)
		},
	)

	t.Run(
		"it does not open an account that is already open",
		func(t *testing.T) {

			testrunner.Runner.
				Begin(t).
				Prepare(
					command.OpenAccount{
						CustomerID:  "C001",
						AccountID:   "A001",
						AccountName: "Anna Smith",
					}).
				ExecuteCommand(
					command.OpenAccount{
						CustomerID:  "C001",
						AccountID:   "A001",
						AccountName: "Anna Smith",
					},
					NoneOf(
						EventTypeRecorded(event.AccountOpened{}),
					),
				)
		},
	)
}
