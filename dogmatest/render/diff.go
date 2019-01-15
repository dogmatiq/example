package render

import (
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// Diff renders a human-readable diff of two strings.
func Diff(a, b string) string {
	d := diffmatchpatch.New()
	var w strings.Builder

	for _, diff := range d.DiffMain(a, b, false) {
		text := diff.Text

		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			w.WriteString(Yellow)
			w.WriteString("{+")
			w.WriteString(text)
			w.WriteString("+}")
			w.WriteString(Reset)
		case diffmatchpatch.DiffDelete:
			w.WriteString(Cyan)
			w.WriteString("[-")
			w.WriteString(text)
			w.WriteString("-]")
			w.WriteString(Reset)
		case diffmatchpatch.DiffEqual:
			w.WriteString(text)
		}
	}

	return w.String()
}
