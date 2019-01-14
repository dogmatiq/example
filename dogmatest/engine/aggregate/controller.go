package aggregate

import (
	"context"
	"fmt"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/examples/dogmatest/engine"
)

type controller struct {
	name      string
	handler   dogma.AggregateMessageHandler
	instances map[string]dogma.AggregateRoot
	describe  engine.MessageDescriber
}

func (c *controller) Name() string {
	return c.name
}

func (c *controller) Handler() interface{} {
	return c.handler
}

func (c *controller) Handle(
	_ context.Context,
	logger engine.Logger,
	env *engine.Envelope,
) error {
	id := c.handler.RouteCommandToInstance(env.Message)
	if id == "" {
		return fmt.Errorf("aggregate '%s' instance ID must not be empty", c.name)
	}

	r, ok := c.instances[id]

	if ok {
		log(logger, c.name, id, "already exists")
	} else {
		log(logger, c.name, id, "does not exist")

		r = c.handler.New()

		if r == nil {
			return fmt.Errorf("aggregate '%s' root must not be nil", c.name)
		}
	}

	s := &scope{
		id:       id,
		name:     c.name,
		root:     r,
		exists:   ok,
		command:  env,
		describe: c.describe,
		logger:   logger,
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

	return nil
}

func (c *controller) Reset() {
	c.instances = nil
}
