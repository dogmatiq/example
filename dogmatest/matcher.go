package dogmatest

import (
	"fmt"

	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/examples/dogmatest/engine"
)

// Matcher is a predicate function that checks a test result against some criteria.
type Matcher func(TestResult) MatchResult

// Command is a matcher that passes only if a message equal to m is executed as
// a command.
func Command(m dogma.Message) Matcher {
	return messageMatcher(m, engine.Command)
}

// CommandType is a matcher that passes only if a message of the same type as m is
// executed as a command.
func CommandType(m dogma.Message) Matcher {
	return messageTypeMatcher(m, engine.Command)
}

// Event is a matcher that passes only if a message equal to m is recorded as
// an event.
func Event(m dogma.Message) Matcher {
	return messageMatcher(m, engine.Event)
}

// EventType is a matcher that passes only if a message of the same type as m is
// recorded as an event.
func EventType(m dogma.Message) Matcher {
	return messageTypeMatcher(m, engine.Event)
}

// All is matcher that passes only if all of the given sub-matchers pass.
func All(matchers ...Matcher) Matcher {
	if m, ok := flattenMatchers(matchers); ok {
		return m
	}

	return logicalMatcher(
		"all of",
		matchers,
		func(n int) (string, bool) {
			c := len(matchers)
			return fmt.Sprintf("%d", c), n == c
		},
	)
}

// Any is matcher that passes if any of the given sub-matchers pass.
func Any(matchers ...Matcher) Matcher {
	if m, ok := flattenMatchers(matchers); ok {
		return m
	}

	return logicalMatcher(
		"any of",
		matchers,
		func(n int) (string, bool) {
			return ">0", n > 0
		},
	)
}

// Not is matcher that passes only if none of the given sub-matchers pass.
func Not(matchers ...Matcher) Matcher {
	flattenMatchers(matchers) // use to verify there are >0 matchers

	return logicalMatcher(
		"none of",
		matchers,
		func(n int) (string, bool) {
			return "0", n == 0
		},
	)
}

// flattenMatchers returns matches[0] if there is only one matcher.
func flattenMatchers(matchers []Matcher) (Matcher, bool) {
	switch len(matchers) {
	case 0:
		panic("no sub-matchers provided")
	case 1:
		return matchers[0], true
	default:
		return nil, false
	}
}
