package models

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

// NotificationModel provides methods for interacting with the notifications table.
type NotificationModel struct {
	DB *sql.DB
}

// Insert inserts a new notification into the database.
func (m *NotificationModel) Insert(ctx context.Context, notification Notification) (string, error) {
	notification.ID = uuid.New().String()
	stmt, err := m.DB.PrepareContext(ctx, "INSERT INTO notifications(id, user_id, type, reference_id, is_read, created_at) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, notification.ID, notification.UserID, notification.Type, notification.ReferenceID, notification.IsRead, notification.CreatedAt)
	if err != nil {
		return "", err
	}

	return notification.ID, nil
}

// GetByUserID retrieves all notifications for a given user.
func (m *NotificationModel) GetByUserID(ctx context.Context, userID string) ([]Notification, error) {
	rows, err := m.DB.QueryContext(ctx, "SELECT id, user_id, type, reference_id, is_read, created_at FROM notifications WHERE user_id = ? ORDER BY created_at DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var n Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.ReferenceID, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

// MarkAsRead marks a notification as read.
func (m *NotificationModel) MarkAsRead(ctx context.Context, notificationID string) error {
	stmt, err := m.DB.PrepareContext(ctx, "UPDATE notifications SET is_read = 1 WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, notificationID)
	return err
}
