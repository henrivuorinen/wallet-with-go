package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"wallet-with-go/internal/database"
	"wallet-with-go/internal/handlers"
	"wallet-with-go/internal/services"
)

func main() {
	// Initialize database
	db, err := database.InitDatabase()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Failed to close database: %v", err)
		}
	}()

	// Initialize services
	walletService := &services.WalletService{DB: db}

	// Set up Gin router and handlers
	router := gin.Default()
	walletHandler := &handlers.WalletHandler{Service: walletService}
	walletHandler.SetupRoutes(router)

	// Create HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start server in a separate goroutine
	go func() {
		log.Println("Starting HTTP server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Set up signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-quit
	log.Println("Shutting down server...")

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
