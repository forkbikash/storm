package db

import (
	"database/sql"
)

type DB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Begin() (*sql.Tx, error)
	Close() error
}

type PostgresDB struct {
	*sql.DB
}

func NewPostgresDB(dsn string) (*PostgresDB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &PostgresDB{db}, nil
}
