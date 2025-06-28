package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"social-nework/pkg/models"

	"github.com/gorilla/mux"
)

// LikePost handles like/unlike requests for posts
func LikePost(db *sql.DB, notificationModel *models.NotificationModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the user ID from context (from auth middleware)
		userID, ok := r.Context().Value("user_id").(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Get post ID from URL params
		vars := mux.Vars(r)
		postID := vars["post_id"]
		if postID == "" {
			http.Error(w, "Post ID is required", http.StatusBadRequest)
			return
		}

		// Create a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Check if the post exists and user has access to it
		var postExists bool
		checkStmt := `
			SELECT EXISTS(
				SELECT 1 FROM posts 
				WHERE id = ? AND deleted_at IS NULL AND (
					privacy = 'public' OR
					user_id = ? OR
					(privacy = 'almost_private' AND EXISTS(
						SELECT 1 FROM follows 
						WHERE follower_id = ? AND followed_id = posts.user_id AND status = 'accepted'
					))
				)
			)
		`
		err := db.QueryRowContext(ctx, checkStmt, postID, userID, userID).Scan(&postExists)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if !postExists {
			http.Error(w, "Post not found or access denied", http.StatusNotFound)
			return
		}

		// Check if we're liking or unliking
		isLike := r.Method == http.MethodPost

		if isLike {
			// Like the post
			like, err := models.CreateLike(db, ctx, notificationModel, userID, "post", postID)
			if err != nil {
				// If the error is because user already liked the post, return a specific status code
				if err.Error() == "user already liked this content" {
					http.Error(w, err.Error(), http.StatusConflict)
					return
				}
				http.Error(w, "Error liking post: "+err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Post liked successfully",
				"like":    like,
			})
		} else {
			// Unlike the post
			err := models.UnlikeContent(db, ctx, userID, "post", postID)
			if err != nil {
				if err.Error() == "like not found or already removed" {
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
				http.Error(w, "Error unliking post: "+err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Post unliked successfully",
			})
		}
	}
}

// GetPostLikes retrieves all likes for a specific post
func GetPostLikes(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context (from auth middleware)
		userID, ok := r.Context().Value("user_id").(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Get post ID from URL params
		vars := mux.Vars(r)
		postID := vars["post_id"]
		if postID == "" {
			http.Error(w, "Post ID is required", http.StatusBadRequest)
			return
		}

		// Create a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Check if the post exists and user has access to it
		var postExists bool
		checkStmt := `
			SELECT EXISTS(
				SELECT 1 FROM posts 
				WHERE id = ? AND deleted_at IS NULL AND (
					privacy = 'public' OR
					user_id = ? OR
					(privacy = 'almost_private' AND EXISTS(
						SELECT 1 FROM follows 
						WHERE follower_id = ? AND followed_id = posts.user_id AND status = 'accepted'
					))
				)
			)
		`
		err := db.QueryRowContext(ctx, checkStmt, postID, userID, userID).Scan(&postExists)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if !postExists {
			http.Error(w, "Post not found or access denied", http.StatusNotFound)
			return
		}

		// Get all likes for the post
		likes, err := models.GetPostLikes(db, ctx, postID)
		if err != nil {
			http.Error(w, "Error getting likes: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Get like count
		count, err := models.GetLikeCount(db, ctx, "post", postID)
		if err != nil {
			http.Error(w, "Error counting likes: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if user has liked the post
		hasLiked, err := models.HasUserLikedPost(db, ctx, userID, postID)
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

// GetUserLikedPosts gets all posts that a user has liked
func GetUserLikedPosts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context (from auth middleware)
		userID, ok := r.Context().Value("user_id").(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Optional: Get target user ID from URL params
		vars := mux.Vars(r)
		targetUserID := vars["user_id"]
		if targetUserID == "" {
			targetUserID = userID // Default to logged-in user
		}

		// Create a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Query to get posts liked by the target user
		stmt := `
			SELECT p.id, p.user_id, p.content, p.privacy, p.created_at,
				(SELECT COUNT(*) FROM likes 
				WHERE likeable_type = 'post' AND likeable_id = p.id AND deleted_at IS NULL) as likes_count,
				EXISTS(SELECT 1 FROM likes 
						WHERE likeable_type = 'post' AND likeable_id = p.id 
						AND user_id = ? AND deleted_at IS NULL) as user_liked
			FROM posts p
			JOIN likes l ON p.id = l.likeable_id
			WHERE l.likeable_type = 'post' AND l.user_id = ? AND l.deleted_at IS NULL
			AND p.deleted_at IS NULL
			AND (
				p.privacy = 'public' OR
				p.user_id = ? OR
				(p.privacy = 'almost_private' AND EXISTS(
					SELECT 1 FROM follows 
					WHERE follower_id = ? AND followed_id = p.user_id AND status = 'accepted'
				))
			)
			ORDER BY l.created_at DESC
		`

		rows, err := db.QueryContext(ctx, stmt, userID, targetUserID, userID, userID)
		if err != nil {
			http.Error(w, "Error getting liked posts: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var posts []models.Post
		for rows.Next() {
			var post models.Post
			err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.Privacy, &post.CreatedAt, &post.LikesCount, &post.UserLiked)
			if err != nil {
				http.Error(w, "Error processing posts: "+err.Error(), http.StatusInternalServerError)
				return
			}
			posts = append(posts, post)
		}

		if posts == nil {
			posts = []models.Post{}
		}

		// Return response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(posts)
	}
}
