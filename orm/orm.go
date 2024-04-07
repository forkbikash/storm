package orm

import (
	"database/sql"
	"reflect"
	"storm/db"
	"strings"
)

type ORM struct {
	db db.DB
}

func New(db db.DB) *ORM {
	return &ORM{db: db}
}

func (o *ORM) Find(dest interface{}, query string, args ...interface{}) error {
	rows, err := o.db.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Hydrate the destination slice
	destValue := reflect.ValueOf(dest).Elem()
	for rows.Next() {
		elem := reflect.New(destValue.Type().Elem()).Interface()
		if err := o.hydrate(elem, rows); err != nil {
			return err
		}
		destValue.Set(reflect.Append(destValue, reflect.ValueOf(elem).Elem()))
	}

	return rows.Err()
}

func (o *ORM) hydrate(dest interface{}, rows *sql.Rows) error {
	// Map database columns to struct fields
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	values := make([]interface{}, len(columns))
	destValue := reflect.ValueOf(dest).Elem()
	for i := range values {
		field := destValue.FieldByName(strings.Title(columns[i]))
		values[i] = field.Addr().Interface()
	}

	return rows.Scan(values...)
}
