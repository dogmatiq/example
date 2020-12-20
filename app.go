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

	AccountProjection  projections.AccountProjectionHandler
	CustomerProjection projections.CustomerProjectionHandler

	// ReadDB is the database to use for read-models. If it is nil the
	// projection message handlers are omitted from the application
	// configuration.
	ReadDB *sql.DB
}

// Configure configures the Dogma engine for this application.
func (a *App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("bank", AppKey)

	c.RegisterAggregate(a.AccountAggregate)
	c.RegisterAggregate(a.CustomerAggregate)
	c.RegisterAggregate(a.DailyDebitLimitAggregate)
	c.RegisterAggregate(a.TransactionAggregate)

	c.RegisterProcess(a.DepositProcess)
	c.RegisterProcess(a.OpenAccountForNewCustomerProcess)
	c.RegisterProcess(a.TransferProcess)
	c.RegisterProcess(a.WithdrawalProcess)

	if a.ReadDB != nil { // TODO: Remove this when testkit is updated
		c.RegisterProjection(
			sqlprojection.New(
				a.ReadDB,
				&a.AccountProjection,
			),
		)
		c.RegisterProjection(
			sqlprojection.New(
				a.ReadDB,
				&a.CustomerProjection,
			),
		)
	}
}
