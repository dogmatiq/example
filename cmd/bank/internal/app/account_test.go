package app_test

import (
	"testing"

	"github.com/dogmatiq/examples/cmd/bank/internal/messages"
	. "github.com/dogmatiq/examples/dogmatest"
)

func TestAccount_OpenAccount(t *testing.T) {
	cmd := messages.OpenAccount{
		AccountID: "A001",
		Name:      "Bob Jones",
	}

	t.Run(
		"it opens the account",
		func(t *testing.T) {
			Begin(t, engine).
				TestCommand(cmd).
				Expect(
					Event(messages.AccountOpened{
						AccountID: "A001",
						Name:      "Bob Jones",
					}),
				)
		},
	)

	t.Run(
		"it does not open an account that is already open",
		func(t *testing.T) {
			Begin(t, engine).
				Reset(cmd).
				TestCommand(cmd).
				Expect(
					Not(
						EventType(messages.AccountOpened{}),
					),
				)
		},
	)
}