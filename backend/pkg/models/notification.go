package models

import (
	"context"
	"database/sql"
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

func (m *NotificationModel) GetByUserID(ctx context.Context, userID string) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			n.id, 
			n.type, 
			n.reference_id,
			n.created_at,
			u.first_name,
			u.last_name,
			u.avatar_url
		FROM notifications n
		LEFT JOIN users u ON u.id = (
			CASE 
				WHEN n.type = 'follow_request' THEN n.reference_id
				WHEN n.type = 'group_invite' THEN n.reference_id
				ELSE n.reference_id
			END
		)
		WHERE n.user_id = ? AND (n.deleted_at IS NULL OR n.deleted_at = 0)
		ORDER BY n.created_at DESC
		LIMIT 50
	`
	
	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var notifications []map[string]interface{}
	for rows.Next() {
		var id, notifType, referenceID string
		var createdAtRaw interface{}
		var firstName, lastName, avatarURL sql.NullString
		
		err := rows.Scan(
			&id,
			&notifType,
			&referenceID,
			&createdAtRaw,
			&firstName,
			&lastName,
			&avatarURL,
		)
		if err != nil {
			return nil, err
		}
		
		// Parse timestamp
		var createdAt time.Time
		switch v := createdAtRaw.(type) {
		case int64:
			createdAt = time.Unix(v, 0)
		case string:
			if t, err := time.Parse(time.RFC3339Nano, v); err == nil {
				createdAt = t
			} else {
				createdAt = time.Now()
			}
		default:
			createdAt = time.Now()
		}
		
		// Format user name
		userName := "Unknown User"
		if firstName.Valid && lastName.Valid {
			userName = firstName.String + " " + lastName.String
		} else if firstName.Valid {
			userName = firstName.String
		}
		
		// Generate message based on type
		message := ""
		href := "#"
		switch notifType {
		case "follow_request":
			message = "sent you a follow request"
			href = "/profile/" + referenceID
		case "group_invite":
			message = "invited you to join a group"
			href = "/groups/" + referenceID
		default:
			message = "sent you a notification"
		}
		
		notification := map[string]interface{}{
			"id":   id,
			"type": notifType,
			"user": map[string]interface{}{
				"name":   userName,
				"avatar": avatarURL.String,
			},
			"message": message,
			"time":    createdAt.Format("2 Jan 2006"),
			"href":    href,
		}
		
		notifications = append(notifications, notification)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return notifications, nil
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
