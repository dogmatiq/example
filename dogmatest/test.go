package dogmatest

import (
	"context"
	"reflect"
	"strings"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/examples/dogmatest/engine"
	"github.com/dogmatiq/examples/dogmatest/internal/ioutil"
)

// Test is an interface for testing a message.
type Test struct {
	ctx    context.Context
	t      TestingT
	engine *Engine
}

// Begin starts a new test.
func Begin(t TestingT, e *Engine) Test {
	return BeginContext(
		context.Background(),
		t,
		e,
	)
}

// BeginContext starts a new test with a context.
func BeginContext(ctx context.Context, t TestingT, e *Engine) Test {
	return Test{
		ctx:    ctx,
		t:      t,
		engine: e,
	}
}

// Reset clears the state of the engine, and then prepares the engine by
// handling the given messages.
func (t Test) Reset(messages ...dogma.Message) Test {
	if err := t.engine.reset(t.ctx, t.t, messages); err != nil {
		t.t.Fatal(err)
	}

	return t
}

// Prepare handles the given messages, without capturing test results.
//
// It is used to place the application into a particular state before handling a
// test message.
func (t Test) Prepare(messages ...dogma.Message) Test {
	if err := t.engine.prepare(t.ctx, t.t, messages); err != nil {
		t.t.Fatal(err)
	}

	return t
}

// TestCommand captures test results describing how the application handles the
// command m.
func (t Test) TestCommand(m dogma.Message) TestResult {
	if err := t.engine.isRoutable(m, engine.Command); err != nil {
		t.t.Fatal(err)
	}

	return t.test(
		engine.NewEnvelope(m, engine.Command),
	)
}

// TestEvent captures test results describing how the application handles the
// event m.
func (t Test) TestEvent(m dogma.Message) TestResult {
	if err := t.engine.isRoutable(m, engine.Event); err != nil {
		t.t.Fatal(err)
	}

	return t.test(
		engine.NewEnvelope(m, engine.Event),
	)
}

func (t Test) test(env *engine.Envelope) TestResult {
	tr := TestResult{
		T:        t.t,
		Envelope: env,
		Compare:  t.engine.compare,
		Describe: t.engine.describe,
	}

	t.t.Logf(
		"testing '%s' %s:\n\n%s\n",
		reflect.TypeOf(env.Message),
		env.Class,
		ioutil.Indent(tr.Describe(env.Message), "| "),
	)

	if err := t.engine.process(t.ctx, t.t, tr.Envelope); err != nil {
		t.t.Fatal(err)
	}

	return tr
}

// TestResult holds information about a message that was tested.
type TestResult struct {
	// T is the *testing.T value, or equivalent under which the test was performed.
	T TestingT

	// Envelope is the message envelope that contains the message under test.
	Envelope *engine.Envelope

	// Compare is the comparator used to test messages for equality.
	Compare engine.MessageComparator

	// Describe is the describer used to render a human-readable representation of
	// a message.
	Describe engine.MessageDescriber
}

// Expect runs the given set of matchers against this test result, and fails the
// test if any of them do not pass.
func (r TestResult) Expect(matchers ...Matcher) {
	mr := All(matchers...)(r)

	var b strings.Builder

	if mr.Passed {
		b.WriteString("expectation passed:\n\n")
	} else {
		b.WriteString("expectation failed:\n\n")
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
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}
