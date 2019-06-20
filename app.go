package example

import (
	"io"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/aggregate/account"
	"github.com/dogmatiq/example/aggregate/customer"
	"github.com/dogmatiq/example/aggregate/transaction"
	"github.com/dogmatiq/example/process/newcustomer/openaccount"
	"github.com/dogmatiq/example/process/transaction/deposit"
	"github.com/dogmatiq/example/process/transaction/transfer"
	"github.com/dogmatiq/example/process/transaction/withdrawal"
	"github.com/dogmatiq/example/projection"
)

// App is an implementation of dogma.Application for the bank example.
type App struct {
	accountAggregate     account.AggregateHandler
	customerAggregate    customer.AggregateHandler
	transactionAggregate transaction.AggregateHandler

	depositProcess     deposit.ProcessHandler
	openAccountProcess openaccount.ProcessHandler
	transferProcess    transfer.ProcessHandler
	withdrawalProcess  withdrawal.ProcessHandler

	accountProjection projection.AccountProjectionHandler
}

// Configure configures the Dogma engine for this application.
func (a *App) Configure(c dogma.ApplicationConfigurer) {
	c.Name("bank")

	c.RegisterAggregate(a.accountAggregate)
	c.RegisterAggregate(a.customerAggregate)
	c.RegisterAggregate(a.transactionAggregate)

	c.RegisterProcess(a.depositProcess)
	c.RegisterProcess(a.openAccountProcess)
	c.RegisterProcess(a.transferProcess)
	c.RegisterProcess(a.withdrawalProcess)

	c.RegisterProjection(&a.accountProjection)
}

// GenerateAccountCSV generates CSV of accounts and their balances, sorted by
// the current balance in descending order.
func (a *App) GenerateAccountCSV(w io.Writer) error {
	return a.accountProjection.GenerateCSV(w)
}
