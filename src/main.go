package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// CREATE ROUTER
	router := createRouter()

	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// GET ENVIRONMENT VARIABLES
	environment := strings.ToUpper(os.Getenv("ENVIRONMENT"))

	db_host := os.Getenv(fmt.Sprintf("%s_DB_HOST", environment))
	db_port := os.Getenv(fmt.Sprintf("%s_DB_PORT", environment))
	db_user := os.Getenv(fmt.Sprintf("%s_DB_USER", environment))
	db_password := os.Getenv(fmt.Sprintf("%s_DB_PASSWORD", environment))
	db_name := os.Getenv(fmt.Sprintf("%s_DB_NAME", environment))

	// Check if environment variables are set
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

	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
		panic(err)
	}

	// CREATE CONTEXT
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// CREATE CHANNEL TO HANDLE OS SIGNALS
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// RUN SERVER
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting or closing listener:: %v", err)
			os.Exit(0)
		}
	}()

	// WAIT FOR OS SIGNAL
	<-quit

	// SHUTDOWN SERVER
	log.Println("Shutting down server...")

	if err := server.Shutdown(context.TODO()); err != nil {
		log.Fatalf("Error closing server:: %v", err)
	}

	log.Println("Server stopped gracefully")
	os.Exit(0)
}
