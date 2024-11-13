package main

import (
	"log"
	"net/http"
	"os"

	"github.com/henok-tesfu/expense-manager/internal/database"
	"github.com/henok-tesfu/expense-manager/internal/routes"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	loadEnvOrPanic()

	// Connect to the database and defer closing the connection
	database.ConnectDatabase()
	defer closeDatabase()

	// Initialize routes
	router := routes.InitRoutes()

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

	requiredVars := []string{"DB_USER", "DB_PASSWORD", "DB_HOST", "DB_PORT", "DB_NAME"}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			log.Fatalf("Environment variable %s is required but not set", v)
		}
	}
}
