package dogmatest

import (
	"github.com/dogmatiq/examples/dogmatest/compare"
	"github.com/dogmatiq/examples/dogmatest/render"
)

// EngineOption is an option that configures the test engine.
type EngineOption func(*Engine)

// Comparator is an engine option that specifies the comparator to use.
func Comparator(c compare.Comparator) EngineOption {
	return func(e *Engine) {
		e.comparator = c
	}
}

// Renderer is an engine option that specifies the renderer to use.
func Renderer(r render.Renderer) EngineOption {
	return func(e *Engine) {
		e.renderer = r
	}
}
