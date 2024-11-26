package jwt

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Config defines JWT configuration settings
type Config struct {
	AccessSecret       []byte
	RefreshSecret      []byte
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

// DefaultConfig provides default JWT settings
var DefaultConfig = Config{
	AccessSecret:       []byte("access_secret_12345"),
	RefreshSecret:      []byte("refresh_secret_67890"),
	AccessTokenExpiry:  10 * time.Second,
	RefreshTokenExpiry: 7 * 24 * time.Hour,
}

// Claims defines custom JWT claims
type Claims struct {
	UserId int `json:"user_id"`
	jwt.RegisteredClaims
}

// TokenService handles JWT operations
type TokenService struct {
	Config Config
}

// NewTokenService creates a new instance of TokenService
func NewTokenService(config Config) *TokenService {
	return &TokenService{Config: config}
}

// GenerateAccessToken generates a short-lived access token
func (ts *TokenService) GenerateAccessToken(userId int) (string, error) {
	claims := &Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ts.Config.AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	log.Println("Access Secret for Signing:", string(ts.Config.AccessSecret)) // Debug
	return token.SignedString(ts.Config.AccessSecret)
}

// GenerateRefreshToken generates a long-lived refresh token
func (ts *TokenService) GenerateRefreshToken(userId int) (string, error) {
	claims := &Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ts.Config.RefreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(ts.Config.RefreshSecret)
}

// ValidateAccessToken validates an access token
func (ts *TokenService) ValidateAccessToken(tokenString string) (*Claims, error) {
	log.Println("Access Secret for Validation:", string(ts.Config.AccessSecret)) // Debug
	return ts.validateToken(tokenString, ts.Config.AccessSecret)
}

// ValidateRefreshToken validates a refresh token
func (ts *TokenService) ValidateRefreshToken(tokenString string) (*Claims, error) {
	return ts.validateToken(tokenString, ts.Config.RefreshSecret)
}

// validateToken validates and parses a JWT token
func (ts *TokenService) validateToken(tokenString string, secret []byte) (*Claims, error) {
	log.Println("Validating token:", tokenString) // Debug
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("invalid token")
		}
		return secret, nil
	})
	if err != nil {
		log.Println("Token validation error:", err) // Debug
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		log.Println("Invalid claims or token not valid") // Debug
		return nil, errors.New("invalid token")
	}
	log.Println("Validated Claims:", claims) // Debug
	return claims, nil
}
