package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/config/runtimeconfig"
	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/database"
	"github.com/dogmatiq/example/messages/commands"
	"github.com/dogmatiq/testkit/engine"
	"github.com/dogmatiq/testkit/fact"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db := database.MustNew()
	defer db.Close()

	app := &example.App{
		ReadDB: db,
	}

	en, err := engine.New(runtimeconfig.FromApplication(app))
	if err != nil {
		panic(err)
	}

	commands := []dogma.Command{
		&commands.OpenAccountForNewCustomer{
			CustomerID:   "cust1",
			CustomerName: "Anna Smith",
			AccountID:    "acct1",
			AccountName:  "Anna Smith",
		},
		&commands.OpenAccountForNewCustomer{
			CustomerID:   "cust2",
			CustomerName: "Bob Jones",
			AccountID:    "acct2",
			AccountName:  "Bob Jones",
		},
		&commands.Deposit{
			TransactionID: "txn1",
			AccountID:     "acct1",
			Amount:        10000,
		},
		&commands.Withdraw{
			TransactionID: "txn2",
			AccountID:     "acct1",
			Amount:        500,
			ScheduledTime: time.Now(),
		},
		&commands.Transfer{
			TransactionID: "txn3",
			FromAccountID: "acct1",
			ToAccountID:   "acct2",
			Amount:        2500,
			ScheduledTime: time.Now(),
		},
		&commands.Transfer{
			TransactionID: "txn4",
			FromAccountID: "acct1",
			ToAccountID:   "acct2",
			Amount:        500,
			ScheduledTime: time.Now().AddDate(0, 0, 1),
		},
	}

	for _, m := range commands {
		err := en.Dispatch(
			context.Background(),
			m,
			engine.WithObserver(
				fact.NewLogger(func(s string) {
					fmt.Println(s)
				}),
			),
			engine.EnableProjections(true),
		)
		if err != nil {
			panic(err)
		}
	}
}
