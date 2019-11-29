package domain_test

import (
	"testing"

	"github.com/dogmatiq/example/internal/testrunner"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
	. "github.com/dogmatiq/testkit/assert"
)

func Test_OpenAccount(t *testing.T) {
	t.Run(
		"when the account does not exist",
		func(t *testing.T) {
			t.Run(
				"the new account is opened",
				func(t *testing.T) {
					testrunner.New(nil).
						Begin(t).
						ExecuteCommand(
							commands.OpenAccount{
								CustomerID:  "C001",
								AccountID:   "A001",
								AccountName: "Anna Smith",
							},
							EventRecorded(
								events.AccountOpened{
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
					testrunner.New(nil).
						Begin(t).
						Prepare(
							commands.OpenAccount{
								CustomerID:  "C001",
								AccountID:   "A001",
								AccountName: "Anna Smith",
							}).
						ExecuteCommand(
							commands.OpenAccount{
								CustomerID:  "C001",
								AccountID:   "A001",
								AccountName: "Anna Smith",
							},
							NoneOf(
								EventTypeRecorded(events.AccountOpened{}),
							),
						)
				},
			)
		},
	)
}
