package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/db"
	"github.com/Travel-Utilities-WWI21SEB/expense-management-service/src/utils"
	"github.com/joho/godotenv"
)

func main() {
	// CREATE ROUTER
	router := createRouter()

	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// CONNECT TO DATABASE
	db.ConnectToDB()
	defer utils.CloseDbConnection()

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
