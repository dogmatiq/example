package dogmatest

import (
	"context"
	"fmt"
	"reflect"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/examples/dogmatest/compare"
	"github.com/dogmatiq/examples/dogmatest/engine"
	"github.com/dogmatiq/examples/dogmatest/engine/aggregate"
	"github.com/dogmatiq/examples/dogmatest/render"
)

// Engine is a Dogma engine used to test Dogma applications.
type Engine struct {
	controllers []engine.Controller
	classes     map[reflect.Type]engine.MessageClass
	routes      map[reflect.Type][]engine.Controller
	comparator  compare.Comparator
	renderer    render.Renderer
}

// NewEngine returns a new test engine for the given application.
func NewEngine(
	a dogma.App,
	options ...EngineOption,
) *Engine {
	cfg := &engine.Configuration{}

	e := &Engine{
		comparator: compare.DefaultComparator,
		renderer:   render.DefaultRenderer,
	}

	for _, opt := range options {
		opt(e)
	}

	for _, h := range a.Aggregates {
		c := &aggregate.Configurer{
			Handler:  h,
			Renderer: e.renderer,
		}

		h.Configure(c)
		c.Apply(cfg)
	}

	e.controllers = cfg.Controllers()
	e.classes = cfg.Classes()
	e.routes = cfg.Routes()

	return e
}

// Reset clears the state of the engine, and then prepares the engine by
// handling the given messages.
//
// This method should only be used outside the context of a test, that is, when
// a *testing.T is not available. Otherwise, use e.Begin(t).Reset().
func (e *Engine) Reset(ctx context.Context, messages ...dogma.Message) *Engine {
	if err := e.reset(ctx, engine.SilentLogger, messages); err != nil {
		panic(err)
	}

	return e
}

func (e *Engine) reset(
	ctx context.Context,
	logger engine.Logger,
	messages []dogma.Message,
) error {
	for _, c := range e.controllers {
		c.Reset()
	}

	return e.prepare(ctx, logger, messages)
}

// Prepare handles the given messages, without capturing test results.
//
// It is used to place the application into a particular state before handling a
// test message.
//
// This method should only be used outside the context of a test, that is, when
// a *testing.T is not available. Otherwise, use e.Begin(t).Prepare().
func (e *Engine) Prepare(ctx context.Context, messages ...dogma.Message) *Engine {
	if err := e.prepare(ctx, engine.SilentLogger, messages); err != nil {
		panic(err)
	}

	return e
}

func (e *Engine) prepare(
	ctx context.Context,
	logger engine.Logger,
	messages []dogma.Message,
) error {
	queue := make([]*engine.Envelope, 0, len(messages))

	for _, m := range messages {
		t := reflect.TypeOf(m)
		cl, ok := e.classes[t]

		if !ok {
			return fmt.Errorf("no route for messages of type %s", t)
		}

		queue = append(
			queue,
			engine.NewEnvelope(m, cl),
		)
	}

	return e.process(ctx, logger, queue...)
}

// isRoutable returns nil if m is routed to at least one handler as a message of
// the given class.
func (e *Engine) isRoutable(m dogma.Message, cl engine.MessageClass) error {
	t := reflect.TypeOf(m)
	actual, ok := e.classes[t]

	if !ok {
		return fmt.Errorf("no route for messages of type %s", t)
	}

	switch actual {
	case cl:
		return nil
	case engine.Command:
		return fmt.Errorf("messages of type %s are commands, not events", t)
	case engine.Event:
		return fmt.Errorf("messages of type %s are events, not commands", t)
	}

	panic("internal error: unrecognised message class: " + actual)
}

func (e *Engine) process(
	ctx context.Context,
	logger engine.Logger,
	queue ...*engine.Envelope,
) error {
	for len(queue) > 0 {
		env := queue[0]
		queue = queue[1:]

		t := reflect.TypeOf(env.Message)

		for _, c := range e.routes[t] {
			if err := c.Handle(ctx, logger, env); err != nil {
				return err
			}

			queue = append(queue, env.Children...)
		}
	}

	return nil
}
