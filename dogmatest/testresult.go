package dogmatest

import (
	"strings"

	"github.com/dogmatiq/examples/dogmatest/internal/types"
)

// TestResult holds information about a message that was tested.
type TestResult struct {
	// T is the *testing.T value, or equivalent under which the test was performed.
	T TestingT

	// Envelope is the message envelope that contains the message under test.
	Envelope *types.Envelope

	// Compare is the comparator used to test messages for equality.
	Compare MessageComparator

	// Describe is the describer used to render a human-readable representation of
	// a message.
	Describe MessageDescriber
}

// Expect runs the given set of matchers against this test result, and fails the
// test if any of them do not pass.
func (r TestResult) Expect(matchers ...Matcher) {
	mr := All(matchers...)(r)

	var b strings.Builder

	if mr.Passed {
		b.WriteString("expectation passed:\n")
	} else {
		b.WriteString("expectation failed:\n")
	}

	mr.WriteTo(&b)
	r.T.Log(b.String())

	if !mr.Passed {
		r.T.FailNow()
	}
}

// TestingT is the interface by which matchers make use of Go's *testing.T type.
//
// It allows use of stand-ins, such as Ginkgo's GinkgoT() value.
type TestingT interface {
	FailNow()
	Log(args ...interface{})
	Logf(format string, args ...interface{})
}
