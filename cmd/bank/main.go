package main

import (
	"context"
	"os"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogmatest"
	"github.com/dogmatiq/examples/cmd/bank/internal/app"
	"github.com/dogmatiq/examples/cmd/bank/internal/messages"
)

func main() {
	a := app.App{}
	e := dogmatest.NewEngine(a.Dogma())

	messages := []dogma.Message{
		messages.OpenAccount{
			AccountID: "acct1",
			Name:      "Anna Smith",
		},
		messages.OpenAccount{
			AccountID: "acct2",
			Name:      "Bob Jones",
		},
		messages.Deposit{
			TransactionID: "txn1",
			AccountID:     "acct1",
			Amount:        10000,
		},
		messages.Transfer{
			TransactionID: "txn2",
			FromAccountID: "acct1",
			ToAccountID:   "acct2",
			Amount:        2500,
		},
	}

	for _, m := range messages {
		err := e.Dispatch(
			context.Background(),
			m,
			// engine.Observe(func(f fact.Fact) {
			// 	dapper.Print(f)
			// 	fmt.Print("\n\n")
			// }),
		)
		if err != nil {
			panic(err)
		}
	}

	if err := a.GenerateAccountCSV(os.Stdout); err != nil {
		panic(err)
	}
}
