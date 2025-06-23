package groups

import (
	"encoding/json"
	"net/http"
	"social-nework/pkg/models"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Create comment on group post
func (gh *GroupHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["postId"]

	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Check if user is member of the group that owns the post
	memberQuery := `SELECT gm.id FROM group_members gm 
					INNER JOIN posts p ON gm.group_id = p.group_id 
					WHERE p.id = ? AND gm.user_id = ? AND gm.deleted_at IS NULL`
	var memberExists string
	err := gh.db.QueryRow(memberQuery, postID, comment.UserID).Scan(&memberExists)
	if err != nil {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	comment.ID = uuid.New().String()
	comment.PostID = postID
	comment.CreatedAt = time.Now().Unix()
	comment.UpdatedAt = time.Now().Unix()

	query := `INSERT INTO comments (id, post_id, user_id, content, image_url, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err = gh.db.Exec(query, comment.ID, comment.PostID, comment.UserID, comment.Content,
		comment.ImageURL, comment.CreatedAt, comment.UpdatedAt)
	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}
