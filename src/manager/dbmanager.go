package manager

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	_ "github.com/lib/pq"
)

type DatabaseMgr interface {
	ExecuteStatement(query string, args ...interface{}) (sql.Result, error)
	ExecuteQuery(query string, args ...interface{}) (*sql.Rows, error)
	ExecuteQueryRow(query string, args ...interface{}) *sql.Row
}

type DatabaseManager struct {
	Connection *sql.DB
}

func (dm *DatabaseManager) ExecuteStatement(query string, args ...interface{}) (sql.Result, error) {
	// Prepare the statement (For Debugging purposes)
	statement, err := dm.Connection.Prepare(query)
	if err != nil {
		log.Printf("Error preparing statement: %v", err)
	}
	defer statement.Close()

	// Execute the statement with the given arguments
	result, err := statement.Exec(args...)
	if err != nil {
		log.Printf("Error executing statement: %v", err)
	}
	return result, err
}

func (dm *DatabaseManager) ExecuteQuery(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := dm.Connection.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (dm *DatabaseManager) ExecuteQueryRow(query string, args ...interface{}) *sql.Row {
	row := dm.Connection.QueryRow(query, args...)
	return row
}

func InitializeDatabaseConnection() *sql.DB {
	// Check if environment variables are set
	var environment = os.Getenv("ENVIRONMENT")

	var (
		db_host     = os.Getenv(fmt.Sprintf("%s_DB_HOST", environment))
		db_port     = os.Getenv(fmt.Sprintf("%s_DB_PORT", environment))
		db_user     = os.Getenv(fmt.Sprintf("%s_DB_USER", environment))
		db_password = os.Getenv(fmt.Sprintf("%s_DB_PASSWORD", environment))
		db_name     = os.Getenv(fmt.Sprintf("%s_DB_NAME", environment))
	)

	if utils.ContainsEmptyString(db_host, db_port, db_user, db_password, db_name) {
		log.Fatalf("Required environment variables are not set")
		os.Exit(1)
	}

	// CONNECT TO DATABASE
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", db_host, db_port, db_user, db_password, db_name)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
		panic(err)
	}

	return db
}
