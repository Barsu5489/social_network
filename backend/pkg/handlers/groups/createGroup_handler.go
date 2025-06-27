// pkg/handlers/group_handlers.go - Updated methods
package groups

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"social-nework/pkg/models"
	"social-nework/pkg/websocket"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		IsPrivate   bool   `json:"is_private"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Group name is required", http.StatusBadRequest)
		return
	}

	group := &models.Group{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		CreatorID:   userID,
		IsPrivate:   req.IsPrivate,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	// Create group with associated chat
	if err := h.groupRepo.CreateGroupWithChat(group); err != nil {
		log.Printf("Error creating group with chat: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get the group's chat ID
	chatID, err := h.groupRepo.GetGroupChatID(group.ID)
	if err != nil {
		log.Printf("Error getting group chat ID: %v", err)
		// Don't fail the request, just log the error
	}

	// Initialize chat room in websocket hub if chat was created successfully
	if chatID != "" && h.h != nil {
		h.h.InitializeChatRoom(chatID, "group", []string{userID})
	}

	response := map[string]interface{}{
		"id":          group.ID,
		"name":        group.Name,
		"description": group.Description,
		"creator_id":  group.CreatorID,
		"is_private":  group.IsPrivate,
		"created_at":  group.CreatedAt,
		"chat_id":     chatID, // Include chat ID in response
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *GroupHandler) JoinGroup(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	vars := mux.Vars(r)
	groupID := vars["groupId"]

	if groupID == "" {
		http.Error(w, "Group ID is required", http.StatusBadRequest)
		return
	}

	// Check if user is already a member
	isMember, err := h.groupRepo.IsUserMember(groupID, userID)
	if err != nil {
		log.Printf("Error checking group membership: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if isMember {
		http.Error(w, "Already a member of this group", http.StatusBadRequest)
		return
	}

	// Add user to group
	if err := h.groupRepo.AddMember(groupID, userID, "member"); err != nil {
		log.Printf("Error adding user to group: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Add user to group chat
	if err := h.chatRepo.AddMemberToGroupChat(groupID, userID); err != nil {
		log.Printf("Error adding user to group chat: %v", err)
		// Don't fail the request, just log the error
	}

	// Add user to websocket chat room if they're online
	if h.h != nil {
		if chatID, err := h.groupRepo.GetGroupChatID(groupID); err == nil {
			h.h.AddUserToChatRoom(chatID, userID)

			// Notify other group members
			notification := websocket.MessagePayload{
				Type:   "member_joined",
				ChatID: chatID,
				Data: map[string]interface{}{
					"user_id":   userID,
					"group_id":  groupID,
					"timestamp": time.Now(),
				},
			}
			h.h.BroadcastToChatRoom(chatID, notification, userID)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Successfully joined group",
	})
}

func (h *GroupHandler) LeaveGroup(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	vars := mux.Vars(r)
	groupID := vars["groupId"]

	if groupID == "" {
		http.Error(w, "Group ID is required", http.StatusBadRequest)
		return
	}

	// Check if user is a member
	isMember, err := h.groupRepo.IsUserMember(groupID, userID)
	if err != nil {
		log.Printf("Error checking group membership: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !isMember {
		http.Error(w, "Not a member of this group", http.StatusBadRequest)
		return
	}

	// Get chat ID before removing user
	chatID, _ := h.groupRepo.GetGroupChatID(groupID)

	// Remove user from group
	if err := h.groupRepo.RemoveMember(groupID, userID); err != nil {
		log.Printf("Error removing user from group: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Remove user from group chat
	if err := h.chatRepo.RemoveMemberFromGroupChat(groupID, userID); err != nil {
		log.Printf("Error removing user from group chat: %v", err)
	}

	// Remove user from websocket chat room
	if h.h != nil && chatID != "" {
		h.h.RemoveUserFromChatRoom(chatID, userID)

		// Notify remaining group members
		notification := websocket.MessagePayload{
			Type:   "member_left",
			ChatID: chatID,
			Data: map[string]interface{}{
				"user_id":   userID,
				"group_id":  groupID,
				"timestamp": time.Now(),
			},
		}
		h.h.BroadcastToChatRoom(chatID, notification, "")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Successfully left group",
	})
}

// Add this method to get group chat information
func (h *GroupHandler) GetGroupChat(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	vars := mux.Vars(r)
	groupID := vars["groupId"]

	if groupID == "" {
		http.Error(w, "Group ID is required", http.StatusBadRequest)
		return
	}

	// Check if user is a member of the group
	isMember, err := h.groupRepo.IsUserMember(groupID, userID)
	if err != nil {
		log.Printf("Error checking group membership: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !isMember {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Get group chat ID
	chatID, err := h.groupRepo.GetGroupChatID(groupID)
	if err != nil {
		log.Printf("Error getting group chat ID: %v", err)
		http.Error(w, "Group chat not found", http.StatusNotFound)
		return
	}

	// Get chat participants
	participants, err := h.chatRepo.GetChatParticipants(chatID)
	if err != nil {
		log.Printf("Error getting chat participants: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"chat_id":      chatID,
		"group_id":     groupID,
		"type":         "group",
		"participants": participants,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
