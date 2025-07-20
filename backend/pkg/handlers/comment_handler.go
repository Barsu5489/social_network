package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"social-nework/pkg/models"
	"social-nework/pkg/websocket"

	"github.com/google/uuid"
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
		fmt.Println("GetPostComments function being called")
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Try to get user ID, but don't require it for public posts
		userID := ""
		if uid := r.Context().Value("user_id"); uid != nil {
			if uidStr, ok := uid.(string); ok {
				userID = uidStr
			}
		}

		// Get post_id from URL parameters
		vars := mux.Vars(r)
		postID := vars["postId"]
		if postID == "" {
			http.Error(w, "Post ID is required", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		comments, err := models.GetPostComments(db, ctx, postID, userID)
		if err != nil {
			log.Printf("ERROR: Failed to get comments: %v", err)
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

// CreateComment creates a new comment on a post with notifications
func CreateComment(db *sql.DB, notificationModel *models.NotificationModel, hub *websocket.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("CreateComment function being called")
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
		postID := vars["postId"]
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

		// Create notification for post owner if different from commenter
		var postOwnerID string
		err = db.QueryRowContext(ctx, "SELECT user_id FROM posts WHERE id = ?", postID).Scan(&postOwnerID)
		if err == nil && postOwnerID != userID {
			notification := models.Notification{
				ID:          uuid.New().String(),
				UserID:      postOwnerID,
				Type:        "new_comment", // Changed from "comment_on_post" to match DB constraint
				ReferenceID: postID,
				IsRead:      false,
				CreatedAt:   time.Now(),
			}

			if notificationModel != nil {
				_, err = notificationModel.Insert(ctx, notification)
				if err != nil {
					log.Printf("ERROR: Failed to create comment notification: %v", err)
				} else {
					log.Printf("SUCCESS: Comment notification created for post owner: %s", postOwnerID)
					
					// Send real-time notification
					if hub != nil {
						hub.SendNotification(postOwnerID, notification, map[string]interface{}{
							"post_id":      postID,
							"comment_id":   comment.ID,
							"commenter_id": userID,
						})
					}
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Comment created successfully",
			"comment": comment,
		})
	}
}
