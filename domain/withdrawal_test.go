package domain_test

import (
	"testing"

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
								Amount:        500,
							},
						).
						ExecuteCommand(
							commands.Withdraw{
								TransactionID: "T002",
								AccountID:     "A001",
								Amount:        500,
								ScheduledDate: businessDateToday,
							},
							EventRecorded(
								events.WithdrawalApproved{
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
								ScheduledDate: businessDateToday,
							},
							EventRecorded(
								events.WithdrawalDeclined{
									TransactionID: "T001",
									AccountID:     "A001",
									Amount:        500,
									Reason:        messages.InsufficientFunds,
								},
							),
						)
				},
			)
		},
	)

	t.Run(
		"when within daily debit limit",
		func(t *testing.T) {
			t.Run(
				"it withdraws funds from the specified account",
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
								ScheduledDate: businessDateToday,
							},
							EventRecorded(
								events.WithdrawalApproved{
									TransactionID: "T002",
									AccountID:     "A001",
									Amount:        500,
								},
							),
						)
				},
			)

			t.Run(
				"it applies the limit per account",
				func(t *testing.T) {
					testrunner.Runner.
						Begin(t).
						Prepare(
							commands.OpenAccount{
								CustomerID:  "C001",
								AccountID:   "A001",
								AccountName: "Anna Smith",
							},
							commands.OpenAccount{
								CustomerID:  "C002",
								AccountID:   "A002",
								AccountName: "Bob Jones",
							},
							commands.Deposit{
								TransactionID: "D001",
								AccountID:     "A001",
								Amount:        expectedDailyDebitLimit + 10000,
							},
							commands.Deposit{
								TransactionID: "D002",
								AccountID:     "A002",
								Amount:        expectedDailyDebitLimit + 10000,
							},
						).
						ExecuteCommand(
							commands.Withdraw{
								TransactionID: "T001",
								AccountID:     "A001",
								Amount:        expectedDailyDebitLimit,
								ScheduledDate: businessDateToday,
							},
							EventRecorded(
								events.WithdrawalApproved{
									TransactionID: "T001",
									AccountID:     "A001",
									Amount:        expectedDailyDebitLimit,
								},
							),
						).
						ExecuteCommand(
							commands.Withdraw{
								TransactionID: "T002",
								AccountID:     "A002",
								Amount:        expectedDailyDebitLimit,
								ScheduledDate: businessDateToday,
							},
							EventRecorded(
								events.WithdrawalApproved{
									TransactionID: "T002",
									AccountID:     "A002",
									Amount:        expectedDailyDebitLimit,
								},
							),
						)
				},
			)

			t.Run(
				"it applies the limit per day",
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
								TransactionID: "D001",
								AccountID:     "A001",
								Amount:        expectedDailyDebitLimit * 2,
							},
							commands.Withdraw{
								TransactionID: "T001",
								AccountID:     "A001",
								Amount:        expectedDailyDebitLimit,
								ScheduledDate: businessDateToday,
							},
						).
						ExecuteCommand(
							commands.Withdraw{
								TransactionID: "T002",
								AccountID:     "A001",
								Amount:        500,
								ScheduledDate: businessDateTomorrow,
							},
							EventRecorded(
								events.WithdrawalApproved{
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
		"when daily debit limit will be exceeded",
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
								TransactionID: "D001",
								AccountID:     "A001",
								Amount:        expectedDailyDebitLimit + 10000,
							},
						).
						ExecuteCommand(
							commands.Withdraw{
								TransactionID: "T001",
								AccountID:     "A001",
								Amount:        expectedDailyDebitLimit + 1,
								ScheduledDate: businessDateToday,
							},
							EventRecorded(
								events.WithdrawalDeclined{
									TransactionID: "T001",
									AccountID:     "A001",
									Amount:        expectedDailyDebitLimit + 1,
									Reason:        messages.DailyDebitLimitExceeded,
								},
							),
						)
				},
			)
		},
	)
}
