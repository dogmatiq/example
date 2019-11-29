package testrunner

import (
	"database/sql"

	"github.com/dogmatiq/example"
	"github.com/dogmatiq/testkit"
)

// New returns a test runner for the example application.
func New(db *sql.DB) *testkit.Runner {
	app, err := example.NewApp(db)
	if err != nil {
		panic(err)
	}

	return testkit.New(app)
}
