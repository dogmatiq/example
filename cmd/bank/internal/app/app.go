package app

import (
	"github.com/dogmatiq/dogma"
)

// App is the Dogma application for the bank example.
var App dogma.App

func init() {
	App = dogma.App{
		Name: "bank",
		Aggregates: []dogma.AggregateMessageHandler{
			AccountHandler,
			TransactionHandler,
		},
		Processes: []dogma.ProcessMessageHandler{
			DepositProcessHandler,
			WithdrawalProcessHandler,
			TransferProcessHandler,
		},
	}
}
