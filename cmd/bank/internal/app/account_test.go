package app_test

import (
	"context"
	"testing"

	. "github.com/dogmatiq/examples/cmd/bank/internal/app"
	"github.com/dogmatiq/examples/cmd/bank/internal/messages"
	"github.com/dogmatiq/examples/dogmatest"
)

func TestAccount_OpenAccount(t *testing.T) {
	engine := dogmatest.New(App)

	cmd := messages.OpenAccount{
		AccountID: "A001",
		Name:      "Bob Jones",
	}

	t.Run(
		"it opens the account",
		func(t *testing.T) {
			engine.
				TestCommand(cmd).
				ExpectEvents(
					messages.AccountOpened{
						AccountID: "A001",
						Name:      "Bob Jones",
					},
				)
		},
	)

	t.Run(
		"it does not opens an account that is already open",
		func(t *testing.T) {
			engine.ExecuteCommand(context.Background(), cmd)
			engine.
				TestCommand(cmd).
				ExpectNoEvents()
		},
	)
}
