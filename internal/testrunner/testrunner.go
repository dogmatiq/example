package testrunner

import (
	"github.com/dogmatiq/example"
	"github.com/dogmatiq/testkit"
)

// Runner is a test runner for the example app.
var Runner = testkit.New(&example.App{})
