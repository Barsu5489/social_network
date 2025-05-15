package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"social-nework/pkg/models"
"social-nework/pkg/auth"
	"github.com/google/uuid"
)

type UserModel interface {
	Insert(user models.User) error
	Authenticate(email, password string) (*models.User, error)
}

type AuthHandler struct {
	UserModel UserModel
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	user.ID = uuid.New().String() // Generate unique ID here

	if err := a.UserModel.Insert(user); err != nil {
		log.Println("Insert error:", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	// Set content type and status code for success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Return success response
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "User registered successfully",
		"data": map[string]string{
			"id": user.ID,
		},
	})
}

// handlers/auth.go

// LoginRequest defines the structure for login request data
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login handles user authentication
func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Authenticate user
	user, err := a.UserModel.Authenticate(req.Email, req.Password)
	if err != nil {
		log.Println("Authentication error:", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create user session
	if err := auth.CreateSession(w, r, user.ID); err != nil {
		log.Println("Session creation error:", err)
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Login successful",
		"data": map[string]interface{}{
			"id":         user.ID,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
		},
	})
}
// Logout terminates a user's session
func (a *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    // Clear the session
    if err := auth.ClearSession(w, r); err != nil {
        log.Println("Logout error:", err)
        http.Error(w, "Failed to logout", http.StatusInternalServerError)
        return
    }
    
    // Return success response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "message": "Logout successful",
    })
}