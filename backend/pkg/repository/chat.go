package repository

import (
	"database/sql"
	"fmt"
	"time"

	"social-nework/pkg/models"

	"github.com/google/uuid"
)

type ChatRepository struct {
	DB *sql.DB
}

// CreateChat creates a new chat (direct or group)
func (r *ChatRepository) CreateChat(chatType string, creatorID string) (*models.Chat, error) {
	chatID := uuid.New().String()
	now := time.Now()

	_, err := r.DB.Exec(`
		INSERT INTO chats (id, type, created_at)
		VALUES (?, ?, ?)`,
		chatID, chatType, now)
	if err != nil {
		return nil, err
	}

	// Add creator as participant
	_, err = r.DB.Exec(`
		INSERT INTO chat_participants (id, chat_id, user_id, joined_at)
		VALUES (?, ?, ?, ?)`,
		uuid.New().String(), chatID, creatorID, now)
	if err != nil {
		return nil, err
	}

	return &models.Chat{
		ID:        chatID,
		Type:      chatType,
		CreatedAt: now,
	}, nil
}

// AddParticipant adds a user to a chat
func (r *ChatRepository) AddParticipant(chatID, userID string) error {
	// Check if user is already a participant (including soft-deleted)
	var existingID string
	var deletedAt sql.NullTime

	err := r.DB.QueryRow(`
		SELECT id, deleted_at FROM chat_participants 
		WHERE chat_id = ? AND user_id = ?`,
		chatID, userID).Scan(&existingID, &deletedAt)

	if err == nil {
		// User exists, check if they're soft-deleted
		if deletedAt.Valid {
			// Restore the participation
			_, err = r.DB.Exec(`
				UPDATE chat_participants 
				SET deleted_at = NULL, joined_at = ?
				WHERE id = ?`,
				time.Now(), existingID)
			return err
		}
		// User is already an active participant
		return fmt.Errorf("user is already a participant in this chat")
	}

	if err != sql.ErrNoRows {
		return err
	}

	// Create new participation
	_, err = r.DB.Exec(`
		INSERT INTO chat_participants (id, chat_id, user_id, joined_at)
		VALUES (?, ?, ?, ?)`,
		uuid.New().String(), chatID, userID, time.Now())
	return err
}

// RemoveParticipant removes a user from a chat (soft delete)
func (r *ChatRepository) RemoveParticipant(chatID, userID string) error {
	_, err := r.DB.Exec(`
		UPDATE chat_participants 
		SET deleted_at = ? 
		WHERE chat_id = ? AND user_id = ? AND deleted_at IS NULL`,
		time.Now(), chatID, userID)
	return err
}

// GetChatParticipants retrieves all active participants of a chat
func (r *ChatRepository) GetChatParticipants(chatID string) ([]string, error) {
	rows, err := r.DB.Query(`
		SELECT user_id FROM chat_participants
		WHERE chat_id = ? AND deleted_at IS NULL`, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		participants = append(participants, userID)
	}
	return participants, nil
}
