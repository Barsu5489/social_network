package groups

import (
	"encoding/json"
	"log"
	"net/http"
	"social-nework/pkg/models"
	"time"

	"github.com/google/uuid"
)

// Create a new group
func (gh *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var group models.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	group.ID = uuid.New().String()
	group.CreatedAt = time.Now().Unix()
	group.UpdatedAt = time.Now().Unix()

	query := `INSERT INTO groups (id, name, description, creator_id, is_private, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := gh.db.Exec(query, group.ID, group.Name, group.Description,
		group.CreatorID, group.IsPrivate, group.CreatedAt, group.UpdatedAt)
	if err != nil {
		log.Printf("Error creating group: %v", err)
		http.Error(w, "Failed to create group", http.StatusInternalServerError)
		return
	}

	// Add creator as admin member
	memberQuery := `INSERT INTO group_members (id, group_id, user_id, role, joined_at) 
					VALUES (?, ?, ?, ?, ?)`
	memberID := uuid.New().String()
	_, err = gh.db.Exec(memberQuery, memberID, group.ID, group.CreatorID, "admin", time.Now().Unix())
	if err != nil {
		log.Printf("Error adding creator as member: %v", err)
		http.Error(w, "Failed to add creator as member", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(group)
}
