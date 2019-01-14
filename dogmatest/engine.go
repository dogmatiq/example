package dogmatest

import (
	"fmt"
	"reflect"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/examples/dogmatest/internal/aggregate"
	"github.com/dogmatiq/examples/dogmatest/internal/types"
)

// Engine is a Dogma engine used to test Dogma applications.
type Engine struct {
	controllers []types.Controller
	classes     map[reflect.Type]types.MessageClass
	routes      map[reflect.Type][]types.Controller
	compare     MessageComparator
	describe    MessageDescriber
}

// NewEngine returns a new test engine for the given application.
func NewEngine(
	a dogma.App,
	options ...EngineOption,
) *Engine {
	cfg := &types.Configuration{}

	for _, h := range a.Aggregates {
		c := &aggregate.Configurer{
			Handler: h,
		}
		h.Configure(c)
		c.Apply(cfg)
	}

	e := &Engine{
		controllers: cfg.Controllers(),
		classes:     cfg.Classes(),
		routes:      cfg.Routes(),
		compare:     DefaultMessageComparator,
		describe:    DefaultMessageDescriber,
	}

	for _, o := range options {
		o(e)
	}

	return e
}

// Reset clears the state of the engine, and then prepares the engine by
// handling the given messages.
func (e *Engine) Reset(messages ...dogma.Message) *Engine {
	for _, c := range e.controllers {
		c.Reset()
	}

	return e.Prepare(messages...)
}

// Prepare handles the given messages, without capturing test results.
//
// It is used to place the application into a particular state before handling a
// test message.
func (e *Engine) Prepare(messages ...dogma.Message) *Engine {
	queue := make([]*types.Envelope, 0, len(messages))

	for _, m := range messages {
		t := reflect.TypeOf(m)
		cl, ok := e.classes[t]

		if !ok {
			panic(fmt.Sprintf("no route for messages of type %s", t))
		}

		queue = append(
			queue,
			types.NewEnvelope(m, cl),
		)
	}

	e.do(queue...)

	return e
}

// TestCommand captures test results describing how the application handles the
// command m.
func (e *Engine) TestCommand(t TestingT, m dogma.Message) TestResult {
	e.assertIsRoutableCommand(m)

	return e.test(
		t,
		types.NewEnvelope(m, types.Command),
	)
}

// TestEvent captures test results describing how the application handles the
// event m.
func (e *Engine) TestEvent(t TestingT, m dogma.Message) TestResult {
	e.assertIsRoutableEvent(m)

	return e.test(
		t,
		types.NewEnvelope(m, types.Event),
	)
}

func (e *Engine) assertIsRoutableCommand(m dogma.Message) {
	t := reflect.TypeOf(m)
	c, ok := e.classes[t]

	if !ok {
		panic(fmt.Sprintf("no route for commands of type %s", t))
	}

	if c == types.Event {
		panic(fmt.Sprintf("messages of type %s are events, not commands", t))
	}
}

func (e *Engine) assertIsRoutableEvent(m dogma.Message) {
	t := reflect.TypeOf(m)
	c, ok := e.classes[t]

	if !ok {
		panic(fmt.Sprintf("no route for events of type %s", t))
	}

	if c == types.Command {
		panic(fmt.Sprintf("messages of type %s are commands, not events", t))
	}
}

func (e *Engine) do(queue ...*types.Envelope) {
	for len(queue) > 0 {
		env := queue[0]
		queue = queue[1:]

		t := reflect.TypeOf(env.Message)

		for _, c := range e.routes[t] {
			c.Handle(env)
			queue = append(queue, env.Children...)
		}
	}
}

func (e *Engine) test(t TestingT, env *types.Envelope) TestResult {
	tr := TestResult{
		T:        t,
		Envelope: env,
		Compare:  e.compare,
		Describe: e.describe,
	}

	e.do(tr.Envelope)

	return tr
}
