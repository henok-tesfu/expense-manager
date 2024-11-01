package main

import (
	"log"
	"net/http"

	"github.com/henok-tesfu/expense-manager/internal/routes"
)

func main() {
	// Initialize routes
	router := routes.InitRoutes()

	// Start the server
	log.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
