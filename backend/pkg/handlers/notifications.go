package handlers

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "time"

    "social-nework/pkg/models"
)

type NotificationHandler struct {
    NotificationModel *models.NotificationModel
}

func NewNotificationHandler(notificationModel *models.NotificationModel) *NotificationHandler {
    return &NotificationHandler{
        NotificationModel: notificationModel,
    }
}

// GetNotifications returns all notifications for the authenticated user
func (h *NotificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("user_id")
    if userID == nil {
        log.Printf("ERROR: GetNotifications - No user_id in context")
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    
    userIDStr, ok := userID.(string)
    if !ok {
        log.Printf("ERROR: GetNotifications - Invalid user_id type")
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    log.Printf("DEBUG: GetNotifications called for user: %s", userIDStr)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    notifications, err := h.NotificationModel.GetByUserID(ctx, userIDStr)
    if err != nil {
        log.Printf("ERROR: Failed to get notifications for user %s: %v", userIDStr, err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    log.Printf("DEBUG: Retrieved %d notifications for user %s", len(notifications), userIDStr)
    for i, notif := range notifications {
        log.Printf("DEBUG: Notification %d - ID: %v, Type: %v", i+1, notif["id"], notif["type"])
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(notifications)
}

// MarkAsRead marks a notification as read
func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("user_id").(string)
    
    var req struct {
        NotificationID string `json:"notification_id"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    err := h.NotificationModel.MarkAsRead(ctx, req.NotificationID, userID)
    if err != nil {
        log.Printf("Failed to mark notification as read: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
