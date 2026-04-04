package integrations_test

import (
	"testing"
	"time"

	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
	. "github.com/dogmatiq/testkit"
)

func Test_ThirdPartyBankIntegrationHandler(t *testing.T) {
	t.Run(
		"when a credit is requested",
		func(t *testing.T) {
			t.Run(
				"it credits the account if the account ID is numeric",
				func(t *testing.T) {
					Begin(t, &example.App{}).
						EnableHandlers("third-party-bank").
						Prepare(
							ExecuteCommand(
								&commands.OpenAccount{
									CustomerID:  "C001",
									AccountID:   "A001",
									AccountName: "Anna Smith",
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
									TransactionID:    "T001",
									FromAccountID:    "A001",
									ToAccountID:      "100001",
									Amount:           100,
									ScheduledTime:    time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC),
									ToThirdPartyBank: true,
								},
							),
							ToRecordEvent(
								&events.ThirdPartyAccountCredited{
									TransactionID: "T001",
									AccountID:     "100001",
									Amount:        100,
								},
							),
						)
				},
			)

			t.Run(
				"it fails the credit if the account ID is not numeric",
				func(t *testing.T) {
					Begin(t, &example.App{}).
						EnableHandlers("third-party-bank").
						Prepare(
							ExecuteCommand(
								&commands.OpenAccount{
									CustomerID:  "C001",
									AccountID:   "A001",
									AccountName: "Anna Smith",
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
									TransactionID:    "T001",
									FromAccountID:    "A001",
									ToAccountID:      "EXT001",
									Amount:           100,
									ScheduledTime:    time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC),
									ToThirdPartyBank: true,
								},
							),
							ToRecordEvent(
								&events.ThirdPartyAccountCreditFailed{
									TransactionID: "T001",
									AccountID:     "EXT001",
									Amount:        100,
								},
							),
						)
				},
			)
		},
	)
}
