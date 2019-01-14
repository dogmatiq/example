package types

import (
	"github.com/dogmatiq/dogma"
)

// Envelope is a container for a message that is processed by the test engine.
type Envelope struct {
	Message  dogma.Message
	Class    MessageClass
	Children []*Envelope
}

// NewEnvelope constructs a new envelope containing the given message.
func NewEnvelope(m dogma.Message, c MessageClass) *Envelope {
	return &Envelope{
		Message: m,
		Class:   c,
	}
}

// NewChild constructs a new envelope as a child of e, indicating that m is
// caused by e.Message.
func (e *Envelope) NewChild(m dogma.Message, c MessageClass) *Envelope {
	env := &Envelope{
		Message: m,
		Class:   c,
	}

	e.Children = append(e.Children, env)

	return env
}

// Walk calls fn() for each of this envelope's children, and their children
// recursively.
//
// fn is not called with e itself.
//
// If fn returns false, iteration stops. It returns true if iteration completes
// fully.
func (e *Envelope) Walk(fn func(*Envelope) bool) bool {
	for _, env := range e.Children {
		if !fn(env) {
			return false
		}

		if !env.Walk(fn) {
			return false
		}
	}

	return true
}

// MessageClass is an enumeration of the "classes" of message.
type MessageClass string

const (
	// Command is the class for command messages.
	Command MessageClass = "command"

	// Event is the class for event messages.
	Event MessageClass = "event"
)
