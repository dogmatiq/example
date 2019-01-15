package engine

import (
	"reflect"

	"github.com/dogmatiq/dogma"
)

// MessageClass is an enumeration of the "classes" of message.
type MessageClass string

const (
	// Command is the class for command messages.
	Command MessageClass = "command"

	// Event is the class for event messages.
	Event MessageClass = "event"
)

// MessageComparator is a function that returns true if two messages are equal.
type MessageComparator func(dogma.Message, dogma.Message) bool

// DefaultMessageComparator is the default message comparator. It returns true
// if the messages are deeply equal.
func DefaultMessageComparator(a, b dogma.Message) bool {
	return reflect.DeepEqual(a, b)
}
