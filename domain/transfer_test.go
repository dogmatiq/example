package domain_test

import (
	"testing"

	"github.com/dogmatiq/example/internal/testrunner"
	"github.com/dogmatiq/example/messages"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
	. "github.com/dogmatiq/testkit/assert"
)

func Test_Transfer(t *testing.T) {
	t.Run(
		"when transfer with sufficient funds",
		func(t *testing.T) {
			t.Run(
				"it transfers the funds from one account to another",
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
								Amount:        500,
							},
						).
						ExecuteCommand(
							commands.Transfer{
								TransactionID: "T001",
								FromAccountID: "A001",
								ToAccountID:   "A002",
								Amount:        100,
							},
							EventRecorded(
								events.TransferApproved{
									TransactionID: "T001",
									FromAccountID: "A001",
									ToAccountID:   "A002",
									Amount:        100,
								},
							),
						)
				},
			)
		},
	)

	t.Run(
		"when transfer with insufficient funds",
		func(t *testing.T) {
			t.Run(
				"it does not transfer any funds from one account to another",
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
								Amount:        500,
							},
						).
						ExecuteCommand(
							commands.Transfer{
								TransactionID: "T001",
								FromAccountID: "A001",
								ToAccountID:   "A002",
								Amount:        1000,
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
						)
				},
			)
		},
	)

	t.Run(
		"when within daily debit limit",
		func(t *testing.T) {
			t.Run(
				"it transfers the funds from one account to another",
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
								ScheduledDate: businessDateToday,
							},
							EventRecorded(
								events.TransferApproved{
									TransactionID: "T002",
									FromAccountID: "A001",
									ToAccountID:   "A002",
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
								ScheduledDate: businessDateToday,
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
						)
				},
			)
		},
	)
}
