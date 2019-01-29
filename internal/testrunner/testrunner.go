package testrunner

import (
	"github.com/dogmatiq/dogmatest"
	"github.com/dogmatiq/example"
)

// Runner is a test runner for the example app.
var Runner = dogmatest.New(&example.App{})
