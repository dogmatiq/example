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

func TestDepositProcess_DepositFunds(t *testing.T) {
	t.Run(
		"it deposits some funds into an account",
		func(t *testing.T) {
			testrunner.Runner.
				Begin(t).
				Prepare(
					commands.OpenAccount{
						CustomerID:  "C001",
						AccountID:   "A001",
						AccountName: "Anna Smith",
					}).
				ExecuteCommand(
					commands.Deposit{
						TransactionID:                 "T001",
						AccountID:                     "A001",
						Amount:                        500,
					},
					EventRecorded(
						events.AccountCreditedForDeposit{
							TransactionID: "T001",
							AccountID:     "A001",
							Amount:        500,
						},
					),
				)
		},
	)
}
