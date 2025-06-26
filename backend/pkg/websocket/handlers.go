	// pkg/websocket/handlers.go
	package websocket

	import (
		"log"
		"time"

		"social-nework/pkg/models"

		"github.com/google/uuid"
	)

	func (h *Hub) handleNewMessage(msg MessagePayload) {
		// Validate chat exists and user is participant
		if !h.validateChatParticipation(msg.ChatID, msg.SenderID) {
			log.Printf("User %s not in chat %s", msg.SenderID, msg.ChatID)
			return
		}

		// Save to database
		message := models.Message{
			ID:       uuid.New().String(),
			ChatID:   msg.ChatID,
			SenderID: msg.SenderID,
			Content:  msg.Content,
			SentAt:   time.Now(),
		}

		if err := h.messageRepo.SaveMessage(&message); err != nil {
			log.Printf("Error saving message: %v", err)
			return
		}

		// Broadcast to chat participants
		chat, ok := h.ChatRooms[msg.ChatID]
		if !ok {
			log.Printf("Chat room %s not active", msg.ChatID)
			return
		}

		response := MessagePayload{
			Type:      "message",
			ChatID:    msg.ChatID,
			SenderID:  msg.SenderID,
			Content:   msg.Content,
			Timestamp: message.SentAt,
		}

		for userID, client := range chat.Members {
			if userID == msg.SenderID {
				continue // Don't echo back to sender
			}
			select {
			case client.Send <- response:
			default:
				log.Printf("Client %s send buffer full", userID)
			}
		}
	}
