package aggregate

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/examples/dogmatest/internal/types"
)

type controller struct {
	name      string
	handler   dogma.AggregateMessageHandler
	instances map[string]dogma.AggregateRoot
}

func (c *controller) Name() string {
	return c.name
}

func (c *controller) Handler() interface{} {
	return c.handler
}

func (c *controller) Handle(env types.Envelope) []types.Envelope {
	id := c.handler.RouteCommandToInstance(env.Message)
	if id == "" {
		panic("aggregate instances ID must not be empty")
	}

	r, ok := c.instances[id]

	if !ok {
		r = c.handler.New()
	}

	s := &scope{
		id:      id,
		root:    r,
		exists:  ok,
		command: env,
	}

	c.handler.HandleCommand(s, env.Message)

	if s.exists {
		if c.instances == nil {
			c.instances = map[string]dogma.AggregateRoot{}
		}
		c.instances[id] = s.root
	} else {
		delete(c.instances, id)
	}

	return s.events
}

func (c *controller) Reset() {
	c.instances = nil
}
