package websocket

import (
	"log"
	"time"
	"social-nework/pkg/models"
)

// NotificationPayload represents a real-time notification
type NotificationPayload struct {
	Type         string                 `json:"type"`
	Notification models.Notification    `json:"notification"`
	Data         map[string]interface{} `json:"data,omitempty"`
}

// SendNotification sends a real-time notification to a specific user
func (h *Hub) SendNotification(userID string, notification models.Notification, additionalData map[string]interface{}) {
	log.Printf("DEBUG: SendNotification called - UserID: %s, Type: %s, NotificationID: %s",
		userID, notification.Type, notification.ID)
	
	h.mu.RLock()
	client, exists := h.Clients[userID]
	h.mu.RUnlock()

	if !exists {
		log.Printf("WARNING: SendNotification: User %s not connected, notification will be stored only", userID)
		return
	}

	log.Printf("DEBUG: User %s is connected, preparing notification payload", userID)

	// Ensure created_at is properly formatted
	createdAtTime := notification.CreatedAt
	if createdAtTime.IsZero() {
		createdAtTime = time.Now()
	}

	payload := NotificationPayload{
		Type:         "notification",
		Notification: notification,
		Data:         additionalData,
	}

	log.Printf("DEBUG: Sending notification payload to user %s: Type=%s, NotificationID=%s, AdditionalData=%+v",
		userID, notification.Type, notification.ID, additionalData)

	messagePayload := MessagePayload{
		Type: "notification",
		Data: payload,
	}

	select {
	case client.Send <- messagePayload:
		log.Printf("SUCCESS: Real-time notification sent to user %s (NotificationID: %s)", userID, notification.ID)
	default:
		log.Printf("ERROR: Client %s send buffer full for notification %s", userID, notification.ID)
	}
}

// BroadcastNotificationToGroup sends a notification to all members of a group
func (h *Hub) BroadcastNotificationToGroup(groupID string, notification models.Notification, excludeUserID string, additionalData map[string]interface{}) {
	// Get group members from database using the correct method
	members, err := h.chatRepo.GetGroupChatMembers(groupID)
	if err != nil {
		log.Printf("BroadcastNotificationToGroup: Error getting group members: %v", err)
		return
	}

	for _, memberID := range members {
		if memberID != excludeUserID {
			h.SendNotification(memberID, notification, additionalData)
		}
	}
}
