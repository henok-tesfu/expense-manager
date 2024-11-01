package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/henok-tesfu/expense-manager/internal/services"
)

var userService = services.NewUserService()

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var user struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&user)

	newUser, err := userService.RegisterUser(user.Username, user.Email, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(newUser)
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&credentials)

	user, err := userService.AuthenticateUser(credentials.Email, credentials.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(user)
}
