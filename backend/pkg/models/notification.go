package models

import (
	"context"
	"database/sql"
	"log"
	"strings"
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
		SELECT n.id, n.type, n.reference_id, n.created_at
		FROM notifications 
		WHERE n.user_id = ? AND n.is_read = false AND n.deleted_at IS NULL
		ORDER BY n.created_at DESC
	`

	rows, err := nm.DB.QueryContext(ctx, query, userID)
	if err != nil {
		log.Printf("ERROR: Failed to query notifications: %v", err)
		return []map[string]interface{}{}, nil
	}
	defer rows.Close()

	var notifications []map[string]interface{}

	for rows.Next() {
		var id, notifType, referenceID string
		var createdAt time.Time

		err := rows.Scan(&id, &notifType, &referenceID, &createdAt)
		if err != nil {
			log.Printf("ERROR: Failed to scan notification: %v", err)
			continue
		}

		// Get user info based on notification type
		var firstName, lastName, avatarURL string

		switch notifType {
		case "new_follower":
			// referenceID is the follower's user ID
			err = nm.DB.QueryRowContext(ctx, 
				"SELECT first_name, last_name, COALESCE(avatar_url, '') FROM users WHERE id = ?", 
				referenceID).Scan(&firstName, &lastName, &avatarURL)
		case "new_like":
			// referenceID is the post ID, find the most recent liker
			err = nm.DB.QueryRowContext(ctx, `
				SELECT u.first_name, u.last_name, COALESCE(u.avatar_url, '') 
				FROM users u 
				JOIN likes l ON u.id = l.user_id 
				WHERE l.likeable_id = ? AND l.likeable_type = 'post' AND l.deleted_at IS NULL 
				ORDER BY l.created_at DESC LIMIT 1
			`, referenceID).Scan(&firstName, &lastName, &avatarURL)
		case "new_comment":
			// referenceID is the post ID, find the most recent commenter
			err = nm.DB.QueryRowContext(ctx, `
				SELECT u.first_name, u.last_name, COALESCE(u.avatar_url, '') 
				FROM users u 
				JOIN comments c ON u.id = c.user_id 
				WHERE c.post_id = ? AND c.deleted_at IS NULL 
				ORDER BY c.created_at DESC LIMIT 1
			`, referenceID).Scan(&firstName, &lastName, &avatarURL)
		case "new_message":
			// referenceID is the chat ID, find the most recent sender
			err = nm.DB.QueryRowContext(ctx, `
				SELECT u.first_name, u.last_name, COALESCE(u.avatar_url, '') 
				FROM users u 
				JOIN messages m ON u.id = m.sender_id 
				WHERE m.chat_id = ? AND m.deleted_at IS NULL 
				ORDER BY m.sent_at DESC LIMIT 1
			`, referenceID).Scan(&firstName, &lastName, &avatarURL)
		default:
			firstName, lastName, avatarURL = "Unknown", "User", ""
		}

		if err != nil {
			log.Printf("ERROR: Failed to get user info for notification %s: %v", id, err)
			firstName, lastName, avatarURL = "Unknown", "User", ""
		}

		// Format user name
		userName := strings.TrimSpace(firstName + " " + lastName)
		if userName == "" {
			userName = "Unknown User"
		}

		// Generate message and href based on type
		message, href := nm.formatNotificationMessage(notifType, referenceID)

		notification := map[string]interface{}{
			"id":   id,
			"type": notifType,
			"user": map[string]interface{}{
				"name":   userName,
				"avatar": avatarURL,
			},
			"message": message,
			"time":    createdAt.Format("2 Jan 2006 15:04"),
			"href":    href,
		}

		notifications = append(notifications, notification)
	}

	log.Printf("DEBUG: Retrieved %d notifications for user %s", len(notifications), userID)
	return notifications, nil
}

func (nm *NotificationModel) formatNotificationMessage(notifType, referenceID string) (string, string) {
	switch notifType {
	case "new_follower":
		return "started following you", "/profile/" + referenceID
	case "new_like":
		return "liked your post", "/posts/" + referenceID
	case "new_comment":
		return "commented on your post", "/posts/" + referenceID
	case "new_message":
		return "sent you a message", "/chat/" + referenceID
	case "follow_request":
		return "sent you a follow request", "/profile/" + referenceID
	case "group_invite":
		return "invited you to join a group", "/groups/" + referenceID
	case "group_join_request":
		return "requested to join your group", "/groups/" + referenceID
	default:
		return "sent you a notification", "#"
	}
}

func (m *NotificationModel) MarkAsRead(ctx context.Context, notificationID, userID string) error {
	query := `
		UPDATE notifications
		SET is_read = 1
		WHERE id = ? AND user_id = ?
	`
	
	_, err := m.DB.ExecContext(ctx, query, notificationID, userID)
	return err
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
