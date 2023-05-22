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

func ExecuteInsert(query string, args ...interface{}) error {
	_, err := connection.Exec(query, args...)
	return err
}

func ExecuteQuery(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := connection.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
