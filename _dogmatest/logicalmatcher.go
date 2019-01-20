package dogmatest

import (
	"fmt"
)

func logicalMatcher(
	title string,
	matchers []Matcher,
	pred func(n int) (string, bool),
) Matcher {
	return func(tr TestResult) MatchResult {
		r := MatchResult{
			Title: title,
		}

		n := 0

		for _, cm := range matchers {
			cr := cm(tr)
			r.Append(cr)

			if cr.Passed {
				n++
			}
		}

		expect, pass := pred(n)

		r.Passed = pass
		r.Message = fmt.Sprintf(
			"%d of %d sub-matchers passed, expected %s",
			n,
			len(matchers),
			expect,
		)

		return r
	}
}
