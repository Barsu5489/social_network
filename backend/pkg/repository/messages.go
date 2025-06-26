package repository

import (
	"database/sql"
	"fmt"
	"time"

	"social-nework/pkg/models"

	"github.com/google/uuid"
)

type MessageRepository struct {
	DB *sql.DB
}

// SaveMessage saves a message to the database
func (r *MessageRepository) SaveMessage(msg *models.Message) error {
	_, err := r.DB.Exec(`
		INSERT INTO messages (id, chat_id, sender_id, content, sent_at)
		VALUES (?, ?, ?, ?, ?)`,
		msg.ID, msg.ChatID, msg.SenderID, msg.Content, msg.SentAt)
	return err
}

// GetChatMessages retrieves messages for a chat with pagination
func (r *MessageRepository) GetChatMessages(chatID string, before time.Time, limit int) ([]models.Message, error) {
	query := `
		SELECT m.id, m.chat_id, m.sender_id, m.content, m.sent_at, m.read_at,
			   u.first_name, u.last_name, u.avatar_url
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.chat_id = ? AND m.deleted_at IS NULL`

	args := []interface{}{chatID}

	if !before.IsZero() {
		query += " AND m.sent_at < ?"
		args = append(args, before)
	}

	query += " ORDER BY m.sent_at DESC LIMIT ?"
	args = append(args, limit)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		var sender models.User
		var readAt sql.NullTime

		err := rows.Scan(
			&msg.ID, &msg.ChatID, &msg.SenderID, &msg.Content, &msg.SentAt, &readAt,
			&sender.FirstName, &sender.LastName, &sender.AvatarURL,
		)
		if err != nil {
			return nil, err
		}

		msg.Sender = sender
		if readAt.Valid {
			msg.ReadAt = readAt.Time
		}

		messages = append(messages, msg)
	}

	return messages, nil
}