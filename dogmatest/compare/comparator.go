package compare

import (
	"reflect"

	"github.com/davecgh/go-spew/spew"
	"github.com/dogmatiq/dogma"
)

// Comparator is an interface for comparing various Dogma values for equality.
type Comparator interface {
	CompareMessage(a, b dogma.Message) bool
	CompareAggregateRoot(a, b dogma.AggregateRoot) bool
	CompareProcessRoot(a, b dogma.ProcessRoot) bool
}

// DefaultComparator is the default Comparator implementation.
var DefaultComparator Comparator = defaultComparator{}

type defaultComparator struct {
	spew spew.ConfigState
}

func (r defaultComparator) CompareMessage(a, b dogma.Message) bool {
	return reflect.DeepEqual(a, b)
}

func (r defaultComparator) CompareAggregateRoot(a, b dogma.AggregateRoot) bool {
	return reflect.DeepEqual(a, b)
}

func (r defaultComparator) CompareProcessRoot(a, b dogma.ProcessRoot) bool {
	return reflect.DeepEqual(a, b)
}
