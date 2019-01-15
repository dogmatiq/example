package render

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
	"unsafe"

	"github.com/dogmatiq/examples/dogmatest/internal/ioutil"
)

func pretty(v interface{}) string {
	var b strings.Builder

	p := printer{}
	p.format(&b, reflect.ValueOf(v), true)

	return b.String()
}

type printer struct {
	visited map[uintptr]struct{}
}

func (p *printer) printf(w io.Writer, f string, v ...interface{}) {
	ioutil.MustWriteString(
		w,
		fmt.Sprintf(f, v...),
	)
}

func (p *printer) format(w io.Writer, v reflect.Value, withType bool) {
	if p.formatScalar(w, v, withType) {
		return
	}

	switch v.Kind() {
	case reflect.Struct:
		p.formatStruct(w, v, withType)
	case reflect.Map:
		p.formatMap(w, v, withType)
	case reflect.Ptr:
		p.formatPtr(w, v, withType)
	case reflect.Interface:
		p.formatInterface(w, v, withType)
	default:
		panic("not supported: " + v.Kind().String())
	}
}

func (p *printer) formatStruct(w io.Writer, v reflect.Value, withType bool) {
	t := v.Type()

	isAnonymous := t.Name() == ""

	if withType {
		if isAnonymous {
			ioutil.MustWriteString(w, "struct")
		} else {
			ioutil.MustWriteString(w, t.String())
		}
	}

	fields := t.NumField()

	if fields == 0 {
		ioutil.MustWriteString(w, "{}")
		return
	}

	ioutil.MustWriteString(w, "{\n")

	width := 0

	for i := 0; i < fields; i++ {
		f := t.Field(i)
		n := len(f.Name)
		if n > width {
			width = n
		}
	}

	iw := ioutil.NewIndenter(w, NestedIndentPrefix)
	for i := 0; i < fields; i++ {
		f := t.Field(i)
		n := len(f.Name)

		ioutil.MustWriteString(iw, f.Name)
		ioutil.MustWriteString(iw, ": ")
		ioutil.MustWriteString(iw, strings.Repeat(" ", width-n))

		fv := v.Field(i)

		p.format(
			iw,
			fv,
			isAnonymous || (f.Type.Kind() == reflect.Interface && !fv.IsNil()),
		)

		ioutil.MustWriteString(iw, "\n")
	}

	ioutil.MustWriteString(w, "}")
}

func (p *printer) formatMap(w io.Writer, v reflect.Value, withType bool) {
	t := v.Type()

	if withType {
		ioutil.MustWriteString(w, t.String())
	}

	size := v.Len()

	if size == 0 {
		ioutil.MustWriteString(w, "{}")
		return
	}

	ioutil.MustWriteString(w, "{\n")

	type key struct {
		value  reflect.Value
		pretty string
	}

	var keys []key
	var b strings.Builder
	width := 0

	for _, k := range v.MapKeys() {
		p.format(&b, k, true) // TODO: only include type if ambiguous
		s := b.String()
		b.Reset()

		// TODO: handle multiline renderings
		if len(s) > width {
			width = len(s)
		}

		keys = append(
			keys,
			key{k, s},
		)
	}

	// sort the keys by their pretty-printed representation
	sort.Slice(
		keys,
		func(i, j int) bool {
			return keys[i].pretty < keys[j].pretty
		},
	)

	iw := ioutil.NewIndenter(w, NestedIndentPrefix)
	for _, k := range keys {
		n := len(k.pretty) // TODO: handle multiline

		ioutil.MustWriteString(iw, k.pretty)
		ioutil.MustWriteString(iw, ": ")
		ioutil.MustWriteString(iw, strings.Repeat(" ", width-n))

		mv := v.MapIndex(k.value)

		p.format(
			iw,
			mv,
			true, // TODO:only render type if ambiguous
		)

		ioutil.MustWriteString(iw, "\n")
	}

	ioutil.MustWriteString(w, "}")
}

func (p *printer) formatPtr(w io.Writer, v reflect.Value, withType bool) {
	t := v.Type()

	if p.isVisited(v) {
		if withType {
			ioutil.MustWriteString(w, "&")
			ioutil.MustWriteString(w, t.Elem().String())
			ioutil.MustWriteString(w, "(*recursive*)")
		} else {
			ioutil.MustWriteString(w, "*recursive*")
		}

		return
	}

	if v.IsNil() {
		if withType {
			ioutil.MustWriteString(w, "&")
			ioutil.MustWriteString(w, t.Elem().String())
			ioutil.MustWriteString(w, "(nil)")
		} else {
			ioutil.MustWriteString(w, "nil")
		}

		return
	}

	ioutil.MustWriteString(w, "&")
	p.format(w, v.Elem(), withType)
}

func (p *printer) formatInterface(w io.Writer, v reflect.Value, withType bool) {
	t := v.Type()

	if v.IsNil() {
		if withType {
			ioutil.MustWriteString(w, t.String())
			ioutil.MustWriteString(w, "(nil)")
		} else {
			ioutil.MustWriteString(w, "nil")
		}

		return
	}

	p.format(w, v.Elem(), withType)
}

func (p *printer) isVisited(v reflect.Value) bool {
	ptr := v.Pointer()

	if _, ok := p.visited[ptr]; ok {
		return true
	}

	if p.visited == nil {
		p.visited = map[uintptr]struct{}{}
	}

	p.visited[ptr] = struct{}{}

	return false
}

func (p *printer) formatScalar(w io.Writer, v reflect.Value, withType bool) bool {
	switch value := v.Interface().(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		if withType {
			p.printf(w, "%T(%v)", value, value)
		} else {
			p.printf(w, "%v", value)
		}

	case uintptr:
		if withType {
			p.printf(w, "%T(%#v)", value, value)
		} else {
			p.printf(w, "%#v", value)
		}

	case complex64, complex128:
		if withType {
			p.printf(w, "%T%#v", value, value)
		} else {
			p.printf(w, "%#v", value)
		}

	case string,
		bool,
		unsafe.Pointer:
		p.printf(w, "%#v", value)

	default:
		return false
	}

	return true
}
