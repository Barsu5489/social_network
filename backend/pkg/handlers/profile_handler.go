package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"social-nework/pkg/models"
)

type ProfileHandler struct {
	UserModel *models.User
}

// GET /profile?target_id=...
func GetProfile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID := r.Context().Value("user_id").(string)
		targetID := r.URL.Query().Get("target_id")

		if targetID == "" {
			targetID = userID
		}

		user, posts, followers, following, err := models.GetUser(db, userID, targetID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		response := map[string]interface{}{
			"user":            user,
			"posts":           posts,
			"followers":       followers,
			"following":       following,
			"follower_count":  len(followers),
			"following_count": len(following),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    response,
		})
	}
}

// // UpdateProfile updates user fields based on provided updates
func UpdateProfile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID := r.Context().Value("user_id").(string)
		var updates map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Invalid request body",
			})
			return
		}
		if err := models.UpdateProfile(db, userID, updates); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]bool{
			"success": true,
		})
	}
}
