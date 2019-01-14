package dogmatest

import "github.com/dogmatiq/examples/dogmatest/engine"

// EngineOption is an option that configures the test engine.
type EngineOption func(*Engine)

// UseMessageComparator is an engine option that specifies the message
// comparator to use.
func UseMessageComparator(c engine.MessageComparator) EngineOption {
	return func(e *Engine) {
		e.compare = c
	}
}

// UseMessageDescriber is an engine option that specifies the message
// describer to use.
func UseMessageDescriber(d engine.MessageDescriber) EngineOption {
	return func(e *Engine) {
		e.describe = d
	}
}
