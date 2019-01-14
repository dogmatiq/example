package dogmatest

import (
	"io"
	"strings"

	"github.com/dogmatiq/examples/dogmatest/internal/ioutil"
)

// MatchResult represents the result of performing a match.
type MatchResult struct {
	Title    string
	Passed   bool
	Message  string
	Hint     string
	Details  string
	Children []MatchResult
}

// Append adds a child result to r.
func (r *MatchResult) Append(c MatchResult) {
	r.Children = append(r.Children, c)
}

// WriteTo writes a human-readable report on the match result to w.
func (r *MatchResult) WriteTo(w io.Writer) (n int64, err error) {
	defer ioutil.Recover(&err)

	if r.Passed {
		ioutil.MustWriteString(w, "\x1b[32m✓ ")
	} else {
		ioutil.MustWriteString(w, "\x1b[31m✗ ")
	}

	ioutil.MustWriteString(w, r.Title)
	ioutil.MustWriteString(w, "\x1b[0m")

	if r.Message != "" {
		ioutil.MustWriteString(w, " (")
		ioutil.MustWriteString(w, r.Message)
		ioutil.MustWriteString(w, ")")
	}

	ioutil.MustWriteString(w, "\n")

	{
		dw := ioutil.NewIndenter(w, "  | ")

		if r.Details != "" {
			ioutil.MustWriteString(dw, "\n")
			ioutil.MustWriteString(dw, strings.TrimSpace(r.Details))
			ioutil.MustWriteString(dw, "\n")
		}

		if r.Hint != "" {
			ioutil.MustWriteString(dw, "\n")
			ioutil.MustWriteString(dw, "Hint: ")
			ioutil.MustWriteString(dw, r.Hint)
			ioutil.MustWriteString(dw, "\n")
		}
	}

	if len(r.Children) > 0 {
		cw := ioutil.NewIndenter(w, "")
		for _, rc := range r.Children {
			n += ioutil.MustWriteTo(cw, &rc)
		}
	}

	return
}
