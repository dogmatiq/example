package transaction_test

import (
	"testing"
	"time"

	"github.com/dogmatiq/example/internal/testrunner"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
	. "github.com/dogmatiq/testkit/assert"
)

func TestTransferProcess_SufficientFunds(t *testing.T) {
	timestamp := time.Now()

	testrunner.Runner.
		Begin(t).
		Prepare(
			commands.OpenAccount{
				AccountID: "A001",
				Name:      "Anna",
			},
			commands.OpenAccount{
				AccountID: "A002",
				Name:      "Bob",
			},
			commands.Deposit{
				TransactionID: "D001",
				AccountID:     "A001",
				Amount:        500,
			},
		).
		ExecuteCommand(
			commands.Transfer{
				TransactionID:        "T002",
				FromAccountID:        "A001",
				ToAccountID:          "A002",
				Amount:               100,
				TransactionTimestamp: timestamp,
			},
			AllOf(
				EventRecorded(
					events.AccountDebitedForTransfer{
						TransactionID:        "T002",
						AccountID:            "A001",
						Amount:               100,
						TransactionTimestamp: timestamp,
					},
				),
				EventRecorded(
					events.AccountCreditedForTransfer{
						TransactionID: "T002",
						AccountID:     "A002",
						Amount:        100,
					},
				),
			),
		)
}

func TestTransferProcess_InsufficientFunds(t *testing.T) {
	timestamp := time.Now()

	testrunner.Runner.
		Begin(t).
		Prepare(
			commands.OpenAccount{
				AccountID: "A001",
				Name:      "Anna",
			},
			commands.OpenAccount{
				AccountID: "A002",
				Name:      "Bob",
			},
			commands.Deposit{
				TransactionID: "D001",
				AccountID:     "A001",
				Amount:        500,
			},
		).
		ExecuteCommand(
			commands.Transfer{
				TransactionID:        "T003",
				FromAccountID:        "A001",
				ToAccountID:          "A002",
				Amount:               1000,
				TransactionTimestamp: timestamp,
			},
			AllOf(
				EventRecorded(
					events.TransferDeclinedDueToInsufficientFunds{
						TransactionID: "T003",
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
}
