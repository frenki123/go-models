package db

import (
	"reflect"
	"testing"
)

func assertEqual[T comparable](t *testing.T, expected T, actual T) {
	t.Helper()
	if expected == actual {
		return
	}
	t.Errorf("expected (%+v) is not equal to actual (%+v)", expected, actual)
}

func TestReflectStruct(t *testing.T) {
	_, err := getStructType("test")
	assertEqual(t, ErrNotStruct, err)
	type A struct{}
	var a A
	pa := new(A)
	typ := reflect.TypeOf(a)
	res, err := getStructType(a)
	assertEqual(t, typ, res)
	res, err = getStructType(pa)
	assertEqual(t, typ, res)
}
func TestPointerAddres(t *testing.T) {
	var p1 string
	_, err := getPtrAddress(p1)
	assertEqual(t, ErrNotPointer, err)
	var p2 *string
	ta, err := getPtrAddress(p2)
	assertEqual(t, nil, err)
	assertEqual(t, reflect.ValueOf(p2).Pointer(), ta)
	type A struct{ f string }
	var s A
	_, err = getPtrAddress(s.f)
	assertEqual(t, ErrNotPointer, err)
	ta, err = getPtrAddress(&s.f)
	assertEqual(t, nil, err)
	assertEqual(t, reflect.ValueOf(&s.f).Pointer(), ta)
}

func TestFieldMaping(t *testing.T) {
	type testStruct struct {
		f0 string
		f1 int
		f2 string
		f3 *string
		f4 string
		f5 float32
		f6 float64
	}
	var ts testStruct
	testTyp := reflect.TypeOf(ts)
	f1 := testTyp.Field(1)
	f3 := testTyp.Field(3)
	f6 := testTyp.Field(6)
	a1, err := getPtrAddress(&ts.f1)
	a3, err := getPtrAddress(&ts.f3)
	a6, err := getPtrAddress(&ts.f6)
	_, err = mapStructFields(ts)
	assertEqual(t, ErrNotPointer, err)
	var testString string
	_, err = mapStructFields(&testString)
	assertEqual(t, ErrNotStruct, err)
	m, err := mapStructFields(&ts)
	_, ok := m[uintptr(99)]
	assertEqual(t, false, ok)
	r, ok := m[a1]
	assertEqual(t, f1.Name, r.Name)
	r, ok = m[a3]
	assertEqual(t, f3.Name, r.Name)
	r, ok = m[a6]
	assertEqual(t, f6.Name, r.Name)
}

type testStruct struct {
	f1 string
	f2 int
	f3 *string
}

func TestIsPointerToStruct(t *testing.T) {
	var s string
	_, err := pointerToStruct(s)
	assertEqual(t, ErrNotPointer, err)
	i := new(int)
	_, err = pointerToStruct(i)
	assertEqual(t, ErrNotStruct, err)
	type St struct{}
	var st St
	var p *St
	r, err := pointerToStruct(p)
	assertEqual(t, nil, err)
	assertEqual(t, reflect.TypeOf(st), r)
}
