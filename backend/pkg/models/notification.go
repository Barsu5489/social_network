package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// The Notification struct is defined in models.go, so it should not be redefined here.

type NotificationModel struct {
	DB *sql.DB
}

func (m *NotificationModel) Insert(ctx context.Context, notification Notification) (*Notification, error) {
	query := `
		INSERT INTO notifications (id, user_id, type, reference_id, is_read, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := m.DB.ExecContext(ctx, query,
		notification.ID,
		notification.UserID,
		notification.Type,
		notification.ReferenceID,
		notification.IsRead,
		notification.CreatedAt.Unix(),
	)

	if err != nil {
		return nil, err
	}

	return &notification, nil
}

func (nm *NotificationModel) GetByUserID(ctx context.Context, userID string) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			notifications.id,
			notifications.user_id,
			notifications.type,
			notifications.reference_id,
			notifications.is_read,
			notifications.created_at
		FROM notifications 
		WHERE notifications.user_id = ? 
		AND notifications.deleted_at IS NULL 
		AND notifications.is_read = 0
		ORDER BY notifications.created_at DESC
		LIMIT 50
	`

	rows, err := nm.DB.QueryContext(ctx, query, userID)
	if err != nil {
		log.Printf("ERROR: Failed to query notifications: %v", err)
		return nil, err
	}
	defer rows.Close()

	var notifications []map[string]interface{}

	for rows.Next() {
		var id, notifUserID, notifType, referenceID string
		var isRead bool
		var createdAt int64
		
		err := rows.Scan(&id, &notifUserID, &notifType, &referenceID, &isRead, &createdAt)
		if err != nil {
			log.Printf("ERROR: Failed to scan notification: %v", err)
			continue
		}
		
		message, link := nm.formatNotificationMessage(notifType, referenceID)
		
		// Get actor info based on notification type
		actorNickname, actorAvatar := nm.getActorInfo(ctx, notifType, referenceID)
		
		notification := map[string]interface{}{
			"id":           id,
			"type":         notifType,
			"message":      message,
			"link":         link,
			"is_read":      isRead,
			"created_at":   createdAt,
			"reference_id": referenceID,
		}
		
		if actorNickname != "" {
			notification["actor_nickname"] = actorNickname
		}
		if actorAvatar != "" {
			notification["actor_avatar"] = actorAvatar
		}
		
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (nm *NotificationModel) formatNotificationMessage(notifType, referenceID string) (string, string) {
	switch notifType {
	case "new_follower":
		return "started following you", "/profile/" + referenceID
	case "follow_request":
		return "sent you a follow request", "/profile/" + referenceID
	case "new_like":
		return "liked your post", "/post/" + referenceID
	case "new_comment":
		return "commented on your post", "/post/" + referenceID
	case "new_message":
		return "sent you a message", "/chat/" + referenceID
	case "group_invite":
		return "invited you to join a group", "/groups/" + referenceID
	case "group_join_request":
		return "requested to join your group", "/groups/" + referenceID
	case "event_created":
		return "created a new event in your group", "/events/" + referenceID
	case "group_join_response":
		return "responded to your group join request", "/groups/" + referenceID
	case "group_invitation_response":
		return "responded to your group invitation", "/groups/" + referenceID
	default:
		return "sent you a notification", "#"
	}
}

func (nm *NotificationModel) MarkAsRead(ctx context.Context, notificationID, userID string) error {
	query := `
		UPDATE notifications 
		SET is_read = 1 
		WHERE id = ? AND user_id = ? AND deleted_at IS NULL
	`
	
	result, err := nm.DB.ExecContext(ctx, query, notificationID, userID)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("notification not found or already read")
	}
	
	return nil
}

func (m *NotificationModel) Delete(ctx context.Context, notificationID, userID string) error {
	query := `
		UPDATE notifications
		SET deleted_at = ?
		WHERE id = ? AND user_id = ?
	`

	_, err := m.DB.ExecContext(ctx, query, time.Now().Unix(), notificationID, userID)
	return err
}

func (nm *NotificationModel) getActorInfo(ctx context.Context, notifType, referenceID string) (string, string) {
	var nickname, avatar string
	
	switch notifType {
	case "new_like", "new_comment":
		// For likes and comments, we need to get the actor from the additional data
		// This is more complex, so for now return empty strings
		// You could store actor_id in notifications table for better performance
		return "", ""
	case "new_follower", "follow_request":
		// For follow notifications, reference_id is the follower's user_id
		query := `SELECT nickname, avatar_url FROM users WHERE id = ?`
		nm.DB.QueryRowContext(ctx, query, referenceID).Scan(&nickname, &avatar)
	case "new_message":
		// For messages, you might want to get sender info from chat context
		return "", ""
	default:
		return "", ""
	}
	
	return nickname, avatar
}
