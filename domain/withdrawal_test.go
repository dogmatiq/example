package domain_test

import (
	"testing"

	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/messages"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
	. "github.com/dogmatiq/testkit"
)

func Test_Withdraw(t *testing.T) {
	t.Run(
		"when there are sufficient funds",
		func(t *testing.T) {
			t.Run(
				"it withdraws the funds from the account",
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
							ExecuteCommand(
								commands.Deposit{
									TransactionID: "T001",
									AccountID:     "A001",
									Amount:        500,
								},
							),
						).
						Expect(
							ExecuteCommand(
								commands.Withdraw{
									TransactionID: "T002",
									AccountID:     "A001",
									Amount:        500,
									ScheduledDate: "2001-02-03",
								},
							),
							ToRecordEvent(
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
								commands.Withdraw{
									TransactionID: "T001",
									AccountID:     "A001",
									Amount:        500,
									ScheduledDate: "2001-02-03",
								},
							),
							ToRecordEvent(
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
					Begin(t, &example.App{}).
						Prepare(
							ExecuteCommand(
								commands.OpenAccount{
									CustomerID:  "C001",
									AccountID:   "A001",
									AccountName: "Anna Smith",
								},
							),
							ExecuteCommand(
								commands.Deposit{
									TransactionID: "T001",
									AccountID:     "A001",
									Amount:        expectedDailyDebitLimit + 10000,
								},
							),
						).
						Expect(
							ExecuteCommand(
								commands.Withdraw{
									TransactionID: "T002",
									AccountID:     "A001",
									Amount:        500,
									ScheduledDate: "2001-02-03",
								},
							),
							ToRecordEvent(
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
					Begin(t, &example.App{}).
						Prepare(
							ExecuteCommand(
								commands.OpenAccount{
									CustomerID:  "C001",
									AccountID:   "A001",
									AccountName: "Anna Smith",
								},
							),
							ExecuteCommand(
								commands.OpenAccount{
									CustomerID:  "C002",
									AccountID:   "A002",
									AccountName: "Bob Jones",
								},
							),
							ExecuteCommand(
								commands.Deposit{
									TransactionID: "D001",
									AccountID:     "A001",
									Amount:        expectedDailyDebitLimit + 10000,
								},
							),
							ExecuteCommand(
								commands.Deposit{
									TransactionID: "D002",
									AccountID:     "A002",
									Amount:        expectedDailyDebitLimit + 10000,
								},
							),
						).
						Expect(
							ExecuteCommand(
								commands.Withdraw{
									TransactionID: "T001",
									AccountID:     "A001",
									Amount:        expectedDailyDebitLimit,
									ScheduledDate: "2001-02-03",
								},
							),
							ToRecordEvent(
								events.WithdrawalApproved{
									TransactionID: "T001",
									AccountID:     "A001",
									Amount:        expectedDailyDebitLimit,
								},
							),
						).
						Expect(
							ExecuteCommand(
								commands.Withdraw{
									TransactionID: "T002",
									AccountID:     "A002",
									Amount:        expectedDailyDebitLimit,
									ScheduledDate: "2001-02-03",
								},
							),
							ToRecordEvent(
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
					Begin(t, &example.App{}).
						Prepare(
							ExecuteCommand(
								commands.OpenAccount{
									CustomerID:  "C001",
									AccountID:   "A001",
									AccountName: "Anna Smith",
								},
							),
							ExecuteCommand(
								commands.Deposit{
									TransactionID: "D001",
									AccountID:     "A001",
									Amount:        expectedDailyDebitLimit * 2,
								},
							),
							ExecuteCommand(
								commands.Withdraw{
									TransactionID: "T001",
									AccountID:     "A001",
									Amount:        expectedDailyDebitLimit,
									ScheduledDate: "2001-02-03",
								},
							),
						).
						Expect(
							ExecuteCommand(
								commands.Withdraw{
									TransactionID: "T002",
									AccountID:     "A001",
									Amount:        500,
									ScheduledDate: "2001-02-04",
								},
							),
							ToRecordEvent(
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
					Begin(t, &example.App{}).
						Prepare(
							ExecuteCommand(
								commands.OpenAccount{
									CustomerID:  "C001",
									AccountID:   "A001",
									AccountName: "Anna Smith",
								},
							),
							ExecuteCommand(
								commands.Deposit{
									TransactionID: "D001",
									AccountID:     "A001",
									Amount:        expectedDailyDebitLimit + 10000,
								},
							),
						).
						Expect(
							ExecuteCommand(
								commands.Withdraw{
									TransactionID: "T001",
									AccountID:     "A001",
									Amount:        expectedDailyDebitLimit + 1,
									ScheduledDate: "2001-02-03",
								},
							),
							ToRecordEvent(
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
