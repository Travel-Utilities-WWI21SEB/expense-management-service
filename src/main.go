package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/managers"
	"github.com/joho/godotenv"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	// LOAD ENVIRONMENT VARIABLES
	log.Println("Loading environment variables...")
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file:: %v", err)
	}
	log.Println("Environment variables loaded successfully")

	// INITIALIZE DATABASE CONNECTION
	log.Println("Initializing database connection...")
	dbConnection := managers.InitializeDatabaseConnection()
	defer dbConnection.Close()

	// CREATE ROUTER
	log.Println("Creating router...")
	router := createRouter(dbConnection)
	log.Println("Router created successfully")

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
		log.Println("Starting server on port 8080...")
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
