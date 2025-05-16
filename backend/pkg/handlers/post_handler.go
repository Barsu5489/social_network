package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"social-nework/pkg/models"
)

type PostHandler struct {
	Post *models.PostModel
}

func (h *PostHandler) NewPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID, ok := r.Context().Value("user_id").(string)

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var req models.Post
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	_, err := h.Post.CreatePost(ctx, userID, req.Content, req.Privacy, req.GroupID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post Created successfully"})
}
