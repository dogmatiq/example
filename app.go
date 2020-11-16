package example

import (
	"database/sql"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/domain"
	"github.com/dogmatiq/example/projections"
	pksql "github.com/dogmatiq/projectionkit/sql"
	pkmysql "github.com/dogmatiq/projectionkit/sql/mysql"
)

// AppKey is the example application's identity key.
const AppKey = "22028264-0bca-43e1-8d9d-cd094efb10b7"

// App is a dogma.Application implementation for the example "bank" domain.
type App struct {
	accountAggregate         domain.AccountHandler
	customerAggregate        domain.CustomerHandler
	dailyDebitLimitAggregate domain.DailyDebitLimitHandler
	transactionAggregate     domain.TransactionHandler

	depositProcess                   domain.DepositProcessHandler
	openAccountForNewCustomerProcess domain.OpenAccountForNewCustomerProcessHandler
	transferProcess                  domain.TransferProcessHandler
	withdrawalProcess                domain.WithdrawalProcessHandler

	accountProjection  projections.AccountProjectionHandler
	customerProjection projections.CustomerProjectionHandler

	// ReadDB is the database to use for read-models. If it is nil the
	// projection message handlers are omitted from the application
	// configuration.
	ReadDB *sql.DB
}

// Configure configures the Dogma engine for this application.
func (a *App) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("bank", AppKey)

	c.RegisterAggregate(a.accountAggregate)
	c.RegisterAggregate(a.customerAggregate)
	c.RegisterAggregate(a.dailyDebitLimitAggregate)
	c.RegisterAggregate(a.transactionAggregate)

	c.RegisterProcess(a.depositProcess)
	c.RegisterProcess(a.openAccountForNewCustomerProcess)
	c.RegisterProcess(a.transferProcess)
	c.RegisterProcess(a.withdrawalProcess)

	if a.ReadDB != nil {
		c.RegisterProjection(
			pksql.MustNew(
				a.ReadDB,
				&a.accountProjection,
				&pkmysql.Driver{},
			),
		)
		c.RegisterProjection(
			pksql.MustNew(
				a.ReadDB,
				&a.customerProjection,
				&pkmysql.Driver{},
			),
		)
	}
}
