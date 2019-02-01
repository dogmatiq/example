package account_test

import (
	"testing"

	. "github.com/dogmatiq/dogmatest/assert"
	"github.com/dogmatiq/example/internal/testrunner"
	"github.com/dogmatiq/example/messages"
)

func TestAccount_OpenAccount(t *testing.T) {
	t.Run(
		"it opens the account",
		func(t *testing.T) {
			testrunner.Runner.
				Begin(t).
				ExecuteCommand(
					messages.OpenAccount{
						AccountID: "A001",
						Name:      "Bob Jones",
					},
					EventRecorded(
						messages.AccountOpened{
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
			cmd := messages.OpenAccount{
				AccountID: "A001",
				Name:      "Bob Jones",
			}

			testrunner.Runner.
				Begin(t).
				Setup(cmd).
				ExecuteCommand(
					cmd,
					NoneOf(
						EventTypeRecorded(messages.AccountOpened{}),
					),
				)
		},
	)
}
