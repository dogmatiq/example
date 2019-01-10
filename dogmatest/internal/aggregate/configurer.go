package aggregate

import (
	"fmt"
	"reflect"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/examples/dogmatest/internal/types"
)

type Configurer struct {
	Handler  dogma.AggregateMessageHandler
	name     string
	commands map[reflect.Type]struct{}
}

func (c *Configurer) Name(n string) {
	if c.name != "" {
		panic(fmt.Sprintf(
			"this handler is already named %s",
			c.name,
		))
	}

	if n == "" {
		panic("handler names can not be empty")
	}

	c.name = n
}

func (c *Configurer) RouteCommandType(m dogma.Message) {
	if c.commands == nil {
		c.commands = map[reflect.Type]struct{}{}
	}

	t := reflect.TypeOf(m)

	if _, ok := c.commands[t]; ok {
		panic(fmt.Sprintf("commands of type %s already routed to this handler", t))
	}

	c.commands[t] = struct{}{}
}

func (c *Configurer) Apply(cfg *types.Configuration) {
	if c.name == "" {
		panic(fmt.Sprintf(
			"%#v did not call AggregateConfigurer.Name() in Configure()",
			c.Handler,
		))
	}

	if len(c.commands) == 0 {
		panic(fmt.Sprintf(
			"%#v did not call AggregateConfigurer.RouteCommandType() in Configure()",
			c.Handler,
		))
	}

	ctrl := &controller{
		name:    c.name,
		handler: c.Handler,
	}

	cfg.RegisterController(ctrl)

	for t := range c.commands {
		cfg.RouteCommand(t, ctrl)
	}
}
