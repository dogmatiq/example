package example

import (
	"io"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/account"
	"github.com/dogmatiq/example/projections"
	"github.com/dogmatiq/example/transaction"
)

// App is an implementation of dogma.Application for the bank example.
type App struct {
	accountAggregate         account.Aggregate
	transactionAggregate     transaction.Aggregate
	dailyDebitLimitAggregate transaction.DailyDebitLimitHandler
	depositProcess           transaction.DepositProcess
	withdrawalProcess        transaction.WithdrawalProcess
	transferProcess          transaction.TransferProcess
	accountProjection        projections.AccountProjectionHandler
}

// Configure configures the Dogma engine for this application.
func (a *App) Configure(c dogma.ApplicationConfigurer) {
	c.Name("bank")
	c.RegisterAggregate(a.accountAggregate)
	c.RegisterAggregate(a.transactionAggregate)
	c.RegisterAggregate(a.dailyDebitLimitAggregate)
	c.RegisterProcess(a.depositProcess)
	c.RegisterProcess(a.withdrawalProcess)
	c.RegisterProcess(a.transferProcess)
	c.RegisterProjection(&a.accountProjection)
}

// GenerateAccountCSV generates CSV of accounts and their balances, sorted by
// the current balance in descending order.
func (a *App) GenerateAccountCSV(w io.Writer) error {
	return a.accountProjection.GenerateCSV(w)
}
