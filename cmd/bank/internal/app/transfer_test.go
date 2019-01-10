package app_test

// import (
// 	. "github.com/dogmatiq/examples/cmd/bank/internal/app"
// 	"github.com/dogmatiq/examples/cmd/bank/internal/messages"
// 	"github.com/dogmatiq/examples/dogmatest"
// 	. "github.com/onsi/ginkgo"
// )

// var _ = Describe("Transfer", func() {
// 	engine := dogmatest.New(App)

// 	BeforeEach(func() {
// 		engine.Reset(
// 			messages.OpenAccount{
// 				AccountID: "A001",
// 				Name:      "Anna",
// 			},
// 			messages.OpenAccount{
// 				AccountID: "A002",
// 				Name:      "Bob",
// 			},
// 			messages.Deposit{
// 				TransactionID: "D001",
// 				AccountID:     "A001",
// 				Amount:        1000,
// 			},
// 		)
// 	})

// 	When("the creditor has sufficient funds", func() {
// 		It("moves funds from one account to another", func() {
// 			engine.
// 				ExecuteCommand(
// 					messages.Transfer{
// 						TransactionID: "T001",
// 						FromAccountID: "A001",
// 						ToAccountID:   "A002",
// 						Amount:        500,
// 					},
// 				).
// 				ExpectEvents(
// 					messages.AccountDebitedForTransfer{
// 						TransactionID: "T001",
// 						AccountID:     "A001",
// 						Amount:        500,
// 					},
// 					messages.AccountCreditedForTransfer{
// 						TransactionID: "T001",
// 						AccountID:     "A002",
// 						Amount:        500,
// 					},
// 				)
// 		})
// 	})

// 	When("the creditor has insufficient funds", func() {
// 		It("declines the transaction", func() {
// 			engine.
// 				ExecuteCommand(
// 					messages.Transfer{
// 						TransactionID: "T001",
// 						FromAccountID: "A001",
// 						ToAccountID:   "A002",
// 						Amount:        2000,
// 					},
// 				).
// 				ExpectEvents(
// 					messages.TransferDeclined{
// 						TransactionID: "T001",
// 						AccountID:     "A001",
// 						Amount:        2000,
// 					},
// 				)
// 		})
// 	})
// })
