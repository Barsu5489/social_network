package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"social-nework/pkg/models"
	"social-nework/pkg/websocket"

	"github.com/gorilla/mux"
)

// LikeComment handles like/unlike requests for comments
func LikeComment(db *sql.DB, notificationModel *models.NotificationModel, hub *websocket.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the user ID from context (from auth middleware)
		userID, ok := r.Context().Value("user_id").(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Get comment ID from URL params
		vars := mux.Vars(r)
		commentID := vars["comment_id"]
		if commentID == "" {
			http.Error(w, "Comment ID is required", http.StatusBadRequest)
			return
		}

		// Create a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Check if the comment exists and user has access to it
		var commentExists bool
		checkStmt := `
			SELECT EXISTS(
				SELECT 1 FROM comments c
				JOIN posts p ON c.post_id = p.id
				WHERE c.id = ? AND c.deleted_at IS NULL AND p.deleted_at IS NULL AND (
					p.privacy = 'public' OR
					p.user_id = ? OR
					c.user_id = ? OR
					(p.privacy = 'almost_private' AND EXISTS(
						SELECT 1 FROM follows 
						WHERE follower_id = ? AND followed_id = p.user_id AND status = 'accepted'
					))
				)
			)
		`
		err := db.QueryRowContext(ctx, checkStmt, commentID, userID, userID, userID).Scan(&commentExists)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if !commentExists {
			http.Error(w, "Comment not found or access denied", http.StatusNotFound)
			return
		}

		// Check if we're liking or unliking
		isLike := r.Method == http.MethodPost

		if isLike {
			// Like the comment
			like, err := models.CreateLike(db, ctx, notificationModel, hub, userID, "comment", commentID)
			if err != nil {
				// If the error is because user already liked the comment, return a specific status code
				if err.Error() == "user already liked this content" {
					http.Error(w, err.Error(), http.StatusConflict)
					return
				}
				http.Error(w, "Error liking comment: "+err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Comment liked successfully",
				"like":    like,
			})
		} else {
			// Unlike the comment
			err := models.UnlikeContent(db, ctx, userID, "comment", commentID)
			if err != nil {
				if err.Error() == "like not found or already removed" {
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
				http.Error(w, "Error unliking comment: "+err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Comment unliked successfully",
			})
		}
	}
}

// GetCommentLikes retrieves all likes for a specific comment
func GetCommentLikes(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context (from auth middleware)
		userID, ok := r.Context().Value("user_id").(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Get comment ID from URL params
		vars := mux.Vars(r)
		commentID := vars["comment_id"]
		if commentID == "" {
			http.Error(w, "Comment ID is required", http.StatusBadRequest)
			return
		}

		// Create a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Check if the comment exists and user has access to it
		var commentExists bool
		checkStmt := `
			SELECT EXISTS(
				SELECT 1 FROM comments c
				JOIN posts p ON c.post_id = p.id
				WHERE c.id = ? AND c.deleted_at IS NULL AND p.deleted_at IS NULL AND (
					p.privacy = 'public' OR
					p.user_id = ? OR
					c.user_id = ? OR
					(p.privacy = 'almost_private' AND EXISTS(
						SELECT 1 FROM follows 
						WHERE follower_id = ? AND followed_id = p.user_id AND status = 'accepted'
					))
				)
			)
		`
		err := db.QueryRowContext(ctx, checkStmt, commentID, userID, userID, userID).Scan(&commentExists)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if !commentExists {
			http.Error(w, "Comment not found or access denied", http.StatusNotFound)
			return
		}

		// Get all likes for the comment
		likes, err := models.GetCommentLikes(db, ctx, commentID)
		if err != nil {
			http.Error(w, "Error getting likes: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Get like count
		count, err := models.GetLikeCount(db, ctx, "comment", commentID)
		if err != nil {
			http.Error(w, "Error counting likes: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if user has liked the comment
		hasLiked, err := models.HasUserLikedComment(db, ctx, userID, commentID)
		if err != nil {
			http.Error(w, "Error checking like status: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Return response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"likes":      likes,
			"count":      count,
			"user_liked": hasLiked,
		})
	}
}
