package db

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"
)

type (
	Query[T any] struct {
		schema     schema[T]
		err        error
		conditions []Condition
		orQuery    bool
	}
	Condition interface {
		ToSQL() string
		Error() error
		Name(obj any) (string, error)
		SetName(string)
	}

	stndValue interface {
		string | int | float32 | float64 | time.Time
	}

	condition[T stndValue] struct {
		structField   T
		pointerAddres uintptr
		collumnName   string
		value         string
		operator      sqlOperator
		isNot         bool
		err           error
	}

	stringDataType interface {
		string
	}

	numericDataType interface {
		int | int32 | float32 | float64
	}
	sqlOperator string

	ErrTable struct {
		desc string
		typ  any
	}
)

var (
	ErrCantGetName = fmt.Errorf("Can't name for the field")
)

const (
	opEQ      sqlOperator = "="
	opGT      sqlOperator = ">"
	opLT      sqlOperator = "<"
	opGE      sqlOperator = ">="
	opLE      sqlOperator = "<="
	opNE      sqlOperator = "<>"
	opBETWEEN sqlOperator = "BETWEEN"
	opLIKE    sqlOperator = "LIKE"
	opIN      sqlOperator = "IN"
)

func (e ErrTable) Error() string {
	name := reflect.TypeOf(e.typ).Name()
	return fmt.Sprintf("Table error '%s' on user type '%s'\n", e.desc, name)
}

func newQuery[T any](schema schema[T]) *Query[T] {
	q := Query[T]{schema: schema}
	return &q
}

func getQuery[T any](userStruct T) (*Query[T], error) {
	var emptyQuery *Query[T]
	table, exists := tablesRegistry[getTableName(reflect.TypeOf(userStruct))]
	if !exists {
		return emptyQuery, ErrTable{desc: "Struct not register", typ: userStruct}
	}
	schema, ok := table.(schema[T])
	if !ok {
		return emptyQuery, ErrTable{desc: "Struct has wrong type", typ: userStruct}
	}
	q := newQuery(schema)
	return q, nil
}

func Get[T any](obj T, id int) (T, error) {
	q, err := getQuery(obj)
	if err != nil {
		var er T
		return er, err
	}
	res := q.schema.data
	stmt := fmt.Sprintf("SELECT * FROM %s WHERE ID=$1", q.schema.TableName())
	err = globalDb.Get(res, stmt, id)
	if err != nil {
		var et T
		return et, err
	}
	return res, nil
}

func Save[T any](obj T) error {
	q, err := getQuery(obj)
	if err != nil {
		return err
	}
	cols := q.schema.fieldNames("")
	values := q.schema.fieldNames(":")
	insertStmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		q.schema.TableName(),
		strings.Join(cols, ","),
		strings.Join(values, ","))
	tx := globalDb.MustBegin()
	defer tx.Rollback()
	if _, err := tx.NamedExec(insertStmt, obj); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (q Query[T]) condToSQL() []string {
	var res []string
	for _, s := range q.conditions {
		res = append(res, s.ToSQL())
	}
	return res
}

func (q *Query[T]) All() ([]T, error) {
	if q.err != nil {
		var errRes []T
		return errRes, q.err
	}
	oper := "AND"
	if q.orQuery {
		oper = "OR"
	}
	cond := strings.Join(q.condToSQL(), " "+oper)
	queryStmt := fmt.Sprintf("SELECT * FROM %s WHERE (%s)", q.schema.TableName(), cond)
	fmt.Println(queryStmt)
	var res []T
	err := globalDb.Select(&res, queryStmt)
	return res, err
}

func Filter[T any](structPtr T, conditions ...Condition) *Query[T] {
	q, err := getQuery(structPtr)
	if err != nil {
		errQ := &Query[T]{err: err}
		return errQ
	}
	fieldMap, err := mapStructFields(structPtr)
	if err != nil {
		q.err = err
		return q
	}
	for _, c := range conditions {
		var errorQuery *Query[T]
		if c.Error() != nil {
			errorQuery.err = c.Error()
			return errorQuery
		}
		colName, err := c.Name(fieldMap)
		if err == nil {
			c.SetName(colName)
			q.conditions = append(q.conditions, c)
		}
	}
	return q
}

func Where[T stndValue](t *T) *condition[T] {
	c := new(condition[T])
	ptrAdd, err := getPtrAddress(t)
	if err != nil {
		c.err = err
		return c
	}
	c.pointerAddres = ptrAdd
	return c
}

func (c condition[T]) ToSQL() string {
	not := ""
	if c.isNot {
		not = "NOT"
	}
	return fmt.Sprintf("%s %s %s %v", not, c.collumnName, c.operator, c.value)
}

func (c condition[T]) Error() error {
	return c.err
}

func (c *condition[T]) SetName(name string) {
	c.collumnName = name
}

func (c condition[T]) Name(obj any) (string, error) {
	fm, ok := obj.(fieldMap)
	if !ok {
		return "", ErrCantGetName
	}
	sf, ex := fm[c.pointerAddres]
	if !ex {
		return "", ErrCantGetName
	}
	name := strings.ToLower(sf.Name)
	return name, nil
}

type num interface {
	int | float32 | float64
}

func (c *condition[numericDataType]) Gt(value any) *condition[numericDataType] {
	c.operator = opGT
	v, ok := value.(int)
	c.value = fmt.Sprintf("%d", value)
	return c
}

func (c *condition[string]) Like(value string) *condition[string] {
	c.operator = opLIKE
	c.value = fmt.Sprintf("'%s'", value)
	return c
}
