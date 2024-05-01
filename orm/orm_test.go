package orm

import (
	"errors"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

type User struct {
	Id    int64  `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
}

func TestORM_Find(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockRows := sqlmock.NewRows([]string{"id", "name", "email"}).
		AddRow(1, "John Doe", "john@example.com").
		AddRow(2, "Jane Smith", "jane@example.com")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email FROM users")).WillReturnRows(mockRows)

	o := New(db)
	var users []User
	err = o.Find(&users, "SELECT id, name, email FROM users")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestORM_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users (name, email) VALUES (?, ?)")).
		WithArgs("John Doe", "john@example.com").
		WillReturnResult(sqlmock.NewResult(1, 1))

	o := New(db)
	user := User{Name: "John Doe", Email: "john@example.com"}
	err = o.Create(user)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestORM_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET name = ?, email = ? WHERE id = ?")).
		WithArgs("John Doe", "john@example.com", int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	o := New(db)
	user := User{Id: 1, Name: "John Doe", Email: "john@example.com"}
	err = o.Update(user)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestORM_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM users WHERE id = ?")).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	o := New(db)
	user := User{Id: 1}
	err = o.Delete(user)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestORM_buildInsertQuery(t *testing.T) {
	o := New(nil)
	user := User{Name: "John Doe", Email: "john@example.com"}
	query, args, err := o.buildInsertQuery(&user)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedQuery := "INSERT INTO users (name, email) VALUES (?, ?)"
	if query != expectedQuery {
		t.Errorf("expected query '%s', got '%s'", expectedQuery, query)
	}

	expectedArgs := []interface{}{"John Doe", "john@example.com"}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("expected args %v, got %v", expectedArgs, args)
	}
}

func TestORM_buildUpdateQuery(t *testing.T) {
	o := New(nil)
	user := User{Id: 1, Name: "John Doe", Email: "john@example.com"}
	query, args, err := o.buildUpdateQuery(&user)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedQuery := "UPDATE users SET name = ?, email = ? WHERE id = ?"
	if query != expectedQuery {
		t.Errorf("expected query '%s', got '%s'", expectedQuery, query)
	}

	expectedArgs := []interface{}{"John Doe", "john@example.com", int64(1)}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("expected args %v, got %v", expectedArgs, args)
	}
}

func TestORM_buildDeleteQuery(t *testing.T) {
	o := New(nil)
	user := User{Id: 1}
	query, args, err := o.buildDeleteQuery(&user)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedQuery := "DELETE FROM users WHERE id = ?"
	if query != expectedQuery {
		t.Errorf("expected query '%s', got '%s'", expectedQuery, query)
	}

	expectedArgs := []interface{}{int64(1)}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("expected args %v, got %v", expectedArgs, args)
	}
}

func TestORM_getTableName(t *testing.T) {
	user := User{}
	tableName := getTableName(&user)
	expectedTableName := "users"
	if tableName != expectedTableName {
		t.Errorf("expected table name '%s', got '%s'", expectedTableName, tableName)
	}
}

func TestORM_buildInsertQuery_InvalidEntity(t *testing.T) {
	o := New(nil)
	_, _, err := o.buildInsertQuery(123)
	if err == nil {
		t.Error("expected an error for invalid entity type, got nil")
	}
}

func TestORM_buildUpdateQuery_InvalidEntity(t *testing.T) {
	o := New(nil)
	_, _, err := o.buildUpdateQuery(123)
	if err == nil {
		t.Error("expected an error for invalid entity type, got nil")
	}
}
func TestORM_buildUpdateQuery_NoIDField(t *testing.T) {
	o := New(nil)
	type NoIDStruct struct {
		Name string `db:"name"`
	}
	_, _, err := o.buildUpdateQuery(NoIDStruct{})
	if err == nil {
		t.Error("expected an error for missing ID field, got nil")
	}
}

func TestORM_buildDeleteQuery_InvalidEntity(t *testing.T) {
	o := New(nil)
	_, _, err := o.buildDeleteQuery(123)
	if err == nil {
		t.Error("expected an error for invalid entity type, got nil")
	}
}
func TestORM_buildDeleteQuery_NoIDField(t *testing.T) {
	o := New(nil)
	type NoIDStruct struct {
		Name string `db:"name"`
	}

	_, _, err := o.buildDeleteQuery(NoIDStruct{})
	if err == nil {
		t.Error("expected an error for missing ID field, got nil")
	}
}
func TestORM_Find_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email FROM users")).WillReturnError(errors.New("query error"))

	o := New(db)
	var users []User
	err = o.Find(&users, "SELECT id, name, email FROM users")
	if err == nil {
		t.Error("expected an error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}
func TestORM_Create_BuildQueryError(t *testing.T) {
	o := New(nil)
	err := o.Create(123)
	if err == nil {
		t.Error("expected an error for invalid entity type, got nil")
	}
}
func TestORM_Update_BuildQueryError(t *testing.T) {
	o := New(nil)
	err := o.Update(123)
	if err == nil {
		t.Error("expected an error for invalid entity type, got nil")
	}
}
func TestORM_Delete_BuildQueryError(t *testing.T) {
	o := New(nil)
	err := o.Delete(123)
	if err == nil {

		t.Error("expected an error for invalid entity type, got nil")
	}
}
