package groups

import (
	"encoding/json"
	"net/http"
	"social-nework/pkg/models"
)

// Get all groups (for browsing)
func (gh *GroupHandler) GetAllGroups(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, name, description, creator_id, is_private, created_at, updated_at 
			  FROM groups WHERE deleted_at IS NULL ORDER BY created_at DESC`

	rows, err := gh.db.Query(query)
	if err != nil {
		http.Error(w, "Failed to fetch groups", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var groups []models.Group
	for rows.Next() {
		var group models.Group
		var isPrivateInt int
		err := rows.Scan(&group.ID, &group.Name, &group.Description, &group.CreatorID,
			&isPrivateInt, &group.CreatedAt, &group.UpdatedAt)
		if err != nil {
			continue
		}
		group.IsPrivate = isPrivateInt == 1
		groups = append(groups, group)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}
