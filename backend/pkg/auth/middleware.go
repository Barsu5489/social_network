package auth

import (
	"context"
	"log"
	"net/http"
)

// RequireAuth middleware checks if a user is authenticated
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("DEBUG: Auth middleware - Path: %s, Method: %s", r.URL.Path, r.Method)

		cookie, err := r.Cookie("social-network-session")
		if err != nil {
			log.Printf("DEBUG: No session cookie found for path: %s", r.URL.Path)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Printf("DEBUG: Session cookie found: %s", cookie.Value[:10]+"...")

		// Validate session and get user ID
		userID, err := GetUserIDFromSession(r)
		if err != nil {
			log.Printf("DEBUG: Invalid session for path: %s, error: %v", r.URL.Path, err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Printf("DEBUG: Valid session for user: %s, path: %s", userID, r.URL.Path)

		// Add user ID to request context
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
