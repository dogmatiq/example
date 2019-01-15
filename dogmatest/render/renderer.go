package render

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/dogmatiq/dogma"
)

// Renderer is an interface for rendering various Dogma values.
type Renderer interface {
	RenderMessage(m dogma.Message) string
	RenderAggregateRoot(r dogma.AggregateRoot) string
	RenderProcessRoot(r dogma.ProcessRoot) string
	RenderAggregateMessageHandler(h dogma.ProcessMessageHandler) string
	RenderProcessMessageHandler(h dogma.ProcessMessageHandler) string
	RenderIntegrationMessageHandler(h dogma.IntegrationMessageHandler) string
	RenderProjectionMessageHandler(h dogma.ProjectionMessageHandler) string
}

// DefaultRenderer is the default Renderer implementation
var DefaultRenderer Renderer = defaultRenderer{
	spew: spew.ConfigState{
		Indent:                  NestedIndentPrefix,
		DisableMethods:          true,
		DisablePointerAddresses: true,
		DisableCapacities:       true,
		SortKeys:                true,
	},
}

type defaultRenderer struct {
	spew spew.ConfigState
}

func (r defaultRenderer) RenderMessage(m dogma.Message) string {
	return r.render(m)
}

func (r defaultRenderer) RenderAggregateRoot(ar dogma.AggregateRoot) string {
	return r.render(ar)
}

func (r defaultRenderer) RenderProcessRoot(pr dogma.ProcessRoot) string {
	return r.render(pr)
}

func (r defaultRenderer) RenderAggregateMessageHandler(h dogma.ProcessMessageHandler) string {
	return r.render(h)
}

func (r defaultRenderer) RenderProcessMessageHandler(h dogma.ProcessMessageHandler) string {
	return r.render(h)
}

func (r defaultRenderer) RenderIntegrationMessageHandler(h dogma.IntegrationMessageHandler) string {
	return r.render(h)
}

func (r defaultRenderer) RenderProjectionMessageHandler(h dogma.ProjectionMessageHandler) string {
	return r.render(h)
}

func (r defaultRenderer) render(v interface{}) string {
	return pretty(v)
}
