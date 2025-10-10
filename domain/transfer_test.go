package domain_test

import (
	"testing"
	"time"

	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/messages"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
	. "github.com/dogmatiq/testkit"
)

func Test_Transfer(t *testing.T) {
	t.Run(
		"when there are sufficient funds",
		func(t *testing.T) {
			t.Run(
				"it transfers the funds from one account to another",
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
							ExecuteCommand(
								&commands.OpenAccount{
									CustomerID:  "C002",
									AccountID:   "A002",
									AccountName: "Bob Jones",
								},
							),
							ExecuteCommand(
								&commands.Deposit{
									TransactionID: "D001",
									AccountID:     "A001",
									Amount:        500,
								},
							),
						).
						Expect(
							ExecuteCommand(
								&commands.Transfer{
									TransactionID: "T001",
									FromAccountID: "A001",
									ToAccountID:   "A002",
									Amount:        100,
									ScheduledTime: time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC),
								},
							),
							ToRecordEvent(
								&events.TransferApproved{
									TransactionID: "T001",
									FromAccountID: "A001",
									ToAccountID:   "A002",
									Amount:        100,
								},
							),
						).
						// verify that funds are availalbe
						Expect(
							ExecuteCommand(
								&commands.Withdraw{
									TransactionID: "W001",
									AccountID:     "A002",
									Amount:        100,
									ScheduledTime: time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC),
								},
							),
							ToRecordEvent(
								&events.WithdrawalApproved{
									TransactionID: "W001",
									AccountID:     "A002",
									Amount:        100,
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
				"it does not transfer any funds from the account",
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
							ExecuteCommand(
								&commands.OpenAccount{
									CustomerID:  "C002",
									AccountID:   "A002",
									AccountName: "Bob Jones",
								},
							),
							ExecuteCommand(
								&commands.Deposit{
									TransactionID: "D001",
									AccountID:     "A001",
									Amount:        500,
								},
							),
						).
						Expect(
							ExecuteCommand(
								&commands.Transfer{
									TransactionID: "T001",
									FromAccountID: "A001",
									ToAccountID:   "A002",
									Amount:        1000,
									ScheduledTime: time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC),
								},
							),
							ToRecordEvent(
								&events.TransferDeclined{
									TransactionID: "T001",
									FromAccountID: "A001",
									ToAccountID:   "A002",
									Amount:        1000,
									Reason:        messages.InsufficientFunds,
								},
							),
						).
						// verify that funds are not availalbe
						Expect(
							ExecuteCommand(
								&commands.Withdraw{
									TransactionID: "W001",
									AccountID:     "A002",
									Amount:        100,
									ScheduledTime: time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC),
								},
							),
							ToRecordEvent(
								&events.WithdrawalDeclined{
									TransactionID: "W001",
									AccountID:     "A002",
									Amount:        100,
									Reason:        messages.InsufficientFunds,
								},
							),
						)
				},
			)
		},
	)

	t.Run(
		"when the transfer does not exceed the daily debit limit",
		func(t *testing.T) {
			t.Run(
				"it transfers the funds from one account to another",
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
							ExecuteCommand(
								&commands.OpenAccount{
									CustomerID:  "C002",
									AccountID:   "A002",
									AccountName: "Bob Jones",
								},
							),
							ExecuteCommand(
								&commands.Deposit{
									TransactionID: "T001",
									AccountID:     "A001",
									Amount:        expectedDailyDebitLimit + 10000,
								},
							),
						).
						Expect(
							ExecuteCommand(
								&commands.Transfer{
									TransactionID: "T002",
									FromAccountID: "A001",
									ToAccountID:   "A002",
									Amount:        500,
									ScheduledTime: time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC),
								},
							),
							ToRecordEvent(
								&events.TransferApproved{
									TransactionID: "T002",
									FromAccountID: "A001",
									ToAccountID:   "A002",
									Amount:        500,
								},
							),
						).
						// verify that funds are availalbe
						Expect(
							ExecuteCommand(
								&commands.Withdraw{
									TransactionID: "W001",
									AccountID:     "A002",
									Amount:        100,
									ScheduledTime: time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC),
								},
							),
							ToRecordEvent(
								&events.WithdrawalApproved{
									TransactionID: "W001",
									AccountID:     "A002",
									Amount:        100,
								},
							),
						)
				},
			)
		},
	)

	t.Run(
		"when the transfer exceeds the daily debit limit",
		func(t *testing.T) {
			t.Run(
				"it does not transfer any funds from the account",
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
							ExecuteCommand(
								&commands.OpenAccount{
									CustomerID:  "C002",
									AccountID:   "A002",
									AccountName: "Bob Jones",
								},
							),
							ExecuteCommand(
								&commands.Deposit{
									TransactionID: "D001",
									AccountID:     "A001",
									Amount:        expectedDailyDebitLimit + 10000,
								},
							),
						).
						Expect(
							ExecuteCommand(
								&commands.Transfer{
									TransactionID: "T001",
									FromAccountID: "A001",
									ToAccountID:   "A002",
									Amount:        expectedDailyDebitLimit + 1,
									ScheduledTime: time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC),
								},
							),
							ToRecordEvent(
								&events.TransferDeclined{
									TransactionID: "T001",
									FromAccountID: "A001",
									ToAccountID:   "A002",
									Amount:        expectedDailyDebitLimit + 1,
									Reason:        messages.DailyDebitLimitExceeded,
								},
							),
						).
						// verify that funds are not availalbe
						Expect(
							ExecuteCommand(
								&commands.Withdraw{
									TransactionID: "W001",
									AccountID:     "A002",
									Amount:        100,
									ScheduledTime: time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC),
								},
							),
							ToRecordEvent(
								&events.WithdrawalDeclined{
									TransactionID: "W001",
									AccountID:     "A002",
									Amount:        100,
									Reason:        messages.InsufficientFunds,
								},
							),
						)
				},
			)
		},
	)

	t.Run(
		"when the transfer is scheduled for a future date",
		func(t *testing.T) {
			t.Run(
				"it transfers the funds after the scheduled time",
				func(t *testing.T) {
					Begin(
						t,
						&example.App{},
						StartTimeAt(
							time.Date(2001, time.February, 3, 11, 22, 33, 0, time.UTC),
						),
					).
						Prepare(
							ExecuteCommand(
								&commands.OpenAccount{
									CustomerID:  "C001",
									AccountID:   "A001",
									AccountName: "Anna Smith",
								},
							),
							ExecuteCommand(
								&commands.OpenAccount{
									CustomerID:  "C002",
									AccountID:   "A002",
									AccountName: "Bob Jones",
								},
							),
							ExecuteCommand(
								&commands.Deposit{
									TransactionID: "D001",
									AccountID:     "A001",
									Amount:        500,
								},
							),
						).
						Expect(
							ExecuteCommand(
								&commands.Transfer{
									TransactionID: "T001",
									FromAccountID: "A001",
									ToAccountID:   "A002",
									Amount:        100,
									ScheduledTime: time.Date(2001, time.February, 4, 0, 0, 0, 0, time.UTC),
								},
							),
							NoneOf(
								ToRecordEventOfType(&events.TransferApproved{}),
							),
						).
						Expect(
							AdvanceTime(
								ToTime(time.Date(2001, time.February, 4, 0, 0, 0, 0, time.UTC)),
							),
							ToRecordEvent(
								&events.TransferApproved{
									TransactionID: "T001",
									FromAccountID: "A001",
									ToAccountID:   "A002",
									Amount:        100,
								},
							),
						).
						// verify that funds are availalbe
						Expect(
							ExecuteCommand(
								&commands.Withdraw{
									TransactionID: "W001",
									AccountID:     "A002",
									Amount:        100,
									ScheduledTime: time.Date(2001, time.February, 4, 0, 0, 0, 0, time.UTC),
								},
							),
							ToRecordEvent(
								&events.WithdrawalApproved{
									TransactionID: "W001",
									AccountID:     "A002",
									Amount:        100,
								},
							),
						)
				},
			)
		},
	)
}
