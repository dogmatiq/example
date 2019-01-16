package render

import (
	"github.com/dogmatiq/iago/indent"
)

const (
	// NestedIndentPrefix is the indent prefix to use for nesting.
	NestedIndentPrefix = "\t"

	// DetailsIndentPrefix is the indent prefix to use when rendering additional
	// details in a test.
	DetailsIndentPrefix = "| "
)

// IndentNested indents s using the NestedIndentPrefix.
func IndentNested(s string) string {
	return indent.String(s, NestedIndentPrefix)
}

// IndentDetails indents s using the DetailsIndent prefix.
func IndentDetails(s string) string {
	return indent.String(s, DetailsIndentPrefix)
}
