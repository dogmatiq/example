package dogmatest

// EngineOption is an option that configures the test engine.
type EngineOption func(*Engine)

// UseMessageComparator is an engine option that specifies the message
// comparator to use.
func UseMessageComparator(c MessageComparator) EngineOption {
	return func(e *Engine) {
		e.compare = c
	}
}

// UseMessageDescriber is an engine option that specifies the message
// describer to use.
func UseMessageDescriber(d MessageDescriber) EngineOption {
	return func(e *Engine) {
		e.describe = d
	}
}
