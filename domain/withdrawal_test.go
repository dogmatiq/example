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
		"when there are sufficient funds",
		func(t *testing.T) {
			t.Run(
				"it withdraws the funds from the account",
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
								ScheduledDate: "2001-02-03",
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
		"when there are insufficient funds",
		func(t *testing.T) {
			t.Run(
				"it does not withdraw funds from the account",
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
								ScheduledDate: "2001-02-03",
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
		"when the withdrawal does not exceed the daily debit limit",
		func(t *testing.T) {
			t.Run(
				"it withdraws the funds from the account",
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
								ScheduledDate: "2001-02-03",
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
				"it enforces the daily debit limit per account",
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
								ScheduledDate: "2001-02-03",
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
								ScheduledDate: "2001-02-03",
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
				"it enforces the daily debit limit per day",
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
								ScheduledDate: "2001-02-03",
							},
						).
						ExecuteCommand(
							commands.Withdraw{
								TransactionID: "T002",
								AccountID:     "A001",
								Amount:        500,
								ScheduledDate: "2001-02-04",
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
		"when the withdrawal exceeds the daily debit limit",
		func(t *testing.T) {
			t.Run(
				"it does not withdraw any funds from the account",
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
								ScheduledDate: "2001-02-03",
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
