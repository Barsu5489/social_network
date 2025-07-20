package models

import (
	"context"
	"database/sql"
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
	// First get all notifications for the user
	query := `
		SELECT id, type, reference_id, created_at
		FROM notifications 
		WHERE user_id = ? AND is_read = false AND deleted_at IS NULL
		ORDER BY created_at DESC
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
		var userQuery string
		var queryParam string

		switch notifType {
		case "new_follower":
			userQuery = "SELECT first_name, last_name, avatar_url FROM users WHERE id = ?"
			queryParam = referenceID
		case "new_like":
			// For likes, we need to find who liked the post
			userQuery = `
				SELECT u.first_name, u.last_name, u.avatar_url 
				FROM users u 
				JOIN likes l ON u.id = l.user_id 
				WHERE l.likeable_id = ? AND l.deleted_at IS NULL 
				ORDER BY l.created_at DESC LIMIT 1
			`
			queryParam = referenceID
		case "new_comment":
			// For comments, we need to find who commented on the post
			userQuery = `
				SELECT u.first_name, u.last_name, u.avatar_url 
				FROM users u 
				JOIN comments c ON u.id = c.user_id 
				WHERE c.post_id = ? AND c.deleted_at IS NULL 
				ORDER BY c.created_at DESC LIMIT 1
			`
			queryParam = referenceID
		default:
			userQuery = "SELECT first_name, last_name, avatar_url FROM users WHERE id = ?"
			queryParam = referenceID
		}

		err = nm.DB.QueryRowContext(ctx, userQuery, queryParam).Scan(&firstName, &lastName, &avatarURL)
		if err != nil {
			log.Printf("ERROR: Failed to get user info for notification %s: %v", id, err)
			firstName, lastName, avatarURL = "Unknown", "User", ""
		}

		// Format user name
		userName := firstName + " " + lastName

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
