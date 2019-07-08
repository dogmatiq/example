package domain_test

import (
	"testing"
	"time"

	"github.com/dogmatiq/example/internal/testrunner"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
	. "github.com/dogmatiq/testkit/assert"
)

func TestDailyDebitLimit_ConsumeDailyDebitLimit(t *testing.T) {
	t.Run(
		"it consumes some of the daily debit limit",
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
					commands.ConsumeDailyDebitLimit{
						TransactionID:                 "T001",
						AccountID:                     "A001",
						Amount:                        500,
						RequestedTransactionTimestamp: time.Unix(12345, 0),
					},
					EventRecorded(
						events.DailyDebitLimitConsumed{
							TransactionID: "T001",
							AccountID:     "A001",
							Amount:        500,
							LimitUsed:     500,
							LimitMaximum:  900000,
						},
					),
				)
		},
	)

	t.Run(
		"it does not consume some of the daily debit limit if it will exceed the limit",
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
					commands.ConsumeDailyDebitLimit{
						TransactionID:                 "T001",
						AccountID:                     "A001",
						Amount:                        900100,
						RequestedTransactionTimestamp: time.Unix(12345, 0),
					},
					EventRecorded(
						events.DailyDebitLimitExceeded{
							TransactionID: "T001",
							AccountID:     "A001",
							Amount:        900100,
							LimitUsed:     0,
							LimitMaximum:  900000,
						},
					),
				)
		},
	)
}
