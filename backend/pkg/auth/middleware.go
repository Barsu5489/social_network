package auth

import (
	"context"
	"encoding/json"
	"net/http"
	
)

// RequireAuth middleware checks if a user is authenticated
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from session
		userID, err := GetUserIDFromSession(r)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Authentication required",
			})
			return
		}

		// Add user ID to request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_id", userID)

		// Call the next handler with the updated context
		next(w, r.WithContext(ctx))
	}
}
