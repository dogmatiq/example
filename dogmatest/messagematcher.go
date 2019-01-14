package dogmatest

import (
	"fmt"
	"math"
	"reflect"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/examples/dogmatest/internal/types"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// messageMatcher returns a matcher that matches messages equal to m, with the
// given class.
func messageMatcher(
	m dogma.Message,
	cl types.MessageClass,
) Matcher {
	return func(tr TestResult) MatchResult {
		t := reflect.TypeOf(m)
		r := MatchResult{
			Details: tr.Describe(m),
		}

		switch cl {
		case types.Command:
			r.Title = fmt.Sprintf("execute specific '%s' command", t)
			r.Message = "this command was not executed"
		case types.Event:
			r.Title = fmt.Sprintf("record specific '%s' event", t)
			r.Message = "this event was not recorded"
		}

		var (
			bestMatch *types.Envelope
			bestSim   = unrelatedTypes
			bestEqual bool
		)

		tr.Envelope.Walk(
			func(env *types.Envelope) bool {
				// look for an identical message
				if tr.Compare(m, env.Message) {
					// if it's the right message class, we've found our match
					if env.Class == cl {
						r.Passed = true
						r.Message = ""
						return false
					}

					// otherwise, this will always be the best match, assuming we don't find an
					// actual match
					bestMatch = env
					bestSim = sameTypes
					bestEqual = true

					return true
				}

				// check to see if this message is of a similar type to our expected message
				sim := typeSimilarity(
					t,
					reflect.TypeOf(env.Message),
				)

				if sim > bestSim {
					bestMatch = env
					bestSim = sim
				}

				return true
			},
		)

		// we found an exact match, nothing more to do
		if r.Passed {
			return r
		}

		// we found an equal message, but it was the wrong class
		if bestEqual {
			r.Hint = chooseHint(
				cl,
				bestMatch.Class,
				"", // classes are guaranteed to differ, otherwise we passed
				"This message was executed as a command, did you mean to use the Command() matcher instead of Event()?",
				"This message was recorded as an event, did you mean to use the Event() matcher instead of Command()?",
			)

			// bail early so we don't build a diff
			return r
		}

		switch bestSim {
		case unrelatedTypes:
			// we didn't find any message of a similar type
			// bail early so we don't build a diff
			return r

		case sameTypes:
			// we found a message of the same type with different content
			r.Hint = chooseHint(
				cl,
				bestMatch.Class,
				"Check the content of the message.",
				"A similar message was executed as a command, did you mean to use the Command() matcher instead of Event()?",
				"A similar message was recorded as an event, did you mean to use the Event() matcher instead of Command()?",
			)

		default:
			// we found a message of a similar type, but not the same type
			r.Hint = chooseHint(
				cl,
				bestMatch.Class,
				"Check the type of the message.",
				"A message of a similar type was executed as a command, did you mean to use the Command() matcher instead of Event()?",
				"A message of a similar type was recorded as an event, did you mean to use the Event() matcher instead of Command()?",
			)
		}

		// use a diff of the message description in the details field
		diff := diffmatchpatch.New()
		r.Details = diff.DiffPrettyText(
			diff.DiffMain(
				r.Details,
				tr.Describe(bestMatch.Message),
				true,
			),
		)

		return r
	}
}

// messageTypeMatcher returns a matcher that matches messages with the same type
// as m, with the given class.
func messageTypeMatcher(
	m dogma.Message,
	cl types.MessageClass,
) Matcher {
	return func(tr TestResult) MatchResult {
		t := reflect.TypeOf(m)
		r := MatchResult{}

		switch cl {
		case types.Command:
			r.Title = fmt.Sprintf("execute any '%s' command", t)
			r.Message = "no commands of this type were executed"
		case types.Event:
			r.Title = fmt.Sprintf("record any '%s' event", t)
			r.Message = "no events of this type were recorded"
		}

		var (
			bestMatch *types.Envelope
			bestSim   = unrelatedTypes
		)

		tr.Envelope.Walk(
			func(env *types.Envelope) bool {
				mt := reflect.TypeOf(env.Message)
				sim := typeSimilarity(t, mt)

				// look for a message of the expected type
				// if it's the right message class, we've found our match
				if sim == sameTypes && env.Class == cl {
					r.Passed = true
					r.Message = ""
					return false
				}

				if sim > bestSim {
					bestMatch = env
					bestSim = sim
				}

				return true
			},
		)

		// we found an exact match, nothing more to do
		if r.Passed {
			return r
		}

		switch bestSim {
		case unrelatedTypes:
			// we didn't find any message of a similar type
			// bail early so we don't build a diff
			return r

		case sameTypes:
			r.Hint = chooseHint(
				cl,
				bestMatch.Class,
				"", // classes are guaranteed to differ, otherwise we passed
				"A message of this type was executed as a command, did you mean to use the CommandType() matcher instead of EventType()?",
				"A message of this type was recorded as an event, did you mean to use the EventType() matcher instead of CommandType()?",
			)

			// bail early so we don't build a diff
			return r

		default:
			// we found a message of a similar type, but not the same type
			r.Hint = chooseHint(
				cl,
				bestMatch.Class,
				"Check the type of the message.",
				"A message of a similar type was executed as a command, did you mean to use the CommandType() matcher instead of EventType()?",
				"A message of a similar type was recorded as an event, did you mean to use the EventType() matcher instead of CommandType()?",
			)
		}

		// use a diff of the message type in the details field
		diff := diffmatchpatch.New()
		r.Details = diff.DiffPrettyText(
			diff.DiffMain(
				t.String(),
				reflect.TypeOf(bestMatch.Message).String(),
				true,
			),
		)

		return r
	}
}

const (
	sameTypes      uint64 = math.MaxUint64
	unrelatedTypes uint64 = 0
)

// typeSimilarity returns the "similarity" between two related types.
//
// A similarity of sameType indicates the types are identical.
// A similarity of unrelatedTypes indicates the that the types are not related at all.
func typeSimilarity(a, b reflect.Type) (v uint64) {
	v = sameTypes

	if a == b {
		return v
	}

	if n, ok := pointerDistance(a, b); ok {
		return v - n
	}

	if n, ok := pointerDistance(b, a); ok {
		return v - n
	}

	return unrelatedTypes
}

// pointerDistance returns the "distance" from the pointer type p, to the
// elemental type t.
func pointerDistance(p, t reflect.Type) (n uint64, ok bool) {
	for p.Kind() == reflect.Ptr {
		p = p.Elem()
		n++

		if p == t {
			ok = true
			break
		}
	}

	return
}

// chooseHint returns a different string message depending on whether the
// expected and actual message class are the same, or if they are different,
// whether the actual message was a command or an event.
func chooseHint(
	expected, actual types.MessageClass,
	same, command, event string,
) string {
	m := ""

	switch actual {
	case expected:
		m = same
	case types.Command:
		m = command
	case types.Event:
		m = event
	}

	if m == "" {
		panic("internal matcher error: no appropriate hint message")
	}

	return m
}
