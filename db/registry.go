package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

const (
	DBTYPE     = "DATABASE"
	CONNSTRING = "DBCONNSTRING"
)

type (
	Table interface {
		ToSQL() string
		TableName() string
	}

	ErrDb struct {
		description string
		sqlxErr     error
	}
	Schema interface {
		Schema() Table
	}
)

var (
	tablesRegistry = make(map[string]Table)
	//error to load env
	ErrEnvLoad = fmt.Errorf("Error loading .env file. Check if the .env file exists!")

	globalDb *sqlx.DB
)

func (e ErrDb) Error() string {
	return fmt.Sprintf("Database error: %s\n Description: %v", e.description, e.sqlxErr)
}

func RegisterModels(schemas ...Schema) error {
	for _, s := range schemas {
		t := s.Schema()
		tablesRegistry[t.TableName()] = t
	}
	return migrate()
}

func MustRegister(schemas ...Schema) {
	err := RegisterModels(schemas...)
	if err != nil {
		panic(err)
	}
}

func migrate() error {
	if err := godotenv.Load(); err != nil {
		return ErrEnvLoad
	}
	db, err := sqlx.Connect(os.Getenv(DBTYPE), os.Getenv(CONNSTRING))
	if err != nil {
		return ErrDb{description: "Failed to connect", sqlxErr: err}
	}
	sqlSchema := ""
	for _, t := range tablesRegistry {
		sqlSchema += t.ToSQL()
	}
	if _, err = db.Exec(sqlSchema); err != nil {
		return ErrDb{description: "Failed to execute migration", sqlxErr: err}
	}
	globalDb = db
	return nil
}
