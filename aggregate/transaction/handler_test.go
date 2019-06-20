package transaction_test

import (
	"testing"

	"github.com/dogmatiq/example/internal/testrunner"
	"github.com/dogmatiq/example/message/command"
	"github.com/dogmatiq/example/message/event"
	. "github.com/dogmatiq/testkit/assert"
)

func TestTransferProcess_SufficientFunds(t *testing.T) {
	testrunner.Runner.
		Begin(t).
		Prepare(
			command.OpenAccount{
				CustomerID:  "C001",
				AccountID:   "A001",
				AccountName: "Anna Smith",
			},
			command.OpenAccount{
				CustomerID:  "C002",
				AccountID:   "A002",
				AccountName: "Bob Jones",
			},
			command.Deposit{
				TransactionID: "D001",
				AccountID:     "A001",
				Amount:        500,
			},
		).
		ExecuteCommand(
			command.Transfer{
				TransactionID: "T001",
				FromAccountID: "A001",
				ToAccountID:   "A002",
				Amount:        100,
			},
			AllOf(
				EventRecorded(
					event.AccountDebitedForTransfer{
						TransactionID: "T001",
						AccountID:     "A001",
						Amount:        100,
					},
				),
				EventRecorded(
					event.AccountCreditedForTransfer{
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
		Prepare(
			command.OpenAccount{
				CustomerID:  "C001",
				AccountID:   "A001",
				AccountName: "Anna Smith",
			},
			command.OpenAccount{
				CustomerID:  "C002",
				AccountID:   "A002",
				AccountName: "Bob Jones",
			},
			command.Deposit{
				TransactionID: "D001",
				AccountID:     "A001",
				Amount:        500,
			},
		).
		ExecuteCommand(
			command.Transfer{
				TransactionID: "T001",
				FromAccountID: "A001",
				ToAccountID:   "A002",
				Amount:        1000,
			},
			AllOf(
				EventRecorded(
					event.TransferDeclinedDueToInsufficientFunds{
						TransactionID: "T001",
						AccountID:     "A001",
						Amount:        1000,
					},
				),
				NoneOf(
					EventTypeRecorded(event.AccountDebitedForTransfer{}),
					EventTypeRecorded(event.AccountCreditedForTransfer{}),
				),
			),
		)
}
