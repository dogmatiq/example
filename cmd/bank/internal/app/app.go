package app

import (
	"io"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/examples/cmd/bank/internal/app/projections"
)

// App is the example bank app.
type App struct {
	accountProjection projections.AccountProjectionHandler
}

// GenerateAccountCSV generates CSV of accounts and their balances, sorted by
// the current balance in descending order.
func (a *App) GenerateAccountCSV(w io.Writer) error {
	return a.accountProjection.GenerateCSV(w)
}

// Dogma returns the Dogma app definition for this app.
func (a *App) Dogma() dogma.App {
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
		Projections: []dogma.ProjectionMessageHandler{
			&a.accountProjection,
		},
	}
}
