package routes

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/henok-tesfu/expense-manager/docs"
	"github.com/henok-tesfu/expense-manager/internal/handlers"
	"github.com/henok-tesfu/expense-manager/internal/jwt"
	"github.com/henok-tesfu/expense-manager/internal/middleware"
	"github.com/henok-tesfu/expense-manager/internal/services"
	httpSwagger "github.com/swaggo/http-swagger"
)

// InitRoutes initializes all application routes
func InitRoutes(tokenService *jwt.TokenService) *mux.Router {
	router := mux.NewRouter()

	// Initialize services
	userService := services.NewUserService()

	// Initialize handlers with dependencies
	userHandler := handlers.NewUserHandler(userService, tokenService)

	// Public routes
	router.HandleFunc("/api/register", userHandler.Register).Methods("POST")
	router.HandleFunc("/api/login", userHandler.Login).Methods("POST")
	router.HandleFunc("/api/auth/refresh", userHandler.Refresh).Methods("POST")

	// Protected routes
	protected := router.PathPrefix("/api").Subrouter()
	protected.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value("user_id").(int)
		w.Write([]byte("Hello, User " + strconv.Itoa(userId)))
	}).Methods("GET")

	// Serve Swagger docs
	docs.SwaggerInfo.BasePath = "/" // Adjust the base path if needed

	// Swagger UI route
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Apply CSP Middleware globally
	router.Use(middleware.ContentSecurityPolicyMiddleware("default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; object-src 'none'"))

	// Apply SecureHeaders middleware globally for protected routes
	protected.Use(middleware.AuthMiddleware(tokenService))
	return router
}
