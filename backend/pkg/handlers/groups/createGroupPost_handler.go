package groups

import (
	"encoding/json"
	"net/http"
	"social-nework/pkg/models"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Create group post
func (gh *GroupHandler) CreateGroupPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["groupId"]

	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Check if user is member of the group
	memberQuery := `SELECT id FROM group_members WHERE group_id = ? AND user_id = ? AND deleted_at IS NULL`
	var memberExists string
	err := gh.db.QueryRow(memberQuery, groupID, post.UserID).Scan(&memberExists)
	if err != nil {
		http.Error(w, "You must be a member to post", http.StatusForbidden)
		return
	}

	post.ID = uuid.New().String()
	post.GroupID = &groupID
	post.Privacy = "private" // Group posts are private to group members
	post.CreatedAt = time.Now().Unix()
	post.UpdatedAt = time.Now().Unix()

	query := `INSERT INTO posts (id, user_id, group_id, content, privacy, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err = gh.db.Exec(query, post.ID, post.UserID, post.GroupID, post.Content,
		post.Privacy, post.CreatedAt, post.UpdatedAt)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	// Add to group_posts table
	groupPostQuery := `INSERT INTO group_posts (id, group_id, post_id, created_at) VALUES (?, ?, ?, ?)`
	groupPostID := uuid.New().String()
	_, err = gh.db.Exec(groupPostQuery, groupPostID, groupID, post.ID, time.Now().Unix())
	if err != nil {
		http.Error(w, "Failed to link post to group in group_posts table", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}
