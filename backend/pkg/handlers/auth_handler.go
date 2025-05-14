package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"social-nework/pkg/models"

	"github.com/google/uuid"
)
type UserModel interface {
    Insert(user models.User) error
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
