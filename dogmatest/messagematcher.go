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
			candidate types.Envelope
			best      uint64
		)

		for _, env := range tr.Output {
			if tr.Compare(m, env.Message) {
				if env.Class == cl {
					r.Passed = true
					r.Message = ""
				} else if env.Class == types.Command {
					r.Hint = "an identical message was executed as a command, are you using the correct matcher?"
				} else {
					r.Hint = "an identical message was recorded as an event, are you using the correct matcher?"
				}

				return r
			}

			sim := typeSimilarity(
				t,
				reflect.TypeOf(env.Message),
			)

			if sim > best {
				candidate = env
				best = sim
			}
		}

		switch best {
		case 0:
			return r
		case math.MaxUint64:
			if candidate.Class == cl {
				r.Hint = "is the message content correct?"
			} else {
				switch candidate.Class {
				case types.Command:
					r.Hint = "a similar message was executed as a command, are you using the correct matcher?"
				case types.Event:
					r.Hint = "a similar message was recorded as an event, are you using the correct matcher?"
				}
			}
		default:
			if candidate.Class == cl {
				r.Hint = "is there a type mismatch?"
			} else {
				switch candidate.Class {
				case types.Command:
					r.Hint = "a message of a similar type was executed as a command, are you using the correct matcher?"
				case types.Event:
					r.Hint = "a message of a similar type was recorded as an event, are you using the correct matcher?"
				}
			}
		}

		diff := diffmatchpatch.New()
		r.Details = diff.DiffPrettyText(
			diff.DiffMain(
				r.Details,
				tr.Describe(candidate.Message),
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
			candidate types.Envelope
			best      uint64
		)

		for _, env := range tr.Output {
			mt := reflect.TypeOf(env.Message)
			sim := typeSimilarity(t, mt)

			if sim == math.MaxUint64 {
				if env.Class == cl {
					r.Passed = true
					r.Message = ""
				} else if env.Class == types.Command {
					r.Hint = "a message of this type was executed as a command, are you using the correct matcher?"
				} else {
					r.Hint = "a message of this type was recorded as an event, are you using the correct matcher?"
				}

				return r
			}

			if sim > best {
				candidate = env
				best = sim
			}
		}

		if best == 0 {
			return r
		}

		if candidate.Class == cl {
			r.Hint = "is there a type mismatch?"
		} else {
			switch candidate.Class {
			case types.Command:
				r.Hint = "a message of a similar type was executed as a command, are you using the correct matcher?"
			case types.Event:
				r.Hint = "a message of a similar type was recorded as an event, are you using the correct matcher?"
			}
		}

		diff := diffmatchpatch.New()

		r.Details = diff.DiffPrettyText(
			diff.DiffMain(
				t.String(),
				reflect.TypeOf(candidate.Message).String(),
				true,
			),
		)

		return r
	}
}

// typeSimilarity returns the "similarity" between two related types.
//
// A similarity of math.Uint64Max indicates the types are identical.
// A similarity of 0 indicates the that the types are not related at all.
func typeSimilarity(a, b reflect.Type) (v uint64) {
	v = math.MaxUint64

	if a == b {
		return v
	}

	if n, ok := pointerDistance(a, b); ok {
		return v - n
	}

	if n, ok := pointerDistance(b, a); ok {
		return v - n
	}

	return 0
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
