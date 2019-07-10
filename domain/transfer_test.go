package domain_test

import (
	"testing"

	"github.com/dogmatiq/example/internal/testrunner"
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
							AllOf(
								EventRecorded(
									events.AccountDebitedForTransfer{
										TransactionID: "T001",
										AccountID:     "A001",
										Amount:        100,
									},
								),
								EventRecorded(
									events.AccountCreditedForTransfer{
										TransactionID: "T001",
										AccountID:     "A002",
										Amount:        100,
									},
								),
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
							AllOf(
								EventRecorded(
									events.TransferDeclinedDueToInsufficientFunds{
										TransactionID: "T001",
										AccountID:     "A001",
										Amount:        1000,
									},
								),
								NoneOf(
									EventTypeRecorded(events.AccountDebitedForTransfer{}),
									EventTypeRecorded(events.AccountCreditedForTransfer{}),
								),
							),
						)
				},
			)
		},
	)
}
