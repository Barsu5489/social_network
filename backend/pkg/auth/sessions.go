// auth/session.go
package auth

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

// Session store with a secure key

// Session name and duration
const (
	SessionName   = "social-network-session"
	SessionMaxAge = 86400 * 7 // 7 days
)

func init() {
	// Configure session store
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   SessionMaxAge,
		HttpOnly: true,
		Secure:   false, 
		SameSite: http.SameSiteStrictMode,
	}
}

// CreateSession creates a new session for a user
func CreateSession(w http.ResponseWriter, r *http.Request, userID string) error {
	session, err := store.Get(r, SessionName)
	if err != nil {
		return err
	}

	// Set session values
	session.Values["user_id"] = userID
	session.Values["authenticated"] = true
	session.Values["session_id"] = uuid.New().String()

	// Save session
	return session.Save(r, w)
}

// ClearSession removes a user's session
func ClearSession(w http.ResponseWriter, r *http.Request) error {
	session, err := store.Get(r, SessionName)
	if err != nil {
		return err
	}

	// Clear session values
	session.Values["user_id"] = nil
	session.Values["authenticated"] = false

	// Set session to expire immediately
	session.Options.MaxAge = -1

	// Save changes
	return session.Save(r, w)
}

// GetUserIDFromSession retrieves the user ID from the session
func GetUserIDFromSession(r *http.Request) (string, error) {
	session, err := store.Get(r, SessionName)
	if err != nil {
		return "", err
	}

	// Check if authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		return "", errors.New("not authenticated")
	}

	// Get user ID
	userID, ok := session.Values["user_id"].(string)
	if !ok {
		return "", errors.New("invalid session")
	}

	return userID, nil
}
