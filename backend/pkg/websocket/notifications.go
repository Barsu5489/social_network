package websocket

import (
    "log"
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
    h.mu.RLock()
    client, exists := h.Clients[userID]
    h.mu.RUnlock()

    if !exists {
        log.Printf("SendNotification: User %s not connected, notification will be stored only", userID)
        return
    }

    payload := NotificationPayload{
        Type:         "notification",
        Notification: notification,
        Data:         additionalData,
    }

    select {
    case client.Send <- MessagePayload{
        Type: "notification",
        Data: payload,
    }:
        log.Printf("SendNotification: Real-time notification sent to user %s", userID)
    default:
        log.Printf("SendNotification: Client %s send buffer full for notification", userID)
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
