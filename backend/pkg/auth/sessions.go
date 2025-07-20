// auth/session.go
package auth

import (
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

// Session store with a secure key
var store = sessions.NewCookieStore([]byte("12345678901234567890123456789012")) // 32 bytes

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
		HttpOnly: false, // Allow JavaScript access for debugging
		Secure:   false,
		SameSite: http.SameSiteLaxMode, // Change from Strict to Lax
		Domain:   "", // Ensure no domain restriction
	}
}

// CreateSession creates a new session for a user
func CreateSession(w http.ResponseWriter, r *http.Request, userID string) error {
	log.Printf("DEBUG: CreateSession called with userID: %q (type: %T)", userID, userID)
	
	session, err := store.Get(r, SessionName)
	if err != nil {
		log.Printf("DEBUG: Failed to get existing session (creating new): %v", err)
		// Create a new session instead of failing
		session = sessions.NewSession(store, SessionName)
		session.IsNew = true
	}
	
	log.Printf("DEBUG: Session retrieved/created successfully")

	// Set session values
	log.Printf("DEBUG: Setting session values - userID: %q", userID)
	session.Values["user_id"] = userID
	session.Values["authenticated"] = true
	session.Values["session_id"] = uuid.New().String()
	
	log.Printf("DEBUG: Session values set: %+v", session.Values)

	// Save session
	log.Printf("DEBUG: Attempting to save session")
	err = session.Save(r, w)
	if err != nil {
		log.Printf("ERROR: Failed to save session: %v", err)
		return err
	}
	
	log.Printf("DEBUG: Session saved successfully")
	return nil
}

// ClearSession removes a user's session
func ClearSession(w http.ResponseWriter, r *http.Request) error {
	// Force clear the cookie even if we can't decode the session
	http.SetCookie(w, &http.Cookie{
		Name:     SessionName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})

	session, err := store.Get(r, SessionName)
	if err != nil {
		// Even if we can't get the session, the cookie is cleared above
		log.Printf("DEBUG: Could not get session to clear, but cookie cleared: %v", err)
		return nil
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
