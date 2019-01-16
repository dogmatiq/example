package dogmatest

import (
	"io"
	"strings"

	"github.com/dogmatiq/examples/dogmatest/render"
	"github.com/dogmatiq/iago"
	"github.com/dogmatiq/iago/indent"
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
func (r *MatchResult) WriteTo(w io.Writer) (_ int64, err error) {
	defer iago.Recover(&err)

	n := 0

	if r.Passed {
		n += iago.MustWriteString(w, render.Green)
		n += iago.MustWriteString(w, "✓ ")
	} else {
		n += iago.MustWriteString(w, render.Red)
		n += iago.MustWriteString(w, "✗ ")
	}

	n += iago.MustWriteString(w, r.Title)
	n += iago.MustWriteString(w, render.Reset)

	if r.Message != "" {
		n += iago.MustWriteString(w, " (")
		n += iago.MustWriteString(w, r.Message)
		n += iago.MustWriteString(w, ")")
	}

	iago.MustWriteString(w, "\n")
	if r.Details != "" || r.Hint != "" {
		n += iago.MustWriteString(w, "\n")
	}

	{
		dw := indent.NewIndenter(w, []byte("  | "))

		if r.Details != "" {
			n += iago.MustWriteString(dw, strings.TrimSpace(r.Details))
			n += iago.MustWriteString(dw, "\n")
		}

		if r.Hint != "" {
			if r.Details != "" {
				n += iago.MustWriteString(dw, "\n")
			}

			n += iago.MustWriteString(dw, "Hint: ")
			n += iago.MustWriteString(dw, r.Hint)
			n += iago.MustWriteString(dw, "\n")
		}
	}

	if len(r.Children) > 0 {
		cw := indent.NewIndenter(w, nil)
		for _, rc := range r.Children {
			n += iago.MustWriteTo(cw, &rc)
		}
	}

	return int64(n), nil
}
