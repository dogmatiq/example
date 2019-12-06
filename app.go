package example

import (
	"database/sql"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/example/domain"
	"github.com/dogmatiq/example/projections"
	pksql "github.com/dogmatiq/projectionkit/sql"
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

	CustomerProjection dogma.ProjectionMessageHandler
	AccountProjection  dogma.ProjectionMessageHandler
}

// NewApp returns the example application.
//
// If db is nil, it omits projection message handlers from the configuration.
func NewApp(db *sql.DB) (*App, error) {
	var (
		app = &App{}
		err error
	)

	if db != nil {
		app.CustomerProjection, err = pksql.New(
			db,
			&projections.CustomerProjectionHandler{},
			nil,
		)
		if err != nil {
			return nil, err
		}

		app.AccountProjection, err = pksql.New(
			db,
			&projections.AccountProjectionHandler{},
			nil,
		)
		if err != nil {
			return nil, err
		}
	}

	return app, nil
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

	if a.CustomerProjection != nil {
		c.RegisterProjection(a.CustomerProjection)
		c.RegisterProjection(a.AccountProjection)
	}
}
