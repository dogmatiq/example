package domain_test

import (
	"testing"

	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/example/messages/events"
	"github.com/dogmatiq/testkit"
	. "github.com/dogmatiq/testkit"
)

func Test_OpeningAnAccount(t *testing.T) {
	annaHasAnAccount := testkit.Scenario(
		"anna has an open account",
		ExecuteCommand(commands.OpenAccountForNewCustomer{
			CustomerID:   "C010",
			CustomerName: "Anna Smith",
			AccountID:    "A010",
			AccountName:  "Anna Smith",
		}),
	)

	annaHasTwoAccounts := annaHasAnAccount.Scenario(
		"anna has two open accounts",
		ExecuteCommand(commands.OpenAccount{
			CustomerID:  "C010",
			AccountID:   "A011",
			AccountName: "Anna Smith",
		}),
	)
	t.Run("opening a customer's first account acquires the customer", func(t *testing.T) {
		Begin(t, &example.App{}).
			Expect(
				EnterScenario(annaHasAnAccount),
				ToRecordEvent(
					events.CustomerAcquired{
						CustomerID:   "C010",
						CustomerName: "Anna Smith",
						AccountID:    "A010",
						AccountName:  "Anna Smith",
					},
				),
			)
	})

	t.Run("opening a second account does not acquire a new customer", func(t *testing.T) {
		Begin(t, &example.App{}).
			Prepare(
				EnterScenario(annaHasAnAccount),
			).
			Expect(
				EnterScenario(annaHasTwoAccounts),
				NoneOf(ToRecordEventOfType(events.CustomerAcquired{})),
			)
	})
}

// t.Run("when the account holder is an existing customer", func(t *testing.T) {
// 	annaHasAnOpenAccount := testkit.
// 		Scenario("anna has an open account").
// 		ExecuteCommand(annaOpensHerFirstAccount)

// 	t.Run("opening a second account does not reacquire the customer", func(t *testing.T) {
// 		Begin(t, &example.App{}).
// 			Given(annaHasAnOpenAccount).
// 			Expect(
// 				ExecuteCommand(annaOpensASecondAccount),
// 				NoneOf(ToRecordEventOfType(events.CustomerAcquired{})),
// 			)
// 	})
// })

// 	t.Run(
// 		"when the customer does not exist",
// 		func(t *testing.T) {
// 			t.Run(
// 				"it acquires the customer",
// 				func(t *testing.T) {
// 					Begin(t, &example.App{}).
// 						Expect(
// 							ExecuteCommand(
// 								commands.OpenAccountForNewCustomer{
// 									CustomerID:   "C001",
// 									CustomerName: "Bob Jones",
// 									AccountID:    "A001",
// 									AccountName:  "Bob Jones",
// 								},
// 							),
// 							ToRecordEvent(
// 								events.CustomerAcquired{
// 									CustomerID:   "C001",
// 									CustomerName: "Bob Jones",
// 									AccountID:    "A001",
// 									AccountName:  "Bob Jones",
// 								},
// 							),
// 						)
// 				},
// 			)
// 		},
// 	)

// 	t.Run(
// 		"when the customer already exists",
// 		func(t *testing.T) {
// 			t.Run(
// 				"it does not reacquire the customer",
// 				func(t *testing.T) {
// 					cmd := commands.OpenAccountForNewCustomer{
// 						CustomerID:   "C001",
// 						CustomerName: "Bob Jones",
// 						AccountID:    "A001",
// 						AccountName:  "Bob Jones",
// 					}

// 					Begin(t, &example.App{}).
// 						Given(bobHasAnOpenAccount).
// 						Prepare(
// 							ExecuteCommand(cmd),
// 						).
// 						Expect(
// 							ExecuteCommand(cmd),
// 							NoneOf(
// 								ToRecordEventOfType(events.CustomerAcquired{}),
// 							),
// 						)
// 				},
// 			)
// 		},
// 	)

// 	t.Run(
// 		"when the account does not exist",
// 		func(t *testing.T) {
// 			t.Run(
// 				"the new account is opened",
// 				func(t *testing.T) {
// 					Begin(t, &example.App{}).
// 						Given(annaHasAnOpenAccount).
// 						Expect(
// 							ExecuteCommand(commands.OpenAccount{
// 								CustomerID:  "C001",
// 								AccountID:   "A002",
// 								AccountName: "Anna Smith's Second Account",
// 							}),
// 							ToRecordEvent(events.AccountOpened{
// 								CustomerID:  "C001",
// 								AccountID:   "A002",
// 								AccountName: "Anna Smith's Second Account",
// 							}),
// 						)
// 				},
// 			)
// 		},
// 	)

// 	t.Run(
// 		"when the account already exists",
// 		func(t *testing.T) {
// 			t.Run(
// 				"nothing happens to the existing account",
// 				func(t *testing.T) {
// 					Begin(t, &example.App{}).
// 						Given(annaHasAnOpenAccount).
// 						Expect(
// 							ExecuteCommand(commands.OpenAccount{
// 								CustomerID:  "C001",
// 								AccountID:   "A001",
// 								AccountName: "Anna Smith",
// 							}),
// 							NoneOf(ToRecordEventOfType(events.AccountOpened{})),
// 						)
// 				},
// 			)
// 		},
// 	)
// }
// }
