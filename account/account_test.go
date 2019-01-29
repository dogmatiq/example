package account_test

import (
	"testing"

	"github.com/dogmatiq/example/internal/testrunner"

	. "github.com/dogmatiq/dogmatest"
	"github.com/dogmatiq/example/messages"
)

func TestAccount_OpenAccount(t *testing.T) {
	cmd := messages.OpenAccount{
		AccountID: "A001",
		Name:      "Bob Jones",
	}

	t.Run(
		"it opens the account",
		func(t *testing.T) {
			testrunner.
				Runner.
				Begin(t).
				ExecuteCommand(
					cmd,
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
			testrunner.
				Runner.
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
