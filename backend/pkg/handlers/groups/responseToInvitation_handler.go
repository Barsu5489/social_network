package groups

import (
	"encoding/json"
	"net/http"
	"social-nework/pkg/models"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Accept/Decline group invitation
func (gh *GroupHandler) RespondToInvitation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	invitationID := vars["id"]

	var response struct {
		Status string `json:"status"` // accepted or declined
	}
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Update invitation status
	updateQuery := `UPDATE invitations SET status = ? WHERE id = ?`
	_, err := gh.db.Exec(updateQuery, response.Status, invitationID)
	if err != nil {
		http.Error(w, "Failed to update invitation", http.StatusInternalServerError)
		return
	}

	// If accepted, add user to group
	if response.Status == "accepted" {
		var invitation models.Invitation
		inviteQuery := `SELECT invitee_id, entity_id FROM invitations WHERE id = ?`
		err := gh.db.QueryRow(inviteQuery, invitationID).Scan(&invitation.InviteeID, &invitation.EntityID)
		if err != nil {
			http.Error(w, "Invitation not found", http.StatusNotFound)
			return
		}

		memberQuery := `INSERT INTO group_members (id, group_id, user_id, role, joined_at) 
						VALUES (?, ?, ?, ?, ?)`
		memberID := uuid.New().String()
		_, err = gh.db.Exec(memberQuery, memberID, invitation.EntityID, invitation.InviteeID,
			"member", time.Now().Unix())
		if err != nil {
			http.Error(w, "Failed to add member to group", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
