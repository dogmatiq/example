package types

import (
	"sync/atomic"

	"github.com/dogmatiq/dogma"
)

type Envelope struct {
	MessageID     uint64
	CausationID   uint64
	CorrelationID uint64
	Message       dogma.Message
	Class         MessageClass
}

var messageID uint64 // atomic

func NewEnvelope(m dogma.Message, c MessageClass) Envelope {
	id := atomic.AddUint64(&messageID, 1)
	e := Envelope{id, id, id, m, c}
	return e
}

func (e Envelope) NewChild(m dogma.Message, c MessageClass) Envelope {
	id := atomic.AddUint64(&messageID, 1)
	return Envelope{id, e.CausationID, e.CorrelationID, m, c}
}

// MessageClass is an enumeration of the "classes" of message.
type MessageClass bool

const (
	// CommandClass is the class for command messages.
	CommandClass MessageClass = true

	// EventClass is the class for event messages.
	EventClass MessageClass = false
)
