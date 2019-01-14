package app_test

import (
	"os"
	"testing"

	. "github.com/dogmatiq/examples/cmd/bank/internal/app"
	"github.com/dogmatiq/examples/dogmatest"
)

var engine *dogmatest.Engine

func TestMain(m *testing.M) {
	engine = dogmatest.NewEngine(App)
	os.Exit(m.Run())
}
