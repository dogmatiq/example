package domain_test

import (
	"testing"

	"github.com/dogmatiq/example/internal/testrunner"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
	. "github.com/dogmatiq/testkit/assert"
)

func Test_Deposit(t *testing.T) {
	t.Run(
		"it deposits the funds into the account",
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
					commands.Deposit{
						TransactionID: "T001",
						AccountID:     "A001",
						Amount:        500,
					},
					EventRecorded(
						events.DepositApproved{
							TransactionID: "T001",
							AccountID:     "A001",
							Amount:        500,
						},
					),
				).
				// verify that funds are availalbe
				ExecuteCommand(
					commands.Withdraw{
						TransactionID: "W001",
						AccountID:     "A001",
						Amount:        100,
						ScheduledDate: "2001-02-03",
					},
					EventRecorded(
						events.WithdrawalApproved{
							TransactionID: "W001",
							AccountID:     "A001",
							Amount:        100,
						},
					),
				)
		},
	)
}
