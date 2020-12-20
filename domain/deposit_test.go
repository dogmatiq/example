package domain_test

import (
	"testing"

	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
	. "github.com/dogmatiq/testkit"
)

func Test_Deposit(t *testing.T) {
	t.Run(
		"when the deposit has not yet started",
		func(t *testing.T) {
			t.Run(
				"it deposits the funds into the account",
				func(t *testing.T) {
					Begin(t, &example.App{}).
						Prepare(
							ExecuteCommand(
								commands.OpenAccount{
									CustomerID:  "C001",
									AccountID:   "A001",
									AccountName: "Anna Smith",
								},
							),
						).
						Expect(
							ExecuteCommand(
								commands.Deposit{
									TransactionID: "T001",
									AccountID:     "A001",
									Amount:        500,
								},
							),
							ToRecordEvent(
								events.DepositApproved{
									TransactionID: "T001",
									AccountID:     "A001",
									Amount:        500,
								},
							),
						).
						// verify that funds are availalbe
						Expect(
							ExecuteCommand(
								commands.Withdraw{
									TransactionID: "W001",
									AccountID:     "A001",
									Amount:        100,
									ScheduledDate: "2001-02-03",
								},
							),
							ToRecordEvent(
								events.WithdrawalApproved{
									TransactionID: "W001",
									AccountID:     "A001",
									Amount:        100,
								},
							),
						)
				},
			)
		},
	)
}
