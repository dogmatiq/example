package domain_test

import (
	"testing"
	"time"

	"github.com/dogmatiq/example/internal/testrunner"
	"github.com/dogmatiq/example/messages"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
	. "github.com/dogmatiq/testkit/assert"
)

func TestWithdrawalProcess_SufficientFunds(t *testing.T) {
	t.Run(
		"it withdraws some funds from an account",
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
						Amount:        5000,
					}).
				ExecuteCommand(
					commands.Withdraw{
						TransactionID:                 "T001",
						AccountID:                     "A001",
						Amount:                        500,
						RequestedTransactionTimestamp: time.Unix(12345, 0),
					},
					EventRecorded(
						events.AccountDebitedForWithdrawal{
							TransactionID: "T001",
							AccountID:     "A001",
							Amount:        500,
						},
					),
				)
		},
	)
}

func TestWithdrawalProcess_InsufficientFunds(t *testing.T) {
	t.Run(
		"it does not withdraw funds from an account if there is insufficient funds",
		func(t *testing.T) {
			testrunner.Runner.
				Begin(t).
				Prepare(
					commands.OpenAccount{
						CustomerID:  "C001",
						AccountID:   "A001",
						AccountName: "Anna Smith",
					}),
				ExecuteCommand(
					commands.Withdraw{
						TransactionID:                 "T001",
						AccountID:                     "A001",
						Amount:                        500,
						RequestedTransactionTimestamp: time.Unix(12345, 0),
					},
					EventRecorded(
						events.WithdrawalDeclined{
							TransactionID: "T001",
							AccountID:     "A001",
							Amount:        500,
							Reason: messages.ReasonInsufficientFunds,
						},
					),
				)
		},
	)
}
