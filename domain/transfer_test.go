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

					whenAnnaHasSufficientFundsToTransferToBob := ExecuteCommand(commands.Transfer{
						TransactionID: "T001",
						FromAccountID: annaAccountID,
						ToAccountID:   bobAccountID,
						Amount:        amount,
						ScheduledTime: time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC),
					})

					theTransferIsApproved := ToRecordEvent(events.TransferApproved{
						TransactionID: "T001",
						FromAccountID: annaAccountID,
						ToAccountID:   bobAccountID,
						Amount:        amount,
					})

					whenBobWithdrawsTheTransferredFunds := ExecuteCommand(commands.Withdraw{
						TransactionID: "W001",
						AccountID:     bobAccountID,
						Amount:        amount,
						ScheduledTime: time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC),
					})

					theWithdrawalIsApproved := ToRecordEvent(events.WithdrawalApproved{
						TransactionID: "W001",
						AccountID:     bobAccountID,
						Amount:        amount,
					})

					Begin(t, &example.App{}).
						Prepare(
							openAccountForAnna,
							openAccountForBob,
							depositForAnna,
						).
						Expect(whenAnnaHasSufficientFundsToTransferToBob, theTransferIsApproved).
						// verify that funds are available
						Expect(whenBobWithdrawsTheTransferredFunds, theWithdrawalIsApproved)
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

					whenAnnaHasInsufficientFundsToTransferToBob := ExecuteCommand(commands.Transfer{
						TransactionID: "T001",
						FromAccountID: annaAccountID,
						ToAccountID:   bobAccountID,
						Amount:        amount,
						ScheduledTime: time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC),
					})

					theTransferIsDeclined := ToRecordEvent(events.TransferDeclined{
						TransactionID: "T001",
						FromAccountID: annaAccountID,
						ToAccountID:   bobAccountID,
						Amount:        amount,
						Reason:        messages.InsufficientFunds,
					})

					whenBobWithdrawsTheTransferredFunds := ExecuteCommand(commands.Withdraw{
						TransactionID: "W001",
						AccountID:     bobAccountID,
						Amount:        amount,
						ScheduledTime: time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC),
					})

					theWithdrawalIsDeclined := ToRecordEvent(events.WithdrawalDeclined{
						TransactionID: "W001",
						AccountID:     bobAccountID,
						Amount:        amount,
						Reason:        messages.InsufficientFunds,
					})

					Begin(t, &example.App{}).
						Prepare(
							openAccountForAnna,
							openAccountForBob,
							depositForAnna,
						).
						Expect(whenAnnaHasInsufficientFundsToTransferToBob, theTransferIsDeclined).
						// verify that funds are not available
						Expect(whenBobWithdrawsTheTransferredFunds, theWithdrawalIsDeclined)
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

					whenAnnaDoesNotExceedDailyDebitLimitToTransferToBob := ExecuteCommand(commands.Transfer{
						TransactionID: "T001",
						FromAccountID: annaAccountID,
						ToAccountID:   bobAccountID,
						Amount:        amount,
						ScheduledTime: time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC),
					})

					theTransferIsApproved := ToRecordEvent(events.TransferApproved{
						TransactionID: "T001",
						FromAccountID: annaAccountID,
						ToAccountID:   bobAccountID,
						Amount:        amount,
					})

					whenBobWithdrawsTheTransferredFunds := ExecuteCommand(commands.Withdraw{
						TransactionID: "W001",
						AccountID:     bobAccountID,
						Amount:        amount,
						ScheduledTime: time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC),
					})

					theWithdrawalIsApproved := ToRecordEvent(events.WithdrawalApproved{
						TransactionID: "W001",
						AccountID:     bobAccountID,
						Amount:        amount,
					})

					Begin(t, &example.App{}).
						Prepare(
							openAccountForAnna,
							openAccountForBob,
							depositAboveDailyDebitLimitForAnna,
						).
						Expect(whenAnnaDoesNotExceedDailyDebitLimitToTransferToBob, theTransferIsApproved).
						// verify that funds are available
						Expect(whenBobWithdrawsTheTransferredFunds, theWithdrawalIsApproved)
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

					whenAnnaDoesExceedDailyDebitLimitToTransferToBob := ExecuteCommand(commands.Transfer{
						TransactionID: "T001",
						FromAccountID: annaAccountID,
						ToAccountID:   bobAccountID,
						Amount:        amount,
						ScheduledTime: time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC),
					})

					theTransferIsDeclined := ToRecordEvent(events.TransferDeclined{
						TransactionID: "T001",
						FromAccountID: annaAccountID,
						ToAccountID:   bobAccountID,
						Amount:        amount,
						Reason:        messages.DailyDebitLimitExceeded,
					})

					whenBobWithdrawsTheTransferredFunds := ExecuteCommand(commands.Withdraw{
						TransactionID: "W001",
						AccountID:     bobAccountID,
						Amount:        amount,
						ScheduledTime: time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC),
					})

					theWithdrawalIsDeclined := ToRecordEvent(events.WithdrawalDeclined{
						TransactionID: "W001",
						AccountID:     bobAccountID,
						Amount:        amount,
						Reason:        messages.InsufficientFunds,
					})

					Begin(t, &example.App{}).
						Prepare(
							openAccountForAnna,
							openAccountForBob,
							depositAboveDailyDebitLimitForAnna,
						).
						Expect(whenAnnaDoesExceedDailyDebitLimitToTransferToBob, theTransferIsDeclined).
						// verify that funds are not available
						Expect(whenBobWithdrawsTheTransferredFunds, theWithdrawalIsDeclined)
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

					whenAnnaHasScheduledAFutureTransferToBob := ExecuteCommand(commands.Transfer{
						TransactionID: "T001",
						FromAccountID: annaAccountID,
						ToAccountID:   bobAccountID,
						Amount:        amount,
						ScheduledTime: time.Date(2001, time.February, 4, 0, 0, 0, 0, time.UTC),
					})

					theTransferIsNotYetApproved := NoneOf(
						ToRecordEventOfType(events.TransferApproved{}),
					)

					whenTimePassesTheScheduledTransferTime := AdvanceTime(
						ToTime(time.Date(2001, time.February, 4, 0, 0, 0, 0, time.UTC)),
					)

					theTransferIsApproved := ToRecordEvent(events.TransferApproved{
						TransactionID: "T001",
						FromAccountID: annaAccountID,
						ToAccountID:   bobAccountID,
						Amount:        amount,
					})

					whenBobWithdrawsTheTransferredFunds := ExecuteCommand(commands.Withdraw{
						TransactionID: "W001",
						AccountID:     bobAccountID,
						Amount:        amount,
						ScheduledTime: time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC),
					})

					theWithdrawalIsApproved := ToRecordEvent(events.WithdrawalApproved{
						TransactionID: "W001",
						AccountID:     bobAccountID,
						Amount:        amount,
					})

					Begin(
						t,
						&example.App{},
						StartTimeAt(
							time.Date(2001, time.February, 3, 11, 22, 33, 0, time.UTC),
						),
					).
						Prepare(
							openAccountForAnna,
							openAccountForBob,
							depositForAnna,
						).
						Expect(whenAnnaHasScheduledAFutureTransferToBob, theTransferIsNotYetApproved).
						Expect(whenTimePassesTheScheduledTransferTime, theTransferIsApproved).
						// verify that funds are available
						Expect(whenBobWithdrawsTheTransferredFunds, theWithdrawalIsApproved)
				},
			)
		},
	)
}
