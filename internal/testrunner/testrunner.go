package testrunner

import (
	"github.com/dogmatiq/example"
	"github.com/dogmatiq/example/internal/database"
	"github.com/dogmatiq/testkit"
)

// Runner is a test runner for the example app.
var Runner *testkit.Runner

func init() {
	app, err := example.NewApp(database.New())
	if err != nil {
		panic(err)
	}

	Runner = testkit.New(app)
}
