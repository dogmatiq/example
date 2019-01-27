package account_test

import (
	"fmt"
	"testing"

	"github.com/dogmatiq/dapper"
	"github.com/dogmatiq/dogmatest/engine"
	"github.com/dogmatiq/dogmatest/engine/fact"
	"github.com/dogmatiq/example"

	"github.com/dogmatiq/dogmatest"
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
			dogmatest.
				New(&example.App{}).
				Begin(
					t,
					engine.WithObserver(
						fact.ObserverFunc(func(f fact.Fact) {
							dapper.Print(f)
							fmt.Print("\n\n")
						}),
					),
				).
				ExecuteCommand(
					cmd,
					dogmatest.ExpectEvent(
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
			// Begin(t, engine).
			// 	Reset(cmd).
			// 	TestCommand(cmd).
			// 	Expect(
			// 		Not(
			// 			EventType(messages.AccountOpened{}),
			// 		),
			// 	)
		},
	)
}
