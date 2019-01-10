package app

import (
	"github.com/dogmatiq/dogma"
)

// New returns a new Dogma app for the banking application.
func New() dogma.App {
	return dogma.App{
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
