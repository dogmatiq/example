package domain_test

import (
	"testing"
	"time"

	"github.com/dogmatiq/example/internal/testrunner"
	"github.com/dogmatiq/example/messages"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
	. "github.com/dogmatiq/testkit/assert"
)

func Test_Withdraw(t *testing.T) {
	t.Run(
		"when sufficient funds",
		func(t *testing.T) {
			t.Run(
				"it withdraws some funds from an account",
				func(t *testing.T) {
					testrunner.Runner.
						Begin(t).
						Prepare(
							commands.OpenAccount{
								CustomerID:  "C001",
								AccountID:   "A001",
								AccountName: "Anna Smith",
							},
							commands.Deposit{
								TransactionID: "T001",
								AccountID:     "A001",
								Amount:        5000,
							},
						).
						ExecuteCommand(
							commands.Withdraw{
								TransactionID: "T002",
								AccountID:     "A001",
								Amount:        500,
								ScheduledDate: time.Unix(12345, 0),
							},
							EventRecorded(
								events.AccountDebitedForWithdrawal{
									TransactionID: "T002",
									AccountID:     "A001",
									Amount:        500,
								},
							),
						)
				},
			)
		},
	)

	t.Run(
		"when insufficient funds",
		func(t *testing.T) {
			t.Run(
				"it does not withdraw funds from an account",
				func(t *testing.T) {
					testrunner.Runner.
						Begin(t).
						Prepare(
							commands.OpenAccount{
								CustomerID:  "C001",
								AccountID:   "A001",
								AccountName: "Anna Smith",
							},
						).
						ExecuteCommand(
							commands.Withdraw{
								TransactionID: "T001",
								AccountID:     "A001",
								Amount:        500,
								ScheduledDate: time.Unix(12345, 0),
							},
							EventRecorded(
								events.WithdrawalDeclined{
									TransactionID: "T001",
									AccountID:     "A001",
									Amount:        500,
									Reason:        messages.ReasonInsufficientFunds,
								},
							),
						)
				},
			)
		},
	)

	// The current expected daily debit limit.
	const expectedDailyDebitLimit = 900000

	t.Run(
		"when within daily limit",
		func(t *testing.T) {
			t.Run(
				"it withdraws funds from an account",
				func(t *testing.T) {
					testrunner.Runner.
						Begin(t).
						Prepare(
							commands.OpenAccount{
								CustomerID:  "C001",
								AccountID:   "A001",
								AccountName: "Anna Smith",
							},
							commands.Deposit{
								TransactionID: "T001",
								AccountID:     "A001",
								Amount:        expectedDailyDebitLimit + 10000,
							},
						).
						ExecuteCommand(
							commands.Withdraw{
								TransactionID: "T002",
								AccountID:     "A001",
								Amount:        500,
								ScheduledDate: time.Unix(12345, 0),
							},
							EventRecorded(
								events.AccountDebitedForWithdrawal{
									TransactionID: "T002",
									AccountID:     "A001",
									Amount:        500,
								},
							),
						)
				},
			)
		},
	)

	t.Run(
		"when daily limit will be exceeded",
		func(t *testing.T) {
			t.Run(
				"it does not withdraw funds from an account",
				func(t *testing.T) {
					testrunner.Runner.
						Begin(t).
						Prepare(
							commands.OpenAccount{
								CustomerID:  "C001",
								AccountID:   "A001",
								AccountName: "Anna Smith",
							},
							commands.Deposit{
								TransactionID: "T001",
								AccountID:     "A001",
								Amount:        expectedDailyDebitLimit + 10000,
							},
						).
						ExecuteCommand(
							commands.Withdraw{
								TransactionID: "T002",
								AccountID:     "A001",
								Amount:        expectedDailyDebitLimit + 1,
								ScheduledDate: time.Unix(12345, 0),
							},
							EventRecorded(
								events.WithdrawalDeclined{
									TransactionID: "T002",
									AccountID:     "A001",
									Amount:        expectedDailyDebitLimit + 1,
									Reason:        messages.ReasonDailyDebitLimitExceeded,
								},
							),
						)
				},
			)
		},
	)
}
