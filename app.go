package example

import (
	"io"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/domain"
	"github.com/dogmatiq/example/projections"
)

// App is an implementation of dogma.Application for the bank example.
type App struct {
	accountAggregate         domain.AccountHandler
	customerAggregate        domain.CustomerHandler
	dailyDebitLimitAggregate domain.DailyDebitLimitHandler
	transactionAggregate     domain.TransactionHandler

	depositProcess                   domain.DepositProcessHandler
	openAccountForNewCustomerProcess domain.OpenAccountForNewCustomerProcessHandler
	transferProcess                  domain.TransferProcessHandler
	withdrawalProcess                domain.WithdrawalProcessHandler

	accountProjection projections.AccountProjectionHandler
}

// Configure configures the Dogma engine for this application.
func (a *App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("bank", "22028264-0bca-43e1-8d9d-cd094efb10b7")

	c.RegisterAggregate(a.accountAggregate)
	c.RegisterAggregate(a.customerAggregate)
	c.RegisterAggregate(a.dailyDebitLimitAggregate)
	c.RegisterAggregate(a.transactionAggregate)

	c.RegisterProcess(a.depositProcess)
	c.RegisterProcess(a.openAccountForNewCustomerProcess)
	c.RegisterProcess(a.transferProcess)
	c.RegisterProcess(a.withdrawalProcess)

	c.RegisterProjection(&a.accountProjection)
}

// GenerateAccountCSV generates CSV of accounts and their balances, sorted by
// the current balance in descending order.
func (a *App) GenerateAccountCSV(w io.Writer) error {
	return a.accountProjection.GenerateCSV(w)
}
