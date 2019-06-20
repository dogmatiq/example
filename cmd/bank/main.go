package main

import (
	"context"
	"os"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/message/command"
	"github.com/dogmatiq/testkit/engine"
)

func main() {
	app := &example.App{}

	en, err := engine.New(&example.App{})
	if err != nil {
		panic(err)
	}

	messages := []dogma.Message{
		command.OpenAccountForNewCustomer{
			CustomerID:   "cust1",
			CustomerName: "Anna Smith",
			AccountID:    "acct1",
			AccountName:  "Anna Smith",
		},
		command.OpenAccountForNewCustomer{
			CustomerID:   "cust2",
			CustomerName: "Bob Jones",
			AccountID:    "acct2",
			AccountName:  "Bob Jones",
		},
		command.Deposit{
			TransactionID: "txn1",
			AccountID:     "acct1",
			Amount:        10000,
		},
		command.Transfer{
			TransactionID: "txn2",
			FromAccountID: "acct1",
			ToAccountID:   "acct2",
			Amount:        2500,
		},
	}

	for _, m := range messages {
		err := en.Dispatch(
			context.Background(),
			m,
			// engine.WithObserver(
			// 	fact.ObserverFunc(func(f fact.Fact) {
			// 		dapper.Print(f)
			// 		fmt.Print("\n\n")
			// 	}),
			// ),
			engine.EnableProjections(true),
		)
		if err != nil {
			panic(err)
		}
	}

	if err := app.GenerateAccountCSV(os.Stdout); err != nil {
		panic(err)
	}
}
