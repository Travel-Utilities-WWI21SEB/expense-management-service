package utils

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var connection *sql.DB

func InitializeDbConnection(db *sql.DB) {
	connection = db
}

func CloseDbConnection() {
	connection.Close()
}

func ExecuteStatement(query string, args ...interface{}) (sql.Result, error) {
	result, err := connection.Exec(query, args...)
	return result, err
}

func ExecuteQuery(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := connection.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func ExecuteQueryRow(query string, args ...interface{}) *sql.Row {
	row := connection.QueryRow(query, args...)
	return row
}
