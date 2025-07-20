package groups

import (
	"encoding/json"
	"log"
	"net/http"
	"social-nework/pkg/models"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Request to join group
func (gh *GroupHandler) RequestToJoinGroup(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	groupID := vars["groupId"]

	var request struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validation fields
	if request.UserID == "" {
		http.Error(w, "Missing required field: user_id", http.StatusBadRequest)
		return
	}
	if request.UserID != userID {
		http.Error(w, "User ID mismatch with authenticated user", http.StatusForbidden)
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
	invitation := models.Invitation{
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
	notification := models.Notification{
		UserID:      creatorID,
		Type:        "group_join_request",
		ReferenceID: request.UserID,
		IsRead:      false,
		CreatedAt:   time.Now(),
	}

	log.Printf("DEBUG: Creating group join request notification - CreatorID: %s, RequesterID: %s, GroupID: %s",
		creatorID, request.UserID, groupID)

	_, err = gh.NotificationModel.Insert(r.Context(), notification)
	if err != nil {
		log.Printf("ERROR: Failed to create group join request notification: %v", err)
	} else {
		log.Printf("SUCCESS: Group join request notification created for creator %s", creatorID)

		// Send real-time notification
		if gh.h != nil {
			log.Printf("DEBUG: Sending real-time group join request notification to creator %s", creatorID)
			gh.h.SendNotification(creatorID, notification, map[string]interface{}{
				"group_id":     groupID,
				"requester_id": request.UserID,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invitation)
}
