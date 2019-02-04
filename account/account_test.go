package account_test

import (
	"testing"

	"github.com/dogmatiq/example/internal/testrunner"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
	. "github.com/dogmatiq/testkit/assert"
)

func TestAccount_OpenAccount(t *testing.T) {
	t.Run(
		"it opens the account",
		func(t *testing.T) {
			testrunner.Runner.
				Begin(t).
				ExecuteCommand(
					commands.OpenAccount{
						AccountID: "A001",
						Name:      "Bob Jones",
					},
					EventRecorded(
						events.AccountOpened{
							AccountID: "A001",
							Name:      "Bob Jones",
						},
					),
				)
		},
	)

	t.Run(
		"it does not open an account that is already open",
		func(t *testing.T) {
			cmd := commands.OpenAccount{
				AccountID: "A001",
				Name:      "Bob Jones",
			}

			testrunner.Runner.
				Begin(t).
				Prepare(cmd).
				ExecuteCommand(
					cmd,
					NoneOf(
						EventTypeRecorded(events.AccountOpened{}),
					),
				)
		},
	)
}
