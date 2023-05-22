package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	_ "github.com/lib/pq"
)

func ConnectToDB() {
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

	utils.InitializeDbConnection(db)
}
