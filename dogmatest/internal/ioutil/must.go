package ioutil

import "io"

// MustWrite writes buf to w.
//
// It panics if w.Write() returns an error.
func MustWrite(w io.Writer, buf []byte) int {
	n, err := w.Write(buf)
	if err != nil {
		panic(errorWrapper{err})
	}

	return n
}

// MustWriteString writes s to w.
//
// It panics if w.Write() returns an error.
func MustWriteString(w io.Writer, s string) int {
	n, err := io.WriteString(w, s)

	if err != nil {
		panic(errorWrapper{err})
	}

	return n
}

// Recover recovers from a panic caused by one of the MustXXX() functions.
//
// The causal error is assigned to *err.
func Recover(err *error) {
	switch v := recover().(type) {
	case errorWrapper:
		*err = v.Err
	case nil:
		return
	default:
		panic(v)
	}
}

type errorWrapper struct {
	Err error
}
