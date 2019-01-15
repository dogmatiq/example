package aggregate

import (
	"fmt"
	"reflect"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/examples/dogmatest/engine"
	"github.com/dogmatiq/examples/dogmatest/render"
)

type scope struct {
	id       string
	name     string
	root     dogma.AggregateRoot
	exists   bool
	command  *engine.Envelope
	renderer render.Renderer
	logger   engine.Logger
}

func (s *scope) InstanceID() string {
	return s.id
}

func (s *scope) Create() bool {
	if s.exists {
		return false
	}

	s.exists = true
	log(s.logger, s.name, s.id, "created")

	return true
}

func (s *scope) Destroy() {
	if !s.exists {
		panic("can not destroy non-existent instance")
	}

	s.exists = false
	log(s.logger, s.name, s.id, "destroyed")
}

func (s *scope) Root() dogma.AggregateRoot {
	if !s.exists {
		panic("can not access aggregate root of non-existent instance")
	}

	return s.root
}

func (s *scope) RecordEvent(m dogma.Message) {
	if !s.exists {
		panic("can not record event against non-existent instance")
	}

	s.root.ApplyEvent(m)
	s.command.NewChild(m, engine.Event)

	log(
		s.logger,
		s.name,
		s.id,
		fmt.Sprintf(
			"recorded '%s' event:\n\n%s\n\n",
			reflect.TypeOf(m),
			render.IndentDetails(
				s.renderer.RenderMessage(m),
			),
		),
	)
}

func (s *scope) Log(f string, v ...interface{}) {
	log(
		s.logger,
		s.name,
		s.id,
		"logged: "+fmt.Sprintf(f, v...),
	)
}
