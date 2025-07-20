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
	query := `
		SELECT n.id, n.type, n.reference_id, n.created_at,
			   u.first_name, u.last_name, u.avatar_url
		FROM notifications n
		LEFT JOIN users u ON (
			CASE 
				WHEN n.type = 'new_follower' THEN n.reference_id = u.id
				WHEN n.type = 'new_like' THEN (
					SELECT l.user_id FROM likes l WHERE l.id = n.reference_id
				) = u.id
				WHEN n.type = 'comment_on_post' THEN (
					SELECT c.user_id FROM comments c WHERE c.id = n.reference_id
				) = u.id
				WHEN n.type = 'new_message' THEN (
					SELECT m.sender_id FROM messages m WHERE m.id = n.reference_id
				) = u.id
				ELSE n.reference_id = u.id
			END
		)
		WHERE n.user_id = ? AND n.is_read = false AND n.deleted_at IS NULL
		ORDER BY n.created_at DESC
	`

	rows, err := nm.DB.QueryContext(ctx, query, userID)
	if err != nil {
		log.Printf("ERROR: Failed to query notifications: %v", err)
		return []map[string]interface{}{}, nil // Return empty array instead of nil
	}
	defer rows.Close()

	var notifications []map[string]interface{}

	for rows.Next() {
		var id, notifType, referenceID string
		var createdAt time.Time
		var firstName, lastName, avatarURL sql.NullString

		err := rows.Scan(&id, &notifType, &referenceID, &createdAt, &firstName, &lastName, &avatarURL)
		if err != nil {
			log.Printf("ERROR: Failed to scan notification: %v", err)
			continue
		}

		// Format user name
		userName := "Unknown User"
		if firstName.Valid && lastName.Valid {
			userName = firstName.String + " " + lastName.String
		} else if firstName.Valid {
			userName = firstName.String
		}

		// Generate message and href based on type
		message, href := nm.formatNotificationMessage(notifType, referenceID)

		notification := map[string]interface{}{
			"id":   id,
			"type": notifType,
			"user": map[string]interface{}{
				"name":   userName,
				"avatar": avatarURL.String,
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
	case "comment_like":
		return "liked your comment", "/posts/" + referenceID
	case "comment_on_post":
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
