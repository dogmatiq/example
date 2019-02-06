package main

import (
	"context"
	"os"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/testkit/engine"
)

func main() {
	app := &example.App{}

	en, err := engine.New(&example.App{})
	if err != nil {
		panic(err)
	}

	messages := []dogma.Message{
		commands.OpenAccount{
			AccountID: "acct1",
			Name:      "Anna Smith",
		},
		commands.OpenAccount{
			AccountID: "acct2",
			Name:      "Bob Jones",
		},
		commands.Deposit{
			TransactionID: "txn1",
			AccountID:     "acct1",
			Amount:        10000,
		},
		commands.Transfer{
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
