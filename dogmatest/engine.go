package dogmatest

import (
	"context"
	"fmt"
	"reflect"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/examples/dogmatest/internal/aggregate"
	"github.com/dogmatiq/examples/dogmatest/internal/types"
)

type EngineOption func(*Engine)

type MessageComparator func(dogma.Message, dogma.Message) bool

func UseMessageComparator(c MessageComparator) EngineOption {
	return func(e *Engine) {
		e.isEqual = c
	}
}

func New(
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
		classes: cfg.Classes(),
		routes:  cfg.Routes(),
		isEqual: func(a, b dogma.Message) bool {
			return reflect.DeepEqual(a, b)
		},
	}

	for _, o := range options {
		o(e)
	}

	return e
}

type Engine struct {
	controllers []types.Controller
	classes     map[reflect.Type]types.MessageClass
	routes      map[reflect.Type][]types.Controller
	isEqual     MessageComparator
}

func (e *Engine) Reset(commands ...dogma.Message) {
	for _, c := range e.controllers {
		c.Reset()
	}
}

func (e *Engine) ExecuteCommand(_ context.Context, m dogma.Message) error {
	e.assertIsRoutableCommand(m)
	e.do(m)
	return nil
}

func (e *Engine) RecordEvent(_ context.Context, m dogma.Message) error {
	e.assertIsRoutableEvent(m)
	e.do(m)
	return nil
}

func (e *Engine) TestCommand(m dogma.Message) TestResult {
	e.assertIsRoutableCommand(m)

	return TestResult{
		Envelopes: e.do(m),
		IsEqual:   e.isEqual,
	}
}

func (e *Engine) assertIsRoutableCommand(m dogma.Message) {
	t := reflect.TypeOf(m)
	c, ok := e.classes[t]

	if !ok {
		panic(fmt.Sprintf("no route for commands of type %s", t))
	}

	if c == types.EventClass {
		panic(fmt.Sprintf("messages of type %s are events, not commands", t))
	}
}

func (e *Engine) assertIsRoutableEvent(m dogma.Message) {
	t := reflect.TypeOf(m)
	c, ok := e.classes[t]

	if !ok {
		panic(fmt.Sprintf("no route for events of type %s", t))
	}

	if c == types.CommandClass {
		panic(fmt.Sprintf("messages of type %s are commands, not events", t))
	}
}

func (e *Engine) do(m dogma.Message) []types.Envelope {
	queue := []types.Envelope{
		types.NewEnvelope(m, types.CommandClass),
	}

	for i := 0; i < len(queue); i++ {
		env := queue[i]
		t := reflect.TypeOf(env.Message)

		for _, c := range e.routes[t] {
			queue = append(
				queue,
				c.Handle(env)...,
			)
		}
	}

	return queue[1:] // don't include the initial message in the results
}

type TestResult struct {
	Envelopes []types.Envelope
	IsEqual   MessageComparator
}

func (r TestResult) ExpectEvents(events ...dogma.Message) {
	actual := r.filter(types.EventClass)

next:
	for _, m := range events {
		for k, x := range actual {
			if r.IsEqual(m, x) {
				delete(actual, k)
				continue next
			}
		}

		panic(fmt.Sprintf(
			"an expected event was not recorded: %#v",
			m,
		))
	}
}

func (r TestResult) ExpectExactEvents(events ...dogma.Message) {
	actual := r.filter(types.EventClass)

next:
	for _, m := range events {
		for id, x := range actual {
			if r.IsEqual(m, x) {
				delete(actual, id)
				continue next
			}
		}

		panic(fmt.Sprintf(
			"an expected event was not recorded: %#v",
			m,
		))
	}

	for _, x := range actual {
		panic(fmt.Sprintf(
			"an unexpected event was recorded: %#v",
			x,
		))
	}
}

func (r TestResult) ExpectNoEvents() {
	r.ExpectExactEvents()
}

func (r TestResult) filter(c types.MessageClass) map[uint64]dogma.Message {
	messages := make(map[uint64]dogma.Message, len(r.Envelopes))

	for _, env := range r.Envelopes {
		if env.Class == c {
			messages[env.MessageID] = env.Message
		}
	}

	return messages
}

func (r TestResult) has(m dogma.Message, c types.MessageClass) bool {
	for _, env := range r.Envelopes {
		if env.Class != c {
			continue
		}

		if r.IsEqual(m, env.Message) {
			return true
		}
	}

	return false
}
