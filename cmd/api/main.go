package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/henok-tesfu/expense-manager/internal/database"
	"github.com/henok-tesfu/expense-manager/internal/jwt"
	"github.com/henok-tesfu/expense-manager/internal/routes"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	loadEnvOrPanic()

	// Connect to the database and defer closing the connection
	database.ConnectDatabase()
	defer closeDatabase()

	// Load environment variables or set default configurations
	jwtConfig := jwt.Config{
		AccessSecret:       []byte(os.Getenv("ACCESS_SECRET")),
		RefreshSecret:      []byte(os.Getenv("REFRESH_SECRET")),
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 7 * 24 * time.Hour,
	}

	// Initialize the TokenService
	tokenService := jwt.NewTokenService(jwtConfig)

	// Initialize routes
	router := routes.InitRoutes(tokenService)

	// Set the server port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Default to port 8000 if no environment variable is set
	}

	// Start the server
	log.Println("Server starting on port", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func closeDatabase() {
	if database.DB != nil {
		log.Println("Closing database connection...")
		database.DB.Close()
	} else {
		log.Println("Database connection was not initialized, nothing to close.")
	}
}

func loadEnvOrPanic() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	requiredVars := []string{"DB_USER", "DB_PASSWORD", "DB_HOST", "DB_PORT", "DB_NAME", "ACCESS_SECRET", "REFRESH_SECRET"}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			log.Fatalf("Environment variable %s is required but not set", v)
		}
	}
}
