package db

import (
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrFieldNotAPointer = fmt.Errorf("Field needs to be a pointer")
	ErrFieldNotInStruct = fmt.Errorf("Field pointer needs to be defined in the struct")
)

type (
	Field interface {
		Name() string
		//Type() reflect.Type
		ToSQL() string
		//GetValidators() []FieldValidator
	}

	field struct {
		ptrAddress uintptr
		name       string
		goType     reflect.Type
		dbType     string
		primaryKey bool
		unique     bool
		blank      bool
		null       bool
		defaultVal fmt.Stringer
		validators []FieldValidator
	}

	FieldValidator func() error
)

func newField(fieldPtr any, dbType string) *field {
	ptrAddress, err := getPtrAddress(fieldPtr)
	if err != nil {
		panic(err)
	}
	return &field{
		ptrAddress: ptrAddress,
		dbType:     dbType,
	}
}

func (f field) Name() string {
	return f.name
}

func (f *field) nameFromField(fieldMap fieldMap) (*field, error) {
	emptyField := new(field)
	sf, ok := fieldMap[f.ptrAddress]
	if !ok {
		return emptyField, ErrFieldNotInStruct
	}
	f.name = strings.ToLower(sf.Name)
	return f, nil
}

func (f field) ToSQL() string {
	var constraints []string

	if !f.null {
		constraints = append(constraints, "NOT NULL")
	}

	if f.unique {
		constraints = append(constraints, "UNIQUE")
	}

	if f.hasDefault() {
		constraints = append(constraints, f.defaultClause())
	}

	return fmt.Sprintf(
		"'%s' %s%s",
		f.name,
		f.dbType,
		formatConstraints(constraints),
	)
}

func (f field) hasDefault() bool {
	return f.defaultVal != nil && f.defaultVal.String() != ""
}

func (f field) defaultClause() string {
	return fmt.Sprintf("DEFAULT '%s'", f.defaultVal.String())
}

func formatConstraints(constraints []string) string {
	if len(constraints) == 0 {
		return ""
	}
	return " " + strings.Join(constraints, " ")
}

func (f *field) Unique() *field {
	f.unique = true
	return f
}

func (f *field) Nullable() *field {
	f.null = true
	return f
}

func (f *field) Default(obj fmt.Stringer) *field {
	f.defaultVal = obj
	return f
}
