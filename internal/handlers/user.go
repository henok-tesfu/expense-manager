package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/henok-tesfu/expense-manager/internal/jwt"
	"github.com/henok-tesfu/expense-manager/internal/services"
	"github.com/henok-tesfu/expense-manager/internal/utils"
)

// ErrorResponse represents the standard error response structure
type ErrorResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

// SuccessResponse represents a standard success response structure
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// UserHandler contains dependencies for user-related operations
type UserHandler struct {
	UserService  *services.UserService
	TokenService *jwt.TokenService
}

// RegisterInput represents the input structure for user registration
type RegisterInput struct {
	Username string `json:"username" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginInput represents the input structure for user login
// @Description Input payload for login
type LoginInput struct {
	Email    string `json:"email" validate:"required,email"` // User email
	Password string `json:"password" validate:"required"`    // User password
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService *services.UserService, tokenService *jwt.TokenService) *UserHandler {
	return &UserHandler{
		UserService:  userService,
		TokenService: tokenService,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Create a new user with a username, email, and password
// @Tags User
// @Accept json
// @Produce json
// @Param RegisterInput body RegisterInput true "Register Input"
// @Success 201 {object} SuccessResponse "User registered successfully"
// @Failure 422 {object} ErrorResponse "Validation or payload errors"
// @Router /api/register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var registerInput RegisterInput

	// Decode and validate the input
	if err := json.NewDecoder(r.Body).Decode(&registerInput); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	if valid, validationErrors := utils.ValidateStruct(&registerInput); !valid {
		respondWithError(w, http.StatusUnprocessableEntity, "Validation failed", validationErrors)
		return
	}

	// Register the user
	newUser, err := h.UserService.RegisterUser(registerInput.Username, registerInput.Email, registerInput.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	respondWithSuccess(w, http.StatusCreated, "User registered successfully", newUser)
}

// Login handles user login
// @Summary Login a user
// @Description Authenticate a user with an email and password
// @Tags User
// @Accept json
// @Produce json
// @Param LoginInput body LoginInput true "Login Input"
// @Success 200 {object} SuccessResponse "Login successful"
// @Failure 400 {object} ErrorResponse "payload errors"
// @Failure 422 {object} ErrorResponse "Validation errors"
// @Failure 401 {object} ErrorResponse "Invalid email or password"
// @Router /api/login [post]
// Login handles user login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials LoginInput

	// Decode JSON payload
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON payload", nil)
		return
	}

	// Validate user input
	if valid, validationErrors := utils.ValidateStruct(&credentials); !valid {
		respondWithError(w, http.StatusUnprocessableEntity, "Validation failed", validationErrors)
		return
	}

	// Authenticate user credentials
	user, err := h.UserService.AuthenticateUser(credentials.Email, credentials.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password", nil)
		return
	}

	// Generate tokens
	if err := h.generateAndSetTokens(w, user.ID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate tokens", nil)
		return
	}

	// Send a success response with user data
	respondWithSuccess(w, http.StatusOK, "Login successful", map[string]interface{}{
		"user": map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

// Refresh handles token refreshing
// @Summary Refresh the access token
// @Description Generate a new access token using a valid refresh token
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} SuccessResponse "Access token refreshed"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Router /api/auth/refresh [post]
func (h *UserHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	// Get the refresh token from cookies
	refreshTokenCookie, err := r.Cookie("refresh_token")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing refresh token", nil)
		return
	}

	// Validate the refresh token
	refreshClaims, err := h.TokenService.ValidateRefreshToken(refreshTokenCookie.Value)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired refresh token", nil)
		return
	}

	// Generate a new access token
	newAccessToken, err := h.TokenService.GenerateAccessToken(refreshClaims.UserId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate new access token", nil)
		return
	}

	// Set the new access token in cookies
	h.setCookie(w, "access_token", newAccessToken, h.TokenService.Config.AccessTokenExpiry, "/")

	respondWithSuccess(w, http.StatusOK, "Access token refreshed", nil)
}

// generateAndSetTokens generates access and refresh tokens and sets them as cookies
func (h *UserHandler) generateAndSetTokens(w http.ResponseWriter, userId int) error {
	accessToken, err := h.TokenService.GenerateAccessToken(userId)
	if err != nil {
		return err
	}

	refreshToken, err := h.TokenService.GenerateRefreshToken(userId)
	if err != nil {
		return err
	}

	// Set cookies
	h.setCookie(w, "access_token", accessToken, h.TokenService.Config.AccessTokenExpiry, "/")
	h.setCookie(w, "refresh_token", refreshToken, h.TokenService.Config.RefreshTokenExpiry, "/api/auth/refresh")

	return nil
}

// setCookie simplifies setting HTTP-only cookies
func (h *UserHandler) setCookie(w http.ResponseWriter, name, value string, expiry time.Duration, path string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     path,
		Expires:  time.Now().Add(expiry),
	})
}

// respondWithError sends a JSON error response
func respondWithError(w http.ResponseWriter, statusCode int, message string, errors map[string]string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Message: message,
		Errors:  errors,
	})
}

// respondWithSuccess sends a JSON success response
func respondWithSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(SuccessResponse{
		Message: message,
		Data:    data,
	})
}
