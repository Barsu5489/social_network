package groups

import (
	"encoding/json"
	"log"
	"net/http"
	"social-nework/pkg/models"
	"time"

	"github.com/google/uuid"
)

// Invite user to group
func (gh *GroupHandler) InviteToGroup(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var invitation models.Invitation
	if err := json.NewDecoder(r.Body).Decode(&invitation); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if invitation.InviteeID == userID {
		http.Error(w, "You can't invite yourself", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if invitation.InviteeID == "" {
		http.Error(w, "Missing required field: invitee_id", http.StatusBadRequest)
		return
	}
	if invitation.EntityID == "" {
		http.Error(w, "Missing required field: entity_id (group ID)", http.StatusBadRequest)
		return
	}

	// Set inviter ID from the authenticated user
	invitation.InviterID = userID

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
	notification := models.Notification{
		UserID:      invitation.InviteeID,
		Type:        "group_invite",
		ReferenceID: invitation.EntityID,
		IsRead:      false,
		CreatedAt:   time.Now(),
	}

	log.Printf("DEBUG: Creating group invitation notification - UserID: %s, Type: %s, GroupID: %s", 
		notification.UserID, notification.Type, notification.ReferenceID)

	_, err = gh.NotificationModel.Insert(r.Context(), notification)
	if err != nil {
		log.Printf("ERROR: Failed to create group invitation notification: %v", err)
		http.Error(w, "Failed to create notification for group invitation", http.StatusInternalServerError)
		return
	} else {
		log.Printf("SUCCESS: Group invitation notification created")
		
		// Send real-time notification
		if gh.h != nil {
			log.Printf("DEBUG: Sending real-time group invitation notification to user %s", invitation.InviteeID)
			gh.h.SendNotification(invitation.InviteeID, notification, map[string]interface{}{
				"group_id": invitation.EntityID,
				"inviter_id": invitation.InviterID,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invitation)
}
