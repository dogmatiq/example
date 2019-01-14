package engine

import (
	"reflect"

	"github.com/davecgh/go-spew/spew"
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

// MessageDescriber is a function that produces human-readable a description of m.
type MessageDescriber func(m dogma.Message) string

// DefaultMessageDescriber is the default message describer.
func DefaultMessageDescriber(m dogma.Message) string {
	return spewConfig.Sdump(m)
}

var spewConfig = spew.ConfigState{
	Indent:                  "    ", // match the ioutil.Indenter default indent
	DisableMethods:          true,
	DisablePointerAddresses: true,
	DisableCapacities:       true,
	SortKeys:                true,
}
