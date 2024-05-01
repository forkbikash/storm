package orm

import (
	"fmt"
	"reflect"
	"strings"
)

func (o *ORM) Create(entity interface{}) error {
	query, args, err := o.buildInsertQuery(entity)
	if err != nil {
		return err
	}

	_, err = o.db.Exec(query, args...)
	return err
}

func (o *ORM) Update(entity interface{}) error {
	query, args, err := o.buildUpdateQuery(entity)
	if err != nil {
		return err
	}

	_, err = o.db.Exec(query, args...)
	return err
}

func (o *ORM) Delete(entity interface{}) error {
	query, args, err := o.buildDeleteQuery(entity)
	if err != nil {
		return err
	}

	_, err = o.db.Exec(query, args...)
	return err
}

func (o *ORM) buildInsertQuery(entity interface{}) (string, []interface{}, error) {
	t := reflect.TypeOf(entity)
	if t.Kind() != reflect.Struct {
		return "", nil, fmt.Errorf("entity must be a struct")
	}

	var columns []string
	var values []interface{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		columnName := field.Tag.Get("db")
		if columnName != "" {
			columns = append(columns, columnName)
			values = append(values, reflect.ValueOf(entity).Field(i).Interface())
		}
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", getTableName(entity), strings.Join(columns, ", "), strings.Repeat("?, ", len(columns)-1)+"?")
	return query, values, nil
}

func (o *ORM) buildUpdateQuery(entity interface{}) (string, []interface{}, error) {
	t := reflect.TypeOf(entity)
	if t.Kind() != reflect.Struct {
		return "", nil, fmt.Errorf("entity must be a struct")
	}

	var setClause []string
	var values []interface{}

	idField, ok := t.FieldByName("Id")
	if !ok {
		return "", nil, fmt.Errorf("entity must have an 'Id' field")
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		columnName := field.Tag.Get("db")
		if columnName != "" && field.Name != idField.Name {
			setClause = append(setClause, fmt.Sprintf("%s = ?", columnName))
			values = append(values, reflect.ValueOf(entity).Field(i).Interface())
		}
	}

	values = append(values, reflect.ValueOf(entity).FieldByName(idField.Name).Interface())
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?", getTableName(entity), strings.Join(setClause, ", "), idField.Tag.Get("db"))
	return query, values, nil
}

func (o *ORM) buildDeleteQuery(entity interface{}) (string, []interface{}, error) {
	t := reflect.TypeOf(entity)
	if t.Kind() != reflect.Struct {
		return "", nil, fmt.Errorf("entity must be a struct")
	}

	idField, ok := t.FieldByName("Id")
	if !ok {
		return "", nil, fmt.Errorf("entity must have an 'Id' field")
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", getTableName(entity), idField.Tag.Get("db"))
	return query, []interface{}{reflect.ValueOf(entity).FieldByName(idField.Name).Interface()}, nil
}

func getTableName(entity interface{}) string {
	// Implement logic to determine the table name based on the entity's struct name
	return strings.ToLower(reflect.TypeOf(entity).Name()) + "s"
}
