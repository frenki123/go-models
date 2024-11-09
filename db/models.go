package db

import (
	"fmt"
	"reflect"
	"strings"
)

type (
	schema[T any] struct {
		data   T
		typ    reflect.Type
		fields []Field
	}
)

var (
	ErrWrongSchemaType = fmt.Errorf("Schema can only be created from a pointer to a struct!")
	ErrFieldNotInit    = fmt.Errorf("Field is not initialized in the schema")
)

func DefSchema[T any](userStruct T, fields ...Field) (schema[T], error) {
	emptySchema := schema[T]{}
	fieldMap, err := mapStructFields(userStruct)
	if err != nil {
		return emptySchema, err
	}
	structTyp, err := pointerToStruct(userStruct)
	if err != nil {
		return emptySchema, err
	}
	schema := schema[T]{
		data: userStruct,
		typ:  structTyp,
	}
	for _, f := range fields {
		// check if passed filed is struct filed or user defined
		sf, ok := f.(*field)
		if ok {
			if nf, err := sf.nameFromField(fieldMap); err == nil {
				schema.fields = append(schema.fields, nf)
			} else {
				return emptySchema, err
			}
		} else {
			schema.fields = append(schema.fields, f)
		}
	}
	return schema, nil
}

func MustDefSchema[T any](userStruct T, fields ...Field) schema[T] {
	schema, err := DefSchema(userStruct, fields...)
	if err != nil {
		panic(err)
	}
	return schema
}

func (schema schema[T]) ToSQL() string {
	var cols []string
	for _, f := range schema.fields {
		cols = append(cols, f.ToSQL())
	}
	s := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ('id' INTEGER PRIMARY KEY,%s);",
		schema.TableName(),
		strings.Join(cols, ","))
	return s
}

func (schema schema[T]) TableName() string {
	t := schema.typ
	return getTableName(t)
}

func (schema schema[T]) fieldNames(prefix string) []string {
	var res []string
	for _, f := range schema.fields {
		res = append(res, prefix+f.Name())
	}
	return res
}

func getTableName(t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return strings.ToLower(t.Name())
}
