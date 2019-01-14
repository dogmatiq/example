package aggregate

import (
	"fmt"
	"reflect"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/examples/dogmatest/engine"
)

// Configurer is the test engine's implementation of dogma.AggregateConfigurer.8
type Configurer struct {
	Handler  dogma.AggregateMessageHandler
	Describe engine.MessageDescriber

	name     string
	commands map[reflect.Type]struct{}
}

// Name sets the name of the handler. Each handler within an application must
// have a unique name.
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

// RouteCommandType configures the engine to route domain command messages of
// the same type as m to the handler.
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

// Apply applies the configurer's properties to the given configuration.
func (c *Configurer) Apply(cfg *engine.Configuration) {
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
		name:     c.name,
		handler:  c.Handler,
		describe: c.Describe,
	}

	cfg.RegisterController(ctrl)

	for t := range c.commands {
		cfg.RouteCommand(t, ctrl)
	}
}
