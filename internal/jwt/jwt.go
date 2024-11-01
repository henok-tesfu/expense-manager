package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Secret keys (use environment variables for production)
var jwtSecret = []byte("your-access-token-secret")
var refreshSecret = []byte("your-refresh-token-secret")

// Access Token Expiration (e.g., 15 minutes)
const accessTokenExpiration = 15 * time.Minute

// Refresh Token Expiration (e.g., 7 days)
const refreshTokenExpiration = 7 * 24 * time.Hour

// Claims structure for access tokens
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateAccessToken generates a short-lived access token
func GenerateAccessToken(username string) (string, error) {
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenExpiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// GenerateRefreshToken generates a long-lived refresh token
func GenerateRefreshToken(username string) (string, error) {
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenExpiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshSecret)
}

// ValidateAccessToken validates an access token
func ValidateAccessToken(tokenString string) (*Claims, error) {
	return validateToken(tokenString, jwtSecret)
}

// ValidateRefreshToken validates a refresh token
func ValidateRefreshToken(tokenString string) (*Claims, error) {
	return validateToken(tokenString, refreshSecret)
}

func validateToken(tokenString string, secret []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, err
	}
	return claims, nil
}
