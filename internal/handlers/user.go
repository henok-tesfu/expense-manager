package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/henok-tesfu/expense-manager/internal/services"
	"github.com/henok-tesfu/expense-manager/internal/utils"
)

type ErrorResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

var userService = services.NewUserService()

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {

	var RegisterInput struct {
		Username string `json:"username" validate:"required,min=3"` // Must be at least 3 characters
		Email    string `json:"email" validate:"required,email"`    // Must be a valid email
		Password string `json:"password" validate:"required,min=6"` // Must be at least 6 characters
	}
	// Decode JSON input
	if err := json.NewDecoder(r.Body).Decode(&RegisterInput); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate the input (pass the pointer)
	valid, validationErrors := utils.ValidateStruct(&RegisterInput)
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Message: "Validation failed",
			Errors:  validationErrors,
		})
		return
	}
	newUser, err := userService.RegisterUser(RegisterInput.Username, RegisterInput.Email, RegisterInput.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(newUser)
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email" validate:"required,email"` // Must be a valid email
		Password string `json:"password" validate:"required"`    // Must not be empty
	}

	// Decode the JSON body
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	log.Println(credentials)

	// Validate the input (pass the pointer)
	valid, validationErrors := utils.ValidateStruct(&credentials)
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Message: "Validation failed",
			Errors:  validationErrors,
		})
		return
	}

	// Authenticate the user
	user, err := userService.AuthenticateUser(credentials.Email, credentials.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Respond with the authenticated user
	json.NewEncoder(w).Encode(user)
}
