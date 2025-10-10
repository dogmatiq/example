package domain_test

import (
	"testing"

	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
	. "github.com/dogmatiq/testkit"
)

func Test_OpenAccount(t *testing.T) {
	t.Run(
		"when the account does not exist",
		func(t *testing.T) {
			t.Run(
				"the new account is opened",
				func(t *testing.T) {
					Begin(t, &example.App{}).
						Expect(
							ExecuteCommand(
								&commands.OpenAccount{
									CustomerID:  "C001",
									AccountID:   "A001",
									AccountName: "Anna Smith",
								},
							),
							ToRecordEvent(
								&events.AccountOpened{
									CustomerID:  "C001",
									AccountID:   "A001",
									AccountName: "Anna Smith",
								},
							),
						)
				},
			)
		},
	)

	t.Run(
		"when the account already exists",
		func(t *testing.T) {
			t.Run(
				"nothing happens to the existing account",
				func(t *testing.T) {
					Begin(t, &example.App{}).
						Prepare(
							ExecuteCommand(
								&commands.OpenAccount{
									CustomerID:  "C001",
									AccountID:   "A001",
									AccountName: "Anna Smith",
								},
							),
						).
						Expect(
							ExecuteCommand(
								&commands.OpenAccount{
									CustomerID:  "C001",
									AccountID:   "A001",
									AccountName: "Anna Smith",
								},
							),
							NoneOf(
								ToRecordEventOfType(&events.AccountOpened{}),
							),
						)
				},
			)
		},
	)
}
