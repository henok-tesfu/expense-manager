package routes

import (
	"github.com/gorilla/mux"
	"github.com/henok-tesfu/expense-manager/internal/handlers"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter()

	// Create a subRouter for /api
	apiRouter := router.PathPrefix("/api").Subrouter()

	// Define the /api routes
	apiRouter.HandleFunc("/register", handlers.RegisterUserHandler).Methods("POST")
	apiRouter.HandleFunc("/login", handlers.LoginUserHandler).Methods("POST")

	return router
}
