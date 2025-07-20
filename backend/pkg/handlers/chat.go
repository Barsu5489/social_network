// pkg/handlers/chat_handlers.go
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"social-nework/pkg/auth"
	"social-nework/pkg/models"
	"social-nework/pkg/repository"
	"social-nework/pkg/websocket"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ChatHandler struct {
	chatRepo         *repository.ChatRepository
	messageRepo      *repository.MessageRepository
	groupRepo        *repository.GroupRepository
	notificationRepo *models.NotificationModel
	hub              *websocket.Hub
}

func NewChatHandler(chatRepo *repository.ChatRepository, messageRepo *repository.MessageRepository, groupRepo *repository.GroupRepository, hub *websocket.Hub, notificationRepo *models.NotificationModel) *ChatHandler {
	return &ChatHandler{
		chatRepo:         chatRepo,
		messageRepo:      messageRepo,
		groupRepo:        groupRepo,
		notificationRepo: notificationRepo,
		hub:              hub,
	}
}

// SendMessage handles sending messages to both direct and group chats
func (h *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	vars := mux.Vars(r)
	chatID := vars["chatId"]

	log.Printf("SendMessage: userID=%s, chatID=%s", userID, chatID)

	if chatID == "" {
		log.Printf("SendMessage: Chat ID is required")
		http.Error(w, "Chat ID is required", http.StatusBadRequest)
		return
	}

	var req struct {
		Content string `json:"content"`
		Type    string `json:"type,omitempty"` // text, image, file, etc.
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("SendMessage: Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("SendMessage: content=%s, type=%s", req.Content, req.Type)

	if req.Content == "" {
		log.Printf("SendMessage: Message content is required")
		http.Error(w, "Message content is required", http.StatusBadRequest)
		return
	}

	// Verify user is in chat
	isInChat, err := h.chatRepo.IsUserInChat(chatID, userID)
	if err != nil {
		log.Printf("SendMessage: Error checking chat membership: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !isInChat {
		log.Printf("SendMessage: User %s not in chat %s", userID, chatID)
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// For direct chats, verify follow relationship still exists
	chatType, err := h.chatRepo.GetChatType(chatID)
	if err != nil {
		log.Printf("SendMessage: Error getting chat type: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if chatType == "direct" {
		participants, err := h.chatRepo.GetChatParticipants(chatID)
		if err != nil {
			log.Printf("SendMessage: Error getting chat participants: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		var recipientID string
		for _, participantID := range participants {
			if participantID != userID {
				recipientID = participantID
				break
			}
		}

		if recipientID != "" {
			canChat, err := h.chatRepo.CanUsersChat(userID, recipientID)
			if err != nil {
				log.Printf("SendMessage: Error checking chat permissions: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if !canChat {
				http.Error(w, "Cannot send message: follow relationship required", http.StatusForbidden)
				return
			}
		}
	}

	log.Printf("SendMessage: User verified in chat, creating message")

	// Create message
	message := &models.Message{
		ID:       uuid.New().String(),
		ChatID:   chatID,
		SenderID: userID,
		Content:  req.Content,
		SentAt:   time.Now().Unix(),
	}

	log.Printf("SendMessage: Created message with ID=%s", message.ID)

	// Save message to database
	if err := h.messageRepo.SaveMessage(message); err != nil {
		log.Printf("SendMessage: Error saving message: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get sender info for response
	var sender models.User
	err = h.chatRepo.DB.QueryRow(`
		SELECT first_name, last_name, avatar_url
		FROM users WHERE id = ?`, userID).Scan(
		&sender.FirstName, &sender.LastName, &sender.AvatarURL)
	if err != nil {
		log.Printf("SendMessage: Error getting sender info: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	message.Sender = sender

	log.Printf("SendMessage: Retrieved sender info: %+v", sender)

	// Broadcast message via websocket
	wsMessage := websocket.MessagePayload{
		Type:   "new_message",
		ChatID: chatID,
		Data: map[string]interface{}{
			"id":        message.ID,
			"chat_id":   message.ChatID,
			"sender_id": message.SenderID,
			"content":   message.Content,
			"sent_at":   message.SentAt,
			"sender": map[string]interface{}{
				"first_name": sender.FirstName,
				"last_name":  sender.LastName,
				"avatar_url": sender.AvatarURL,
			},
		},
	}

	log.Printf("SendMessage: Broadcasting WebSocket message to chat %s", chatID)
	h.hub.BroadcastToChatRoom(chatID, wsMessage, userID)

	// Return the saved message
	w.Header().Set("Content-Type", "application/json")
	responseData := map[string]interface{}{
		"message": map[string]interface{}{
			"id":        message.ID,
			"chat_id":   message.ChatID,
			"sender_id": message.SenderID,
			"content":   message.Content,
			"sent_at":   message.SentAt,
			"sender": map[string]interface{}{
				"first_name": sender.FirstName,
				"last_name":  sender.LastName,
				"avatar_url": sender.AvatarURL,
			},
		},
		"success": true,
	}
	
	log.Printf("SendMessage: Returning response: %+v", responseData)
	json.NewEncoder(w).Encode(responseData)
}

// CreateDirectChat creates a direct chat between two users
func (h *ChatHandler) CreateDirectChat(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req struct {
		RecipientID string `json:"recipient_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.RecipientID == "" {
		http.Error(w, "Recipient ID is required", http.StatusBadRequest)
		return
	}

	if req.RecipientID == userID {
		http.Error(w, "Cannot create chat with yourself", http.StatusBadRequest)
		return
	}

	// Check if users can chat (follow relationship required)
	canChat, err := h.chatRepo.CanUsersChat(userID, req.RecipientID)
	if err != nil {
		log.Printf("Error checking chat permissions: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !canChat {
		http.Error(w, "Cannot create chat: users must follow each other or recipient must have public profile", http.StatusForbidden)
		return
	}

	// Create or get existing direct chat
	chatID, err := h.messageRepo.CreateDirectChat(userID, req.RecipientID)
	if err != nil {
		log.Printf("Error creating direct chat: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Initialize chat room in hub if not exists
	h.hub.InitializeChatRoom(chatID, "direct", []string{userID, req.RecipientID})

	response := map[string]interface{}{
		"chat_id": chatID,
		"type":    "direct",
		"created": true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateGroupChat creates a group chat (legacy method, consider using group creation instead)
func (h *ChatHandler) CreateGroupChat(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req struct {
		Name           string   `json:"name"`
		Description    string   `json:"description,omitempty"`
		ParticipantIDs []string `json:"participant_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Group name is required", http.StatusBadRequest)
		return
	}

	if len(req.ParticipantIDs) == 0 {
		http.Error(w, "At least one participant is required", http.StatusBadRequest)
		return
	}

	// Create group chat
	chat, err := h.chatRepo.CreateChat("group", userID)
	if err != nil {
		log.Printf("Error creating group chat: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Add participants
	allParticipants := []string{userID}
	for _, participantID := range req.ParticipantIDs {
		if participantID != userID {
			if err := h.chatRepo.AddParticipant(chat.ID, participantID); err != nil {
				log.Printf("Error adding participant %s: %v", participantID, err)
				continue
			}
			allParticipants = append(allParticipants, participantID)
		}
	}

	// Initialize chat room in hub
	h.hub.InitializeChatRoom(chat.ID, "group", allParticipants)

	response := map[string]interface{}{
		"chat_id":      chat.ID,
		"type":         "group",
		"name":         req.Name,
		"description":  req.Description,
		"participants": allParticipants,
		"created_at":   chat.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetUserChats returns all chats for the authenticated user (both direct and group)
func (h *ChatHandler) GetUserChats(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	// Get direct chats
	log.Printf("GetUserChats: starting for user %s", userID)
	directChats, err := h.chatRepo.GetUserChats(userID)
	if err != nil {
		log.Printf("GetUserChats: Error getting user chats: %v", err)
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("GetUserChats: retrieved direct chats: %v", directChats)

	var enhancedChats []map[string]interface{}

	// Process direct chats
	for _, chat := range directChats {
		if chat.Type == "direct" {
			log.Printf("GetUserChats: processing direct chat: %v", chat)

			// Get participant details (excluding current user for display name)
			var otherParticipant models.User
			participants, err := h.chatRepo.GetChatParticipants(chat.ID)
			if err != nil {
				log.Printf("GetUserChats: Error getting chat participants: %v", err)
				continue
			}
			log.Printf("GetUserChats: retrieved chat participants: %v", participants)

			for _, participantID := range participants {
				if participantID != userID {
					err = h.chatRepo.DB.QueryRow(`
						SELECT first_name, last_name, avatar_url
						FROM users WHERE id = ?`, participantID).Scan(
						&otherParticipant.FirstName, &otherParticipant.LastName, &otherParticipant.AvatarURL)
					if err != nil {
						log.Printf("GetUserChats: Error getting participant info: %v", err)
					}
					break
				}
			}
			log.Printf("GetUserChats: retrieved participant info: %v", otherParticipant)

			// Get last message
			var lastMessage *models.Message
			messages, err := h.messageRepo.GetChatMessages(chat.ID, time.Time{}, 1)
			if err == nil && len(messages) > 0 {
				lastMessage = &messages[0]
			}
			log.Printf("GetUserChats: retrieved last message: %v", lastMessage)

			enhancedChat := map[string]interface{}{
				"id":           chat.ID,
				"type":         chat.Type,
				"created_at":   chat.CreatedAt,
				"name":         otherParticipant.FirstName + " " + otherParticipant.LastName,
				"avatar_url":   otherParticipant.AvatarURL,
				"participants": participants,
			}

			if lastMessage != nil {
				enhancedChat["last_message"] = map[string]interface{}{
					"content":   lastMessage.Content,
					"sender_id": lastMessage.SenderID,
					"sent_at":   lastMessage.SentAt,
					"sender":    lastMessage.Sender,
				}
			}

			enhancedChats = append(enhancedChats, enhancedChat)
		}
	}

	// Get group chats
	rows, err := h.chatRepo.DB.Query(`
		SELECT c.id, c.type, c.created_at, g.id, g.name, g.description
		FROM chats c
		JOIN group_chats gc ON c.id = gc.chat_id
		JOIN groups g ON gc.group_id = g.id
		JOIN group_members gm ON g.id = gm.group_id
		WHERE gm.user_id = ? AND c.type = 'group'
		AND c.deleted_at IS NULL AND gm.deleted_at IS NULL
		ORDER BY c.created_at DESC`, userID)
	if err != nil {
		log.Printf("GetUserChats: Error getting group chats: %v", err)
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
	} else {
		log.Printf("GetUserChats: retrieved group chats")
		defer rows.Close()

		for rows.Next() {
			var chatID, chatType, groupID, groupName, groupDescription string
			var createdAtUnix int64

			if err := rows.Scan(&chatID, &chatType, &createdAtUnix, &groupID, &groupName, &groupDescription); err != nil {
				log.Printf("GetUserChats: Error scanning group chat: %v", err)
				continue
			}

			createdAt := time.Unix(createdAtUnix, 0)
			log.Printf("GetUserChats: scanned group chat: %s, %s, %v, %s, %s, %s", chatID, chatType, createdAt, groupID, groupName, groupDescription)

			// Get participants
			participants, err := h.chatRepo.GetChatParticipants(chatID)
			if err != nil {
				log.Printf("GetUserChats: Error getting chat participants: %v", err)
				continue
			}
			log.Printf("GetUserChats: retrieved chat participants: %v", participants)

			// Get last message
			messages, err := h.messageRepo.GetChatMessages(chatID, time.Time{}, 1)
			var lastMessage *models.Message
			if err == nil && len(messages) > 0 {
				lastMessage = &messages[0]
			}
			log.Printf("GetUserChats: retrieved last message: %v", lastMessage)

			groupChat := map[string]interface{}{
				"id":          chatID,
				"type":        chatType,
				"created_at":  createdAt,
				"name":        groupName,
				"description": groupDescription,
				"group": map[string]interface{}{
					"id":          groupID,
					"name":        groupName,
					"description": groupDescription,
				},
				"participants": participants,
			}

			if lastMessage != nil {
				groupChat["last_message"] = map[string]interface{}{
					"content":   lastMessage.Content,
					"sender_id": lastMessage.SenderID,
					"sent_at":   lastMessage.SentAt,
					"sender":    lastMessage.Sender,
				}
			}

			enhancedChats = append(enhancedChats, groupChat)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"chats": enhancedChats,
	})
}
// AddParticipant adds a user to a group chat
func (h *ChatHandler) AddParticipant(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	vars := mux.Vars(r)
	chatID := vars["chatId"]

	var req struct {
		UserID string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify requester is in chat
	isInChat, err := h.chatRepo.IsUserInChat(chatID, userID)
	if err != nil || !isInChat {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Add participant
	if err := h.chatRepo.AddParticipant(chatID, req.UserID); err != nil {
		log.Printf("Error adding participant: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Add to active chat room if exists
	h.hub.AddUserToChatRoom(chatID, req.UserID)

	// Notify other participants
	notification := websocket.MessagePayload{
		Type:   "participant_added",
		ChatID: chatID,
		Data: map[string]interface{}{
			"user_id":   req.UserID,
			"added_by":  userID,
			"timestamp": time.Now(),
		},
	}
	h.hub.BroadcastToChatRoom(chatID, notification, "")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Participant added successfully",
	})
}

// GetChatMessages returns paginated messages for a specific chat
func (h *ChatHandler) GetChatMessages(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	vars := mux.Vars(r)
	chatID := vars["chatId"]

	if chatID == "" {
		http.Error(w, "Chat ID is required", http.StatusBadRequest)
		return
	}

	// Verify user is in chat
	isInChat, err := h.chatRepo.IsUserInChat(chatID, userID)
	if err != nil {
		log.Printf("Error checking chat membership: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !isInChat {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Parse pagination parameters
	limitStr := r.URL.Query().Get("limit")
	beforeStr := r.URL.Query().Get("before")

	limit := 50 // default
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	var before time.Time
	if beforeStr != "" {
		if parsedTime, err := time.Parse(time.RFC3339, beforeStr); err == nil {
			before = parsedTime
		}
	}

	messages, err := h.messageRepo.GetChatMessages(chatID, before, limit)
	if err != nil {
		log.Printf("Error getting chat messages: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Reverse messages to show oldest first
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"messages": messages,
		"has_more": len(messages) == limit,
	})
}

// GetGroupChatForGroup returns the chat ID for a specific group (helper method)
func (h *ChatHandler) GetGroupChatForGroup(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	vars := mux.Vars(r)
	groupID := vars["groupId"]

	if groupID == "" {
		http.Error(w, "Group ID is required", http.StatusBadRequest)
		return
	}

	// Verify user is member of group
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
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get group info
	var groupName, groupDescription string
	err = h.chatRepo.DB.QueryRow(`
		SELECT name, description FROM groups WHERE id = ?`, groupID).Scan(&groupName, &groupDescription)
	if err != nil {
		log.Printf("Error getting group info: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get participants
	participants, err := h.chatRepo.GetChatParticipants(chatID)
	if err != nil {
		log.Printf("Error getting chat participants: %v", err)
		participants = []string{} // Return empty array on error
	}

	// Initialize chat room in hub if needed
	h.hub.InitializeChatRoom(chatID, "group", participants)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"chat_id": chatID,
		"group": map[string]interface{}{
			"id":          groupID,
			"name":        groupName,
			"description": groupDescription,
		},
		"participants": participants,
		"type":         "group",
	})
 	}

// RegisterChatRoutes registers all chat-related routes
func RegisterChatRoutes(router *mux.Router, handler *ChatHandler) {
	chatRouter := router.PathPrefix("/api/chats").Subrouter()

	// Chat management - all protected with auth middleware
	chatRouter.HandleFunc("", auth.RequireAuth(handler.GetUserChats)).Methods("GET")
	chatRouter.HandleFunc("/direct", auth.RequireAuth(handler.CreateDirectChat)).Methods("POST")
	chatRouter.HandleFunc("/group", auth.RequireAuth(handler.CreateGroupChat)).Methods("POST")

	// Message management - also protected
	chatRouter.HandleFunc("/{chatId}/messages", auth.RequireAuth(handler.GetChatMessages)).Methods("GET")
	chatRouter.HandleFunc("/{chatId}/messages", auth.RequireAuth(handler.SendMessage)).Methods("POST")
	chatRouter.HandleFunc("/{chatId}/participants", auth.RequireAuth(handler.AddParticipant)).Methods("POST")

	// Group chat helper route
	groupRouter := router.PathPrefix("/api/groups").Subrouter()
	groupRouter.HandleFunc("/{groupId}/chat", auth.RequireAuth(handler.GetGroupChatForGroup)).Methods("GET")
}
