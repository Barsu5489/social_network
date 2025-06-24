package groups

import (
	"encoding/json"
	"net/http"
	"social-nework/pkg/models"

	"github.com/gorilla/mux"
)

// Get group posts
func (gh *GroupHandler) GetGroupPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["groupId"]
	userID := r.URL.Query().Get("user_id")

	// Check if user is member of the group
	memberQuery := `SELECT id FROM group_members WHERE group_id = ? AND user_id = ? AND deleted_at IS NULL`
	var memberExists string
	err := gh.db.QueryRow(memberQuery, groupID, userID).Scan(&memberExists)
	if err != nil {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	query := `SELECT p.id, p.user_id, p.group_id, p.content, p.privacy, p.created_at, p.updated_at 
			  FROM posts p 
			  INNER JOIN group_posts gp ON p.id = gp.post_id 
			  WHERE gp.group_id = ? AND p.deleted_at IS NULL 
			  ORDER BY p.created_at DESC`

	rows, err := gh.db.Query(query, groupID)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.GroupID, &post.Content,
			&post.Privacy, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			continue
		}
		posts = append(posts, post)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
