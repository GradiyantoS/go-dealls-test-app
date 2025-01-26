package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/GradiyantoS/go-dealls-test-app/utils"
)

// Context key to store the user ID in the request context
type contextKey string

const UserContextKey contextKey = "user_id"

// AuthMiddleware validates the JWT and adds the user ID to the request context
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		userID, err := utils.ValidateJWT(token)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add the user ID to the request context
		r = r.WithContext(context.WithValue(r.Context(), UserContextKey, userID))
		next.ServeHTTP(w, r)
	})
}

// GetUserIDFromContext retrieves the user ID from the request context
func GetUserIDFromContext(r *http.Request) (int, bool) {
	userID, ok := r.Context().Value(UserContextKey).(int)
	return userID, ok
}
