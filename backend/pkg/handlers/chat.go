// pkg/handlers/chat_handlers.go
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"social-nework/pkg/models"
	"social-nework/pkg/repository"
	"social-nework/pkg/websocket"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ChatHandler struct {
	chatRepo    *repository.ChatRepository
	messageRepo *repository.MessageRepository
	groupRepo   *repository.GroupRepository // Add this for group chat integration
	hub         *websocket.Hub
}

func NewChatHandler(chatRepo *repository.ChatRepository, messageRepo *repository.MessageRepository, groupRepo *repository.GroupRepository, hub *websocket.Hub) *ChatHandler {
	return &ChatHandler{
		chatRepo:    chatRepo,
		messageRepo: messageRepo,
		groupRepo:   groupRepo,
		hub:         hub,
	}
}

// SendMessage handles sending messages to both direct and group chats
func (h *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	vars := mux.Vars(r)
	chatID := vars["chatId"]

	if chatID == "" {
		http.Error(w, "Chat ID is required", http.StatusBadRequest)
		return
	}

	var req struct {
		Content string `json:"content"`
		Type    string `json:"type,omitempty"` // text, image, file, etc.
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		http.Error(w, "Message content is required", http.StatusBadRequest)
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

	// Create message
	message := &models.Message{
		ID:       uuid.New().String(),
		ChatID:   chatID,
		SenderID: userID,
		Content:  req.Content,
		SentAt:   time.Now(),
	}

	// Save message to database
	if err := h.messageRepo.SaveMessage(message); err != nil {
		log.Printf("Error saving message: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get sender info for websocket broadcast
	var sender models.User
	err = h.chatRepo.DB.QueryRow(`
		SELECT first_name, last_name, avatar_url 
		FROM users WHERE id = ?`, userID).Scan(
		&sender.FirstName, &sender.LastName, &sender.AvatarURL)
	if err != nil {
		log.Printf("Error getting sender info: %v", err)
	}
	message.Sender = sender

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

	h.hub.BroadcastToChatRoom(chatID, wsMessage, userID)

	// Return the saved message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
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
	})
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
