package domain_test

import (
	"testing"
	"time"

	"github.com/dogmatiq/example/internal/testrunner"
	"github.com/dogmatiq/example/messages"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
	. "github.com/dogmatiq/testkit"
	. "github.com/dogmatiq/testkit/assert"
)

func Test_Transfer(t *testing.T) {
	t.Run(
		"when there are sufficient funds",
		func(t *testing.T) {
			t.Run(
				"it transfers the funds from one account to another",
				func(t *testing.T) {
					testrunner.New(nil).
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
								Amount:        500,
							},
						).
						ExecuteCommand(
							commands.Transfer{
								TransactionID: "T001",
								FromAccountID: "A001",
								ToAccountID:   "A002",
								Amount:        100,
								ScheduledDate: "2001-02-03",
							},
							EventRecorded(
								events.TransferApproved{
									TransactionID: "T001",
									FromAccountID: "A001",
									ToAccountID:   "A002",
									Amount:        100,
								},
							),
						).
						// verify that funds are availalbe
						ExecuteCommand(
							commands.Withdraw{
								TransactionID: "W001",
								AccountID:     "A002",
								Amount:        100,
								ScheduledDate: "2001-02-03",
							},
							EventRecorded(
								events.WithdrawalApproved{
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
					testrunner.New(nil).
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
								Amount:        500,
							},
						).
						ExecuteCommand(
							commands.Transfer{
								TransactionID: "T001",
								FromAccountID: "A001",
								ToAccountID:   "A002",
								Amount:        1000,
								ScheduledDate: "2001-02-03",
							},
							EventRecorded(
								events.TransferDeclined{
									TransactionID: "T001",
									FromAccountID: "A001",
									ToAccountID:   "A002",
									Amount:        1000,
									Reason:        messages.InsufficientFunds,
								},
							),
						).
						// verify that funds are not availalbe
						ExecuteCommand(
							commands.Withdraw{
								TransactionID: "W001",
								AccountID:     "A002",
								Amount:        100,
							},
							EventRecorded(
								events.WithdrawalDeclined{
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
					testrunner.New(nil).
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
								TransactionID: "T001",
								AccountID:     "A001",
								Amount:        expectedDailyDebitLimit + 10000,
							},
						).
						ExecuteCommand(
							commands.Transfer{
								TransactionID: "T002",
								FromAccountID: "A001",
								ToAccountID:   "A002",
								Amount:        500,
								ScheduledDate: "2001-02-03",
							},
							EventRecorded(
								events.TransferApproved{
									TransactionID: "T002",
									FromAccountID: "A001",
									ToAccountID:   "A002",
									Amount:        500,
								},
							),
						).
						// verify that funds are availalbe
						ExecuteCommand(
							commands.Withdraw{
								TransactionID: "W001",
								AccountID:     "A002",
								Amount:        100,
								ScheduledDate: "2001-02-03",
							},
							EventRecorded(
								events.WithdrawalApproved{
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
					testrunner.New(nil).
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
						).
						ExecuteCommand(
							commands.Transfer{
								TransactionID: "T001",
								FromAccountID: "A001",
								ToAccountID:   "A002",
								Amount:        expectedDailyDebitLimit + 1,
								ScheduledDate: "2001-02-03",
							},
							EventRecorded(
								events.TransferDeclined{
									TransactionID: "T001",
									FromAccountID: "A001",
									ToAccountID:   "A002",
									Amount:        expectedDailyDebitLimit + 1,
									Reason:        messages.DailyDebitLimitExceeded,
								},
							),
						).
						// verify that funds are not availalbe
						ExecuteCommand(
							commands.Withdraw{
								TransactionID: "W001",
								AccountID:     "A002",
								Amount:        100,
								ScheduledDate: "2001-02-03",
							},
							EventRecorded(
								events.WithdrawalDeclined{
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
					testrunner.New(nil).
						Begin(
							t,
							WithStartTime(
								time.Date(2001, time.February, 3, 11, 22, 33, 0, time.UTC),
							),
						).
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
								Amount:        500,
							},
						).
						ExecuteCommand(
							commands.Transfer{
								TransactionID: "T001",
								FromAccountID: "A001",
								ToAccountID:   "A002",
								Amount:        100,
								ScheduledDate: "2001-02-04",
							},
							NoneOf(EventRecorded(&events.TransferApproved{})),
						).
						AdvanceTime(
							ToTime(time.Date(2001, time.February, 4, 0, 0, 0, 0, time.UTC)),
							EventRecorded(
								events.TransferApproved{
									TransactionID: "T001",
									FromAccountID: "A001",
									ToAccountID:   "A002",
									Amount:        100,
								},
							),
						).
						// verify that funds are availalbe
						ExecuteCommand(
							commands.Withdraw{
								TransactionID: "W001",
								AccountID:     "A002",
								Amount:        100,
								ScheduledDate: "2001-02-04",
							},
							EventRecorded(
								events.WithdrawalApproved{
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
		"when the transfer is to the same account",
		func(t *testing.T) {
			t.Run(
				"it does not start the transfer",
				func(t *testing.T) {
					cmd := commands.Transfer{
						TransactionID: "T001",
						FromAccountID: "A001",
						ToAccountID:   "A001",
						Amount:        100,
						ScheduledDate: "2001-02-04",
					}

					testrunner.New(nil).
						Begin(t).
						Prepare(cmd).
						ExecuteCommand(
							cmd,
							NoneOf(
								EventTypeRecorded(events.WithdrawalApproved{}),
							),
						)
				},
			)
		},
	)

	t.Run(
		"when the transfer has already started",
		func(t *testing.T) {
			t.Run(
				"it does not start the transfer again",
				func(t *testing.T) {
					cmd := commands.Transfer{
						TransactionID: "T001",
						FromAccountID: "A001",
						ToAccountID:   "A002",
						Amount:        100,
						ScheduledDate: "2001-02-04",
					}

					testrunner.New(nil).
						Begin(t).
						Prepare(cmd).
						ExecuteCommand(
							cmd,
							NoneOf(
								EventTypeRecorded(events.WithdrawalApproved{}),
							),
						)
				},
			)
		},
	)
}
