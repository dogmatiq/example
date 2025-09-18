package example

import (
	"database/sql"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/domain"
	"github.com/dogmatiq/example/projections"
	"github.com/dogmatiq/projectionkit/sqlprojection"
)

// AppKey is the example application's identity key.
const AppKey = "22028264-0bca-43e1-8d9d-cd094efb10b7"

// App is a dogma.Application implementation for the example "bank" domain.
type App struct {
	AccountAggregate         domain.AccountHandler
	CustomerAggregate        domain.CustomerHandler
	DailyDebitLimitAggregate domain.DailyDebitLimitHandler
	TransactionAggregate     domain.TransactionHandler

	DepositProcess                   domain.DepositProcessHandler
	OpenAccountForNewCustomerProcess domain.OpenAccountForNewCustomerProcessHandler
	TransferProcess                  domain.TransferProcessHandler
	WithdrawalProcess                domain.WithdrawalProcessHandler

	ReadDB             *sql.DB
	AccountProjection  projections.AccountProjectionHandler
	CustomerProjection projections.CustomerProjectionHandler
}

// Configure configures the Dogma engine for this application.
func (a *App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("bank", AppKey)

	c.Routes(
		dogma.ViaAggregate(a.AccountAggregate),
		dogma.ViaAggregate(a.CustomerAggregate),
		dogma.ViaAggregate(a.DailyDebitLimitAggregate),
		dogma.ViaAggregate(a.TransactionAggregate),

		dogma.ViaProcess(a.DepositProcess),
		dogma.ViaProcess(a.OpenAccountForNewCustomerProcess),
		dogma.ViaProcess(a.TransferProcess),
		dogma.ViaProcess(a.WithdrawalProcess),

		dogma.ViaProjection(sqlprojection.New(a.ReadDB, sqlprojection.SQLiteDriver, &a.AccountProjection)),
		dogma.ViaProjection(sqlprojection.New(a.ReadDB, sqlprojection.SQLiteDriver, &a.CustomerProjection)),
	)
}
