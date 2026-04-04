package domain_test

import (
	"testing"
	"time"

	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/messages"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
	. "github.com/dogmatiq/testkit"
)

func Test_Transfer(t *testing.T) {
	cases := []struct {
		Name     string
		Transfer *commands.Transfer
	}{
		{"it transfers to an in-house account", &commands.Transfer{
			ToAccountID: "A002",
		}},
		{"it transfers to a third-party account", &commands.Transfer{
			ToAccountID:      "100001",
			ToThirdPartyBank: true,
		}},
	}

	for _, c := range cases {
		t.Run(
			c.Name,
			func(t *testing.T) {
				app := func(opts ...TestOption) *Test {
					a := Begin(t, &example.App{}, opts...)
					if c.Transfer.ToThirdPartyBank {
						a = a.EnableHandlers("third-party-bank")
					}
					return a
				}

				t.Run(
					"when there are sufficient funds",
					func(t *testing.T) {
						transfer := *c.Transfer
						transfer.TransactionID = "T001"
						transfer.FromAccountID = "A001"
						transfer.Amount = 100
						transfer.ScheduledTime = time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC)

						app().
							Prepare(
								ExecuteCommand(
									&commands.OpenAccount{
										CustomerID:  "C001",
										AccountID:   "A001",
										AccountName: "Anna Smith",
									},
								),
								ExecuteCommand(
									&commands.OpenAccount{
										CustomerID:  "C002",
										AccountID:   "A002",
										AccountName: "Bob Jones",
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
								ExecuteCommand(&transfer),
								ToRecordEvent(
									&events.TransferApproved{
										TransactionID: "T001",
										FromAccountID: "A001",
										ToAccountID:   c.Transfer.ToAccountID,
										Amount:        100,
									},
								),
							)
					},
				)

				t.Run(
					"when the transfer does not exceed the daily debit limit",
					func(t *testing.T) {
						transfer := *c.Transfer
						transfer.TransactionID = "T002"
						transfer.FromAccountID = "A001"
						transfer.Amount = 500
						transfer.ScheduledTime = time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC)

						app().
							Prepare(
								ExecuteCommand(
									&commands.OpenAccount{
										CustomerID:  "C001",
										AccountID:   "A001",
										AccountName: "Anna Smith",
									},
								),
								ExecuteCommand(
									&commands.OpenAccount{
										CustomerID:  "C002",
										AccountID:   "A002",
										AccountName: "Bob Jones",
									},
								),
								ExecuteCommand(
									&commands.Deposit{
										TransactionID: "D001",
										AccountID:     "A001",
										Amount:        expectedDailyDebitLimit + 10000,
									},
								),
							).
							Expect(
								ExecuteCommand(&transfer),
								ToRecordEvent(
									&events.TransferApproved{
										TransactionID: "T002",
										FromAccountID: "A001",
										ToAccountID:   c.Transfer.ToAccountID,
										Amount:        500,
									},
								),
							)
					},
				)

				t.Run(
					"when the transfer is scheduled for a future date",
					func(t *testing.T) {
						transfer := *c.Transfer
						transfer.TransactionID = "T001"
						transfer.FromAccountID = "A001"
						transfer.Amount = 100
						transfer.ScheduledTime = time.Date(2001, time.February, 4, 0, 0, 0, 0, time.UTC)

						app(
							StartTimeAt(
								time.Date(2001, time.February, 3, 11, 22, 33, 0, time.UTC),
							),
						).
							Prepare(
								ExecuteCommand(
									&commands.OpenAccount{
										CustomerID:  "C001",
										AccountID:   "A001",
										AccountName: "Anna Smith",
									},
								),
								ExecuteCommand(
									&commands.OpenAccount{
										CustomerID:  "C002",
										AccountID:   "A002",
										AccountName: "Bob Jones",
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
								ExecuteCommand(&transfer),
								NoneOf(
									ToRecordEventOfType(&events.TransferApproved{}),
								),
							).
							Expect(
								AdvanceTime(
									ToTime(time.Date(2001, time.February, 4, 0, 0, 0, 0, time.UTC)),
								),
								ToRecordEvent(
									&events.TransferApproved{
										TransactionID: "T001",
										FromAccountID: "A001",
										ToAccountID:   c.Transfer.ToAccountID,
										Amount:        100,
									},
								),
							)
					},
				)
			},
		)
	}

	t.Run(
		"when there are insufficient funds",
		func(t *testing.T) {
			t.Run(
				"it does not transfer any funds from the account",
				func(t *testing.T) {
					Begin(t, &example.App{}).
						Prepare(
							ExecuteCommand(
								&commands.OpenAccount{
									CustomerID:  "C001",
									AccountID:   "A001",
									AccountName: "Anna Smith",
								},
							),
							ExecuteCommand(
								&commands.OpenAccount{
									CustomerID:  "C002",
									AccountID:   "A002",
									AccountName: "Bob Jones",
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
									TransactionID: "T001",
									FromAccountID: "A001",
									ToAccountID:   "A002",
									Amount:        1000,
									ScheduledTime: time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC),
								},
							),
							ToRecordEvent(
								&events.TransferDeclined{
									TransactionID: "T001",
									FromAccountID: "A001",
									ToAccountID:   "A002",
									Amount:        1000,
									Reason:        messages.InsufficientFunds,
								},
							),
						).
						// verify that funds are not available
						Expect(
							ExecuteCommand(
								&commands.Withdraw{
									TransactionID: "W001",
									AccountID:     "A002",
									Amount:        100,
									ScheduledTime: time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC),
								},
							),
							ToRecordEvent(
								&events.WithdrawalDeclined{
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
		"when the transfer exceeds the daily debit limit",
		func(t *testing.T) {
			t.Run(
				"it does not transfer any funds from the account",
				func(t *testing.T) {
					Begin(t, &example.App{}).
						Prepare(
							ExecuteCommand(
								&commands.OpenAccount{
									CustomerID:  "C001",
									AccountID:   "A001",
									AccountName: "Anna Smith",
								},
							),
							ExecuteCommand(
								&commands.OpenAccount{
									CustomerID:  "C002",
									AccountID:   "A002",
									AccountName: "Bob Jones",
								},
							),
							ExecuteCommand(
								&commands.Deposit{
									TransactionID: "D001",
									AccountID:     "A001",
									Amount:        expectedDailyDebitLimit + 10000,
								},
							),
						).
						Expect(
							ExecuteCommand(
								&commands.Transfer{
									TransactionID: "T001",
									FromAccountID: "A001",
									ToAccountID:   "A002",
									Amount:        expectedDailyDebitLimit + 1,
									ScheduledTime: time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC),
								},
							),
							ToRecordEvent(
								&events.TransferDeclined{
									TransactionID: "T001",
									FromAccountID: "A001",
									ToAccountID:   "A002",
									Amount:        expectedDailyDebitLimit + 1,
									Reason:        messages.DailyDebitLimitExceeded,
								},
							),
						).
						// verify that funds are not available
						Expect(
							ExecuteCommand(
								&commands.Withdraw{
									TransactionID: "W001",
									AccountID:     "A002",
									Amount:        100,
									ScheduledTime: time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC),
								},
							),
							ToRecordEvent(
								&events.WithdrawalDeclined{
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
		"when the third-party credit fails",
		func(t *testing.T) {
			t.Run(
				"it refunds the source account",
				func(t *testing.T) {
					// Integration handler is intentionally not enabled so the
					// process stalls after issuing CreditThirdPartyAccount, allowing
					// us to inject ThirdPartyAccountCreditFailed directly.
					Begin(t, &example.App{}).
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
						).
						Expect(
							RecordEvent(
								&events.ThirdPartyAccountCreditFailed{
									TransactionID: "T001",
									AccountID:     "100001",
									Amount:        100,
								},
							),
							ToRecordEvent(
								&events.TransferFailed{
									TransactionID: "T001",
									FromAccountID: "A001",
									ToAccountID:   "100001",
									Amount:        100,
								},
							),
						).
						// verify that the funds were returned
						Expect(
							ExecuteCommand(
								&commands.Withdraw{
									TransactionID: "W001",
									AccountID:     "A001",
									Amount:        500,
									ScheduledTime: time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC),
								},
							),
							ToRecordEvent(
								&events.WithdrawalApproved{
									TransactionID: "W001",
									AccountID:     "A001",
									Amount:        500,
								},
							),
						)
				},
			)
		},
	)
}
