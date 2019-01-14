package aggregate

import (
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/examples/dogmatest/internal/types"
)

type scope struct {
	id      string
	root    dogma.AggregateRoot
	exists  bool
	command *types.Envelope
}

func (s *scope) InstanceID() string {
	return s.id
}

func (s *scope) Create() bool {
	if s.exists {
		return false
	}

	s.exists = true
	return true
}

func (s *scope) Destroy() {
	if !s.exists {
		panic("can not destroy non-existent instance")
	}

	s.exists = false
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
	s.command.NewChild(m, types.Event)
}

func (s *scope) Log(f string, v ...interface{}) {
	// TODO:
}
