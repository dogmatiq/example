package render

import (
	"fmt"
	"testing"
)

func TestPretty_ScalarTypes(t *testing.T) {
	var (
		vint8       = int8(-100)
		vint16      = int16(-100)
		vint32      = int32(-100)
		vint64      = int64(-100)
		vuint8      = uint8(100)
		vuint16     = uint16(100)
		vuint32     = uint32(100)
		vuint64     = uint64(100)
		vcomplex64  = complex64(100 + 5i)
		vcomplex128 = complex128(100 + 5i)
		vfloat32    = float32(1.23)
		vfloat64    = float64(1.23)
		vstring     = "foo\nbar"
		vbool       = true
	)

	scalars := map[interface{}]string{
		vint8:       `int8(-100)`,
		vint16:      `int16(-100)`,
		vint32:      `int32(-100)`,
		vint64:      `int64(-100)`,
		vuint8:      `uint8(100)`,
		vuint16:     `uint16(100)`,
		vuint32:     `uint32(100)`,
		vuint64:     `uint64(100)`,
		vcomplex64:  `complex64(100+5i)`,
		vcomplex128: `complex128(100+5i)`,
		vfloat32:    `float32(1.23)`,
		vfloat64:    `float64(1.23)`,
		vstring:     `"foo\nbar"`,
		vbool:       `true`,

		&vint8:       `&int8(-100)`,
		&vint16:      `&int16(-100)`,
		&vint32:      `&int32(-100)`,
		&vint64:      `&int64(-100)`,
		&vuint8:      `&uint8(100)`,
		&vuint16:     `&uint16(100)`,
		&vuint32:     `&uint32(100)`,
		&vuint64:     `&uint64(100)`,
		&vcomplex64:  `&complex64(100+5i)`,
		&vcomplex128: `&complex128(100+5i)`,
		&vfloat32:    `&float32(1.23)`,
		&vfloat64:    `&float64(1.23)`,
		&vstring:     `&"foo\nbar"`,
		&vbool:       `&true`,
	}

	for v, s := range scalars {
		testPretty(
			t,
			fmt.Sprintf("it formats %T correctly", v),
			v,
			s,
		)
	}
}

func TestPretty_Struct(t *testing.T) {
	value := 100

	type static struct {
		Value int
	}

	testPretty(
		t,
		"it renders a struct as expected",
		static{Value: value},
		`render.static{
	Value: 100
}`,
	)

	type empty struct{}

	testPretty(
		t,
		"it renders an empty struct on a single line",
		empty{},
		`render.empty{}`,
	)

	type dynamic struct {
		Value interface{}
	}

	testPretty(
		t,
		"it renders the types of values in interface fields",
		dynamic{Value: value},
		`render.dynamic{
	Value: int(100)
}`,
	)

	testPretty(
		t,
		"it renders the types of pointers in interface fields",
		dynamic{Value: &value},
		`render.dynamic{
	Value: &int(100)
}`,
	)

	testPretty(
		t,
		"it renders nil interface fields",
		dynamic{},
		`render.dynamic{
	Value: nil
}`,
	)

	testPretty(
		t,
		"it renders the types of values in anonymous structs",
		struct{ Value int }{Value: value},
		`struct{
	Value: int(100)
}`,
	)

	type pointer struct {
		Value *int
	}

	testPretty(
		t,
		"it renders pointer fields as expected",
		pointer{Value: &value},
		`render.pointer{
	Value: &100
}`,
	)

	testPretty(
		t,
		"it renders nil pointer fields as expected",
		pointer{},
		`render.pointer{
	Value: nil
}`,
	)

	type nested struct {
		Value1 static
		Value2 interface{}
	}

	testPretty(
		t,
		"it renders nil pointer fields as expected",
		nested{
			Value1: static{
				Value: value,
			},
			Value2: static{
				Value: value * 2,
			},
		},
		`render.nested{
	Value1: {
		Value: 100
	}
	Value2: render.static{
		Value: 200
	}
}`,
	)
}

func TestPretty_Map(t *testing.T) {
	static := map[string]string{
		"foo": "bar",
		"baz": "qux",
	}

	testPretty(
		t,
		"it sorts map keys",
		static,
		`map[string]string{
	"baz": "qux"
	"foo": "bar"
}`,
	)
}

func testPretty(
	t *testing.T,
	n string,
	v interface{},
	x string,
) {
	t.Run(
		n,
		func(t *testing.T) {
			p := pretty(v)

			t.Log("expected:\n\n" + x + "\n")
			t.Log("actual:\n\n" + p + "\n")

			if p != x {
				t.FailNow()
			}
		},
	)
}
