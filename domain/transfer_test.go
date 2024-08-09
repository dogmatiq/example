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

func Test_Transfer_Refactor(t *testing.T) {
	annaCustomerID := "C001"
	bobCustomerID := "C002"

	annaAccountID := "A001"
	bobAccountID := "A002"

	transactionID := "T001"
	scheduledTime := time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC)

	openAccountForAnna := ExecuteCommand(commands.OpenAccount{
		CustomerID:  annaCustomerID,
		AccountID:   annaAccountID,
		AccountName: "Anna Smith",
	})

	openAccountForBob := ExecuteCommand(commands.OpenAccount{
		CustomerID:  bobCustomerID,
		AccountID:   bobAccountID,
		AccountName: "Bob Jones",
	})

	depositForAnna := ExecuteCommand(commands.Deposit{
		TransactionID: "D001",
		AccountID:     annaAccountID,
		Amount:        500,
	})

	depositAboveDailyDebitLimitForAnna := ExecuteCommand(commands.Deposit{
		TransactionID: "D001",
		AccountID:     annaAccountID,
		Amount:        expectedDailyDebitLimit + 1,
	})

	t.Run(
		"when there are sufficient funds",
		func(t *testing.T) {
			t.Run(
				"it transfers the funds from one account to another",
				func(t *testing.T) {
					var amount int64 = 100

					whenAnnaSendsTransferToBobWithSufficientFunds := ExecuteCommand(commands.Transfer{
						TransactionID: transactionID,
						FromAccountID: annaAccountID,
						ToAccountID:   bobAccountID,
						Amount:        amount,
						ScheduledTime: scheduledTime,
					})

					theTransferIsApproved := ToRecordEvent(events.TransferApproved{
						TransactionID: transactionID,
						FromAccountID: annaAccountID,
						ToAccountID:   bobAccountID,
						Amount:        amount,
					})

					annaIsDebited := ToRecordEvent(events.AccountDebited{
						TransactionID:   transactionID,
						AccountID:       annaAccountID,
						TransactionType: messages.Transfer,
						Amount:          amount,
						ScheduledTime:   scheduledTime,
					})

					bobIsCredited := ToRecordEvent(events.AccountCredited{
						TransactionID:   transactionID,
						AccountID:       bobAccountID,
						TransactionType: messages.Transfer,
						Amount:          amount,
					})

					Begin(t, &example.App{}).
						Prepare(
							openAccountForAnna,
							openAccountForBob,
							depositForAnna,
						).
						Expect(
							whenAnnaSendsTransferToBobWithSufficientFunds,
							AllOf(
								theTransferIsApproved,
								annaIsDebited,
								bobIsCredited,
							),
						)
				},
			)
		},
	)

	t.Run(
		"when there are insufficient funds",
		func(t *testing.T) {
			t.Run(
				"it does not transfer any funds from the account",
				func(t *testing.T) {
					var amount int64 = 1000

					whenAnnaSendsTransferToBobWithInsufficientFunds := ExecuteCommand(commands.Transfer{
						TransactionID: transactionID,
						FromAccountID: annaAccountID,
						ToAccountID:   bobAccountID,
						Amount:        amount,
						ScheduledTime: scheduledTime,
					})

					theTransferIsDeclined := ToRecordEvent(events.TransferDeclined{
						TransactionID: transactionID,
						FromAccountID: annaAccountID,
						ToAccountID:   bobAccountID,
						Amount:        amount,
						Reason:        messages.InsufficientFunds,
					})

					annaIsDebited := ToRecordEventOfType(events.AccountDebited{})
					bobIsCredited := ToRecordEventOfType(events.AccountCredited{})

					Begin(t, &example.App{}).
						Prepare(
							openAccountForAnna,
							openAccountForBob,
							depositForAnna,
						).
						Expect(
							whenAnnaSendsTransferToBobWithInsufficientFunds,
							AllOf(
								theTransferIsDeclined,
								NoneOf(
									annaIsDebited,
									bobIsCredited,
								),
							),
						)
				},
			)
		},
	)

	t.Run(
		"when the transfer does not exceed the daily debit limit",
		func(t *testing.T) {
			t.Run(
				"it transfers the funds from one account to another",
				func(t *testing.T) {
					var amount int64 = expectedDailyDebitLimit

					whenAnnaSendsTransferToBobWithoutExceedingDailyDebitLimit := ExecuteCommand(commands.Transfer{
						TransactionID: transactionID,
						FromAccountID: annaAccountID,
						ToAccountID:   bobAccountID,
						Amount:        amount,
						ScheduledTime: scheduledTime,
					})

					theTransferIsApproved := ToRecordEvent(events.TransferApproved{
						TransactionID: transactionID,
						FromAccountID: annaAccountID,
						ToAccountID:   bobAccountID,
						Amount:        amount,
					})

					annaIsDebited := ToRecordEvent(events.AccountDebited{
						TransactionID:   transactionID,
						AccountID:       annaAccountID,
						TransactionType: messages.Transfer,
						Amount:          amount,
						ScheduledTime:   scheduledTime,
					})

					bobIsCredited := ToRecordEvent(events.AccountCredited{
						TransactionID:   transactionID,
						AccountID:       bobAccountID,
						TransactionType: messages.Transfer,
						Amount:          amount,
					})

					Begin(t, &example.App{}).
						Prepare(
							openAccountForAnna,
							openAccountForBob,
							depositAboveDailyDebitLimitForAnna,
						).
						Expect(
							whenAnnaSendsTransferToBobWithoutExceedingDailyDebitLimit,
							AllOf(
								theTransferIsApproved,
								annaIsDebited,
								bobIsCredited,
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
					var amount int64 = expectedDailyDebitLimit + 1

					whenAnnaSendsTransferToBobAndExceedsDailyDebitLimit := ExecuteCommand(commands.Transfer{
						TransactionID: transactionID,
						FromAccountID: annaAccountID,
						ToAccountID:   bobAccountID,
						Amount:        amount,
						ScheduledTime: scheduledTime,
					})

					theTransferIsDeclined := ToRecordEvent(events.TransferDeclined{
						TransactionID: transactionID,
						FromAccountID: annaAccountID,
						ToAccountID:   bobAccountID,
						Amount:        amount,
						Reason:        messages.DailyDebitLimitExceeded,
					})

					bobIsCredited := ToRecordEvent(events.AccountCredited{
						TransactionID:   transactionID,
						AccountID:       bobAccountID,
						TransactionType: messages.Transfer,
						Amount:          amount,
					})

					Begin(t, &example.App{}).
						Prepare(
							openAccountForAnna,
							openAccountForBob,
							depositAboveDailyDebitLimitForAnna,
						).
						Expect(
							whenAnnaSendsTransferToBobAndExceedsDailyDebitLimit,
							AllOf(
								theTransferIsDeclined,
								NoneOf(
									bobIsCredited,
								),
							),
						)
				},
			)
		},
	)

	t.Run(
		"when the transfer is scheduled for a future date",
		func(t *testing.T) {
			t.Run(
				"it transfers the funds after the scheduled time",
				func(t *testing.T) {
					var amount int64 = 100
					startTime := time.Date(2001, time.February, 3, 11, 22, 33, 0, time.UTC)
					timeInFuture := time.Date(2001, time.February, 4, 0, 0, 0, 0, time.UTC)

					whenAnnaSendsTransferToBobScheduledInTheFuture := ExecuteCommand(commands.Transfer{
						TransactionID: transactionID,
						FromAccountID: annaAccountID,
						ToAccountID:   bobAccountID,
						Amount:        amount,
						ScheduledTime: timeInFuture,
					})

					theTransferIsNotYetApproved := NoneOf(
						ToRecordEventOfType(events.TransferApproved{}),
					)

					whenTimePassesTheScheduledTransferTime := AdvanceTime(
						ToTime(timeInFuture),
					)

					theTransferIsApproved := ToRecordEvent(events.TransferApproved{
						TransactionID: transactionID,
						FromAccountID: annaAccountID,
						ToAccountID:   bobAccountID,
						Amount:        amount,
					})

					annaIsDebited := ToRecordEvent(events.AccountDebited{
						TransactionID:   transactionID,
						AccountID:       annaAccountID,
						TransactionType: messages.Transfer,
						Amount:          amount,
						ScheduledTime:   timeInFuture,
					})

					bobIsCredited := ToRecordEvent(events.AccountCredited{
						TransactionID:   transactionID,
						AccountID:       bobAccountID,
						TransactionType: messages.Transfer,
						Amount:          amount,
					})

					Begin(
						t,
						&example.App{},
						StartTimeAt(startTime),
					).
						Prepare(
							openAccountForAnna,
							openAccountForBob,
							depositForAnna,
						).
						Expect(
							whenAnnaSendsTransferToBobScheduledInTheFuture,
							AllOf(
								theTransferIsNotYetApproved,
								NoneOf(
									annaIsDebited,
									bobIsCredited,
								),
							),
						).
						Expect(
							whenTimePassesTheScheduledTransferTime,
							AllOf(
								theTransferIsApproved,
								annaIsDebited,
								bobIsCredited,
							),
						)
				},
			)
		},
	)
}
