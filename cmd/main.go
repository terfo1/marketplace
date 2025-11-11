package main

import (
	"PROJECTTEST/internal/database"
	"PROJECTTEST/internal/middleware"
	"PROJECTTEST/internal/routes"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file")
	}

	uri, exists := os.LookupEnv("URI")
	if !exists {
		fmt.Println("No URI provided in .env")
	}

	secret, exists := os.LookupEnv("SECRET")
	if !exists {
		fmt.Println("No SECRET provided in .env")
	}

	middleware.JWTSecret = []byte(secret)
	database.ConnectDB(uri)

	// Setup HTTP server
	r := routes.InitRoutes()
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	<-quit
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server shutdown error: %v\n", err)
	}

	database.CloseDB() // Disconnect MongoDB
	fmt.Println("Server gracefully stopped")
}
