package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/henok-tesfu/expense-manager/internal/jwt"
)

// AuthMiddleware validates access tokens and adds the user ID to the request context
func AuthMiddleware(tokenService *jwt.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Retrieve access token from cookies
			accessTokenCookie, err := r.Cookie("access_token")
			if err != nil {
				http.Error(w, "Unauthorized: Missing access token", http.StatusUnauthorized)
				return
			}

			// Validate the access token
			claims, err := tokenService.ValidateAccessToken(accessTokenCookie.Value)
			log.Println(claims)
			if err != nil {
				if err.Error() == "token is expired" {
					http.Error(w, "Unauthorized: Access token expired", http.StatusUnauthorized)
					return
				}
				http.Error(w, "Unauthorized: Invalid access token", http.StatusUnauthorized)
				return
			}

			// Add user ID to the context
			ctx := context.WithValue(r.Context(), "user_id", claims.UserId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
