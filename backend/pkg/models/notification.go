package models

import (
	"database/sql"
	"context"
)

// NotificationModel provides methods for interacting with the notifications table.
type NotificationModel struct {
	DB *sql.DB
}

// Insert inserts a new notification into the database.
func (m *NotificationModel) Insert(ctx context.Context, notification Notification) (int, error) {
	stmt, err := m.DB.PrepareContext(ctx, "INSERT INTO notifications(user_id, type, source_id, content, read_status) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, notification.UserID, notification.Type, notification.SourceID, notification.Content, notification.ReadStatus)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// GetByUserID retrieves all notifications for a given user.
func (m *NotificationModel) GetByUserID(ctx context.Context, userID int) ([]Notification, error) {
	rows, err := m.DB.QueryContext(ctx, "SELECT id, user_id, type, source_id, content, read_status, created_at FROM notifications WHERE user_id = ? ORDER BY created_at DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var n Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.SourceID, &n.Content, &n.ReadStatus, &n.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

// MarkAsRead marks a notification as read.
func (m *NotificationModel) MarkAsRead(ctx context.Context, notificationID int) error {
	stmt, err := m.DB.PrepareContext(ctx, "UPDATE notifications SET read_status = TRUE WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, notificationID)
	return err
}
