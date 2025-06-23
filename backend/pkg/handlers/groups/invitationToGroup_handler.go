package groups

import (
	"encoding/json"
	"net/http"
	"social-nework/pkg/models"
	"time"

	"github.com/google/uuid"
)

// Invite user to group
func (gh *GroupHandler) InviteToGroup(w http.ResponseWriter, r *http.Request) {
	var invitation models.Invitation
	if err := json.NewDecoder(r.Body).Decode(&invitation); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Check if inviter is member of the group
	memberQuery := `SELECT id FROM group_members WHERE group_id = ? AND user_id = ? AND deleted_at IS NULL`
	var memberExists string
	err := gh.db.QueryRow(memberQuery, invitation.EntityID, invitation.InviterID).Scan(&memberExists)
	if err != nil {
		http.Error(w, "You must be a member to invite others", http.StatusForbidden)
		return
	}

	invitation.ID = uuid.New().String()
	invitation.EntityType = "group"
	invitation.Status = "pending"
	invitation.CreatedAt = time.Now().Unix()

	query := `INSERT INTO invitations (id, inviter_id, invitee_id, entity_type, entity_id, status, created_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err = gh.db.Exec(query, invitation.ID, invitation.InviterID, invitation.InviteeID,
		invitation.EntityType, invitation.EntityID, invitation.Status, invitation.CreatedAt)
	if err != nil {
		http.Error(w, "Failed to send invitation", http.StatusInternalServerError)
		return
	}

	// Create notification
	notificationQuery := `INSERT INTO notifications (id, user_id, type, reference_id, created_at) 
						  VALUES (?, ?, ?, ?, ?)`
	notificationID := uuid.New().String()
	_, err = gh.db.Exec(notificationQuery, notificationID, invitation.InviteeID,
		"group_invite", invitation.ID, time.Now().Unix())
		if err != nil {
			http.Error(w, "Failed to create notification for group invitation", http.StatusInternalServerError)
			return
		}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invitation)
}
