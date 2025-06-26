package repository

import (
	"database/sql"
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

// CreateDirectChat creates or returns existing direct chat between two users
func (r *MessageRepository) CreateDirectChat(user1, user2 string) (string, error) {
	// Check if chat already exists
	var chatID string
	err := r.DB.QueryRow(`
		SELECT cp1.chat_id
		FROM chat_participants cp1
		JOIN chat_participants cp2 ON cp1.chat_id = cp2.chat_id
		JOIN chats c ON cp1.chat_id = c.id
		WHERE cp1.user_id = ? AND cp2.user_id = ?
		AND c.type = 'direct' AND c.deleted_at IS NULL
		AND cp1.deleted_at IS NULL AND cp2.deleted_at IS NULL`,
		user1, user2).Scan(&chatID)

	if err == nil {
		return chatID, nil
	}

	if err != sql.ErrNoRows {
		return "", err
	}

	// Create new chat
	tx, err := r.DB.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	chatID = uuid.New().String()

	_, err = tx.Exec(`
		INSERT INTO chats (id, type, created_at)
		VALUES (?, 'direct', ?)`,
		chatID, time.Now())
	if err != nil {
		return "", err
	}

	_, err = tx.Exec(`
		INSERT INTO chat_participants (id, chat_id, user_id, joined_at)
		VALUES (?, ?, ?, ?)`,
		uuid.New().String(), chatID, user1, time.Now())
	if err != nil {
		return "", err
	}

	_, err = tx.Exec(`
		INSERT INTO chat_participants (id, chat_id, user_id, joined_at)
		VALUES (?, ?, ?, ?)`,
		uuid.New().String(), chatID, user2, time.Now())
	if err != nil {
		return "", err
	}

	return chatID, tx.Commit()
}

// MarkMessageAsRead marks a message as read by a user
func (r *MessageRepository) MarkMessageAsRead(messageID, userID string) error {
	_, err := r.DB.Exec(`
		UPDATE messages 
		SET read_at = ? 
		WHERE id = ? AND read_at IS NULL`,
		time.Now(), messageID)
	return err
}

// GetUnreadMessageCount returns the count of unread messages for a user in a chat
func (r *MessageRepository) GetUnreadMessageCount(chatID, userID string) (int, error) {
	var count int
	err := r.DB.QueryRow(`
		SELECT COUNT(*) 
		FROM messages m
		WHERE m.chat_id = ? AND m.sender_id != ? 
		AND m.read_at IS NULL AND m.deleted_at IS NULL`,
		chatID, userID).Scan(&count)
	return count, err
}
