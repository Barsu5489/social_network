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
        log.Printf("ERROR: GetNotifications - Invalid user_id type: %T", userID)
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

    // Return empty array instead of null if no notifications
    if notifications == nil {
        notifications = []map[string]interface{}{}
    }

    log.Printf("DEBUG: Retrieved %d notifications for user %s", len(notifications), userIDStr)
    
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(notifications); err != nil {
        log.Printf("ERROR: Failed to encode notifications response: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    
    log.Printf("SUCCESS: Notifications response sent for user %s", userIDStr)
}

// MarkNotificationAsRead marks a notification as read
func (h *NotificationHandler) MarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value("user_id").(string)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var req struct {
        NotificationID string `json:"notification_id"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        log.Printf("ERROR: Invalid request body: %v", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if req.NotificationID == "" {
        http.Error(w, "Notification ID is required", http.StatusBadRequest)
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    log.Printf("DEBUG: Marking notification %s as read for user %s", req.NotificationID, userID)

    err := h.NotificationModel.MarkAsRead(ctx, req.NotificationID, userID)
    if err != nil {
        log.Printf("ERROR: Failed to mark notification as read: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    log.Printf("SUCCESS: Notification %s marked as read", req.NotificationID)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
