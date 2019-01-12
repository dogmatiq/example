package ioutil_test

import (
	"strings"
	"testing"

	. "github.com/dogmatiq/examples/dogmatest/internal/ioutil"
)

func TestIndenter(t *testing.T) {
	b := &strings.Builder{}
	w := NewIndenter(b, "")

	n := MustWriteString(w, "fo")
	n += MustWriteString(w, "o\nb")
	n += MustWriteString(w, "ar\n")
	n += MustWriteString(w, "baz")

	expected := "    foo\n    bar\n    baz"

	result := b.String()
	if result != expected {
		t.Fatalf(
			"unexpected output: %s, expected %s",
			result,
			expected,
		)
	}

	if n != len(expected) {
		t.Fatalf(
			"unexpected byte count: %d, expected %d",
			n,
			len(expected),
		)
	}
}
