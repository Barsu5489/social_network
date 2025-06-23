package groups

import (
	"encoding/json"
	"net/http"
	"social-nework/pkg/models"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Request to join group
func (gh *GroupHandler) RequestToJoinGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["groupId"]

	var request struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Get group creator
	var creatorID string
	creatorQuery := `SELECT creator_id FROM groups WHERE id = ?`
	err := gh.db.QueryRow(creatorQuery, groupID).Scan(&creatorID)
	if err != nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	// Create invitation from user to themselves (request)
	invitation := models.Invitation {
		ID:         uuid.New().String(),
		InviterID:  request.UserID,
		InviteeID:  request.UserID,
		EntityType: "group",
		EntityID:   groupID,
		Status:     "pending",
		CreatedAt:  time.Now().Unix(),
	}

	query := `INSERT INTO invitations (id, inviter_id, invitee_id, entity_type, entity_id, status, created_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err = gh.db.Exec(query, invitation.ID, invitation.InviterID, invitation.InviteeID,
		invitation.EntityType, invitation.EntityID, invitation.Status, invitation.CreatedAt)
	if err != nil {
		http.Error(w, "Failed to send join request", http.StatusInternalServerError)
		return
	}

	// Notify group creator
	notificationQuery := `INSERT INTO notifications (id, user_id, type, reference_id, created_at) 
						  VALUES (?, ?, ?, ?, ?)`
	notificationID := uuid.New().String()
	_, err = gh.db.Exec(notificationQuery, notificationID, creatorID,
		"group_invite", invitation.ID, time.Now().Unix())
		if err != nil {
			http.Error(w, "Failed to send notification", http.StatusInternalServerError)
			return
		}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invitation)
}
