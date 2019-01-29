package transaction_test

import (
	"testing"

	. "github.com/dogmatiq/dogmatest"
	"github.com/dogmatiq/example/internal/testrunner"
	"github.com/dogmatiq/example/messages"
)

func TestTransferProcess_SufficientFunds(t *testing.T) {
	testrunner.Runner.
		Begin(t).
		Setup(
			messages.OpenAccount{
				AccountID: "A001",
				Name:      "Anna",
			},
			messages.OpenAccount{
				AccountID: "A002",
				Name:      "Bob",
			},
			messages.Deposit{
				TransactionID: "D001",
				AccountID:     "A001",
				Amount:        500,
			},
		).
		ExecuteCommand(
			messages.Transfer{
				TransactionID: "T001",
				FromAccountID: "A001",
				ToAccountID:   "A002",
				Amount:        100,
			},
			AllOf(
				EventRecorded(
					messages.AccountDebitedForTransfer{
						TransactionID: "T001",
						AccountID:     "A001",
						Amount:        100,
					},
				),
				EventRecorded(
					messages.AccountCreditedForTransfer{
						TransactionID: "T001",
						AccountID:     "A002",
						Amount:        100,
					},
				),
			),
		)
}

func TestTransferProcess_InsufficientFunds(t *testing.T) {
	testrunner.Runner.
		Begin(t).
		Setup(
			messages.OpenAccount{
				AccountID: "A001",
				Name:      "Anna",
			},
			messages.OpenAccount{
				AccountID: "A002",
				Name:      "Bob",
			},
			messages.Deposit{
				TransactionID: "D001",
				AccountID:     "A001",
				Amount:        500,
			},
		).
		ExecuteCommand(
			messages.Transfer{
				TransactionID: "T001",
				FromAccountID: "A001",
				ToAccountID:   "A002",
				Amount:        1000,
			},
			AllOf(
				EventRecorded(
					messages.TransferDeclined{
						TransactionID: "T001",
						AccountID:     "A001",
						Amount:        1000,
					},
				),
				NoneOf(
					EventTypeRecorded(messages.AccountDebitedForTransfer{}),
					EventTypeRecorded(messages.AccountCreditedForTransfer{}),
				),
			),
		)
}
