package aggregate

import (
	"fmt"

	"github.com/dogmatiq/examples/dogmatest/engine"
)

// log records a log message about the handling of a dogma message by an aggregate message handler.
//
// Note: we perform all logging within "aggregate.go" so that it's more
// obvious in the test logs that the message does not originate in one of the
// user's test files.
func log(logger engine.Logger, name, id string, f string, v ...interface{}) {
	logger.Logf(
		"%s '%s' %s",
		name,
		id,
		fmt.Sprintf(f, v...),
	)
}
