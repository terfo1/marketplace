package main

import (
	"PROJECTTEST/internal/database"
	"PROJECTTEST/internal/middleware"
	"PROJECTTEST/internal/routes"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file")
	}
	uri, exists := os.LookupEnv("URI")
	if !exists {
		fmt.Println("No var called uri")
	}
	secret, exists := os.LookupEnv("SECRET")
	if !exists {
		fmt.Println("No var called secret")
	}
	middleware.JWTSecret = []byte(secret)
	database.ConnectDB(uri)
	r := routes.InitRoutes()
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println("❌ Server error:", err)
	}
}
