package middleware

import (
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request has a valid token (or session, depending on your auth system)
		token := r.Header.Get("Authorization")
		if token == "" || !isValidToken(token) { // Replace isValidToken with your actual token validation logic
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Token is valid; proceed with the request
		next.ServeHTTP(w, r)
	})
}

// Dummy token validation function (replace with actual logic)
func isValidToken(token string) bool {
	// Here you would typically parse and validate the JWT token or check session
	return token == "valid-token" // This is just an example; replace with actual validation
}
