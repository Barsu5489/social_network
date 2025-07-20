// pkg/websocket/handlers.go
package websocket

import (
	"context"
	"log"
	"time"

	"social-nework/pkg/models"

	"github.com/google/uuid"
)

func (h *Hub) handleNewMessage(msg MessagePayload) {
	log.Printf("DEBUG: handleNewMessage called - ChatID: %s, SenderID: %s", msg.ChatID, msg.SenderID)
	
	// Validate chat exists and user is participant
	if !h.validateChatParticipation(msg.ChatID, msg.SenderID) {
		log.Printf("ERROR: User %s not in chat %s", msg.SenderID, msg.ChatID)
		return
	}

	// Save to database
	message := models.Message{
		ID:       uuid.New().String(),
		ChatID:   msg.ChatID,
		SenderID: msg.SenderID,
		Content:  msg.Content,
		SentAt:   time.Now().Unix(),
	}

	if err := h.messageRepo.SaveMessage(&message); err != nil {
		log.Printf("ERROR: Error saving message: %v", err)
		return
	}
	
	log.Printf("SUCCESS: Message saved to database - ID: %s", message.ID)

	// Get chat participants for notifications
	participants, err := h.chatRepo.GetChatParticipants(msg.ChatID)
	if err != nil {
		log.Printf("ERROR: Failed to get chat participants: %v", err)
	} else {
		log.Printf("DEBUG: Chat participants: %v", participants)
		
		// Create notifications for other participants
		for _, participantID := range participants {
			if participantID != msg.SenderID {
				log.Printf("DEBUG: Creating message notification for participant: %s", participantID)
				
				notification := models.Notification{
					ID:          uuid.New().String(),
					UserID:      participantID,
					Type:        "new_message",
					ReferenceID: message.ID,
					IsRead:      false,
					CreatedAt:   time.Now(),
				}
				
				// Save notification to database
				if h.notificationModel != nil {
					_, err := h.notificationModel.Insert(context.Background(), notification)
					if err != nil {
						log.Printf("ERROR: Failed to save message notification: %v", err)
					} else {
						log.Printf("SUCCESS: Message notification saved for user: %s", participantID)
						
						// Send real-time notification
						h.SendNotification(participantID, notification, map[string]interface{}{
							"chat_id":    msg.ChatID,
							"sender_id":  msg.SenderID,
							"message_id": message.ID,
						})
					}
				}
			}
		}
	}

	// Broadcast to chat participants
	chat, ok := h.ChatRooms[msg.ChatID]
	if !ok {
		log.Printf("WARNING: Chat room %s not active", msg.ChatID)
		return
	}

	// Get full message with sender details for broadcasting
	fullMessage, err := h.messageRepo.GetMessageByID(message.ID)
	if err != nil {
		log.Printf("ERROR: Failed to get full message details: %v", err)
		return
	}

	broadcastMsg := MessagePayload{
		Type:   "new_message",
		ChatID: msg.ChatID,
		Data:   fullMessage,
	}

	log.Printf("DEBUG: Broadcasting message to %d participants", len(chat.Members))
	for participantID := range chat.Members {
		if client, exists := h.Clients[participantID]; exists {
			select {
			case client.Send <- broadcastMsg:
				log.Printf("SUCCESS: Message broadcasted to participant: %s", participantID)
			default:
				log.Printf("ERROR: Failed to send message to participant %s - buffer full", participantID)
			}
		} else {
			log.Printf("WARNING: Participant %s not connected", participantID)
		}
	}
}

func (h *Hub) handleHistoryRequest(msg MessagePayload) {
	// Validate chat participation
	if !h.validateChatParticipation(msg.ChatID, msg.SenderID) {
		return
	}

	var before time.Time
	if data, ok := msg.Data.(map[string]interface{}); ok {
		if ts, ok := data["before"].(string); ok {
			before, _ = time.Parse(time.RFC3339, ts)
		}
	}

	messages, err := h.messageRepo.GetChatMessages(msg.ChatID, before, 50)
	if err != nil {
		log.Printf("Error getting chat history: %v", err)
		return
	}

	client, ok := h.Clients[msg.SenderID]
	if !ok {
		return
	}

	response := MessagePayload{
		Type:   "history_response",
		ChatID: msg.ChatID,
		Data:   messages,
	}

	select {
	case client.Send <- response:
	default:
		log.Printf("Client %s send buffer full", msg.SenderID)
	}
}

func (h *Hub) validateChatParticipation(chatID, userID string) bool {
	// Check in-memory first
	if chat, ok := h.ChatRooms[chatID]; ok {
		_, exists := chat.Members[userID]
		return exists
	}

	// Fallback to database check
	var exists bool
	err := h.db.QueryRow(`
			SELECT EXISTS(
				SELECT 1 FROM chat_participants 
				WHERE chat_id = ? AND user_id = ? AND deleted_at IS NULL
			)`, chatID, userID).Scan(&exists)
	if err != nil {
		log.Printf("Error validating chat participation: %v", err)
		return false
	}

	return exists
}
