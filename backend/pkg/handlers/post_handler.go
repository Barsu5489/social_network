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

	"github.com/gorilla/mux"
)

func NewPost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		_, err := models.CreatePost(db, ctx, userID, req.Content, req.Privacy, req.GroupID, req.AllowedUserIDs, req.ImageURL)
		if err != nil {
			fmt.Print(err)
			http.Error(w, "Error creating Post", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Post Created successfully"})
	}
}

func FollowingPosts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {

			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		userID := r.Context().Value("user_id").(string)

		posts, err := models.GetFollowingPosts(db, userID)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error getting Post", http.StatusMethodNotAllowed)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Encode posts as JSON
		json.NewEncoder(w).Encode(posts)
	}
}

// Get all posts
func AllPosts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get the authenticated user ID from context
		userID, ok := r.Context().Value("user_id").(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Get all posts
		posts, err := models.GetAllPosts(db, userID)
		if err != nil {
			http.Error(w, "Error getting posts", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(posts)
	}
}

func DeletPost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("user_id").(string)
		vars := mux.Vars(r)
		postID := vars["post_id"]

		err := models.DeletePost(db, postID, userID)

		if err != nil {
			http.Error(w, "Error deleting post: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("message: Post deleted successfully")
	}
}

// GetSinglePost retrieves a single post by ID
func GetSinglePost(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID, ok := r.Context().Value("user_id").(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		vars := mux.Vars(r)
		postID := vars["post_id"]
		if postID == "" {
			http.Error(w, "Post ID is required", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Query to get the post with privacy checks
		query := `
			SELECT p.id, p.user_id, p.content, p.privacy, p.created_at,
				(SELECT COUNT(*) FROM likes 
				WHERE likeable_type = 'post' AND likeable_id = p.id AND deleted_at IS NULL) as likes_count,
				EXISTS(SELECT 1 FROM likes 
						WHERE likeable_type = 'post' AND likeable_id = p.id 
						AND user_id = ? AND deleted_at IS NULL) as user_liked
			FROM posts p
			WHERE p.id = ? AND p.deleted_at IS NULL AND (
				p.privacy = 'public' OR
				p.user_id = ? OR
				(p.privacy = 'almost_private' AND EXISTS(
					SELECT 1 FROM follows 
					WHERE follower_id = ? AND followed_id = p.user_id AND status = 'accepted'
				))
			)
		`

		var post models.Post
		err := db.QueryRowContext(ctx, query, userID, postID, userID, userID).Scan(
			&post.ID, &post.UserID, &post.Content, &post.Privacy,
			&post.CreatedAt, &post.LikesCount, &post.UserLiked,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Post not found or access denied", http.StatusNotFound)
				return
			}
			log.Printf("Error fetching post: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(post)
	}
}
