package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"social-nework/pkg/models"

	"github.com/gorilla/mux"
)

// NewComment creates a new comment on a post
func NewComment(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Is this function being called?")
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        userID, ok := r.Context().Value("user_id").(string)
        if !ok {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Get post_id from URL parameters
        vars := mux.Vars(r)
        postID := vars["post_id"]
        if postID == "" {
            http.Error(w, "Post ID is required", http.StatusBadRequest)
            return
        }

        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        var req struct {
            Content  string  `json:"content"`
            ImageURL *string `json:"image_url,omitempty"`
        }

        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        // Validate content
        if req.Content == "" {
            http.Error(w, "Content is required", http.StatusBadRequest)
            return
        }

        comment, err := models.CreateComment(db, ctx, postID, userID, req.Content, req.ImageURL)
        if err != nil {
            http.Error(w, "Error creating comment", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "message": "Comment created successfully",
            "comment": comment,
        })
    }
}

// GetPostComments retrieves all comments for a specific post
func GetPostComments(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println()
        if r.Method != http.MethodGet {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        userID, ok := r.Context().Value("user_id").(string)
        if !ok {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        // Get post_id from URL parameters
        vars := mux.Vars(r)
        postID := vars["post_id"]
        if postID == "" {
            http.Error(w, "Post ID is required", http.StatusBadRequest)
            return
        }

        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        comments, err := models.GetPostComments(db, ctx, postID, userID)
        if err != nil {
            http.Error(w, "Error getting comments", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "comments": comments,
            "count":    len(comments),
        })
    }
}