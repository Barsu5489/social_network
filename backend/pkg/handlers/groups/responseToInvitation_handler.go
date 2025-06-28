package groups

import (
	"encoding/json"
	"fmt"
	"net/http"
	"social-nework/pkg/models"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Accept/Decline group invitation
func (gh *GroupHandler) RespondToInvitation(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	invitationID := vars["id"]
	if invitationID == "" {
		http.Error(w, "Missing invitation ID", http.StatusBadRequest)
		return
	}

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

	var invitation models.Invitation
	inviteQuery := `SELECT inviter_id, invitee_id, entity_id FROM invitations WHERE id = ?`
	err = gh.db.QueryRow(inviteQuery, invitationID).Scan(&invitation.InviterID, &invitation.InviteeID, &invitation.EntityID)
	if err != nil {
		http.Error(w, "Invitation not found", http.StatusNotFound)
		return
	}

	// If accepted, add user to group
	if response.Status == "accepted" {
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

	var notification models.Notification
	if invitation.InviterID == invitation.InviteeID {
		// Scenario 1: User requested to join a group (inviter_id == invitee_id)
		notification = models.Notification{
			UserID:      invitation.InviteeID,
			Type:        "group_join_response",
			ReferenceID: invitation.EntityID,
			IsRead:      false,
			CreatedAt:   time.Now(),
		}
	} else {
		// Scenario 2: Direct group invitation (inviter_id != invitee_id)
		notification = models.Notification{
			UserID:      invitation.InviterID,
			Type:        "group_invitation_response",
			ReferenceID: invitation.EntityID,
			IsRead:      false,
			CreatedAt:   time.Now(),
		}
	}

	_, err = gh.NotificationModel.Insert(r.Context(), notification)
	if err != nil {
		fmt.Printf("Failed to create notification for group invitation response: %v\n", err)
		// Decide if you should rollback the invitation status update or just log the error
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
