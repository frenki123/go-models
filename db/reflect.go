package db

import (
	"fmt"
	"reflect"
)

type fieldMap map[uintptr]reflect.StructField

var (
	ErrCantAddress = fmt.Errorf("Can't get address of passed variable")
	ErrNotPointer  = fmt.Errorf("Pass object is not a pointer")
	ErrNotStruct   = fmt.Errorf("Object is not a struct")
)

func getStructType(i any) (reflect.Type, error) {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		var u reflect.Type
		return u, ErrNotStruct
	}
	return t, nil
}

func getPtrAddress(i any) (uintptr, error) {
	var errRes uintptr
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Pointer {
		return errRes, ErrNotPointer
	}
	return v.Pointer(), nil
}

func mapStructFields(userStruct any) (fieldMap, error) {
	fieldMap := make(fieldMap)
	structVal := reflect.ValueOf(userStruct)
	if structVal.Kind() != reflect.Pointer {
		return fieldMap, ErrNotPointer
	}
	structTyp := reflect.TypeOf(userStruct).Elem()
	elem := structVal.Elem()
	if elem.Kind() != reflect.Struct {
		return fieldMap, ErrNotStruct
	}
	for i := 0; i < elem.NumField(); i++ {
		pointerField := elem.Field(i)
		if pointerField.CanAddr() {
			pointerAddres := pointerField.Addr().Pointer()
			fieldMap[pointerAddres] = structTyp.Field(i)
		}
	}
	return fieldMap, nil
}

func pointerToStruct(userStruct any) (reflect.Type, error) {
	var et reflect.Type
	t := reflect.TypeOf(userStruct)
	if t.Kind() != reflect.Pointer {
		return et, ErrNotPointer
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return et, ErrNotStruct
	}
	return t, nil
}
