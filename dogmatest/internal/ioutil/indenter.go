package ioutil

import (
	"bytes"
	"io"
	"strings"
)

type indenter struct {
	w        io.Writer
	prefix   []byte
	indented bool
}

// NewIndenter returns a writer that write to w, indenting each line with the
// prefix p.
//
// If p is empty, it defaults to four spaces.
func NewIndenter(w io.Writer, p string) io.Writer {
	if p == "" {
		p = "    "
	}

	return &indenter{
		w:      w,
		prefix: []byte(p),
	}
}

// Indent returns s, indented using the prefix p.
//
// If p is empty, it defaults to four spaces.
func Indent(s, p string) string {
	var b strings.Builder
	w := NewIndenter(&b, p)
	MustWriteString(w, s)
	return b.String()
}

func (w *indenter) Write(buf []byte) (n int, err error) {
	defer Recover(&err)

	// keep writing so long as there's something in the buffer
	for len(buf) > 0 {
		// indent if we're ready to do so
		if !w.indented {
			n += MustWrite(w.w, w.prefix)
			w.indented = true
		}

		// find the next line break character
		i := bytes.IndexByte(buf, '\n')

		// if there are no more line break characters, simply write the remainder of
		// the buffer and we're done
		if i == -1 {
			n += MustWrite(w.w, buf)
			break
		}

		// otherwise, write the remainder of this line, including the line break
		// character, and trim the beginning of the buffer
		n += MustWrite(w.w, buf[:i+1])
		buf = buf[i+1:]

		// we're ready for another indent if/when there is more content
		w.indented = false
	}

	return
}
