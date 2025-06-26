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

// GetChatParticipantsWithDetails retrieves participants with their user details
func (r *ChatRepository) GetChatParticipantsWithDetails(chatID string) ([]models.User, error) {
	rows, err := r.DB.Query(`
		SELECT u.id, u.first_name, u.last_name, u.email, u.avatar_url
		FROM chat_participants cp
		JOIN users u ON cp.user_id = u.id
		WHERE cp.chat_id = ? AND cp.deleted_at IS NULL
		ORDER BY u.first_name, u.last_name`, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.FirstName, &user.LastName,
			&user.Email, &user.AvatarURL,
		)
		if err != nil {
			return nil, err
		}
		participants = append(participants, user)
	}
	return participants, nil
}

// GetUserChats retrieves all chats for a user
func (r *ChatRepository) GetUserChats(userID string) ([]models.Chat, error) {
	rows, err := r.DB.Query(`
		SELECT c.id, c.type, c.created_at
		FROM chats c
		JOIN chat_participants cp ON c.id = cp.chat_id
		WHERE cp.user_id = ? AND c.deleted_at IS NULL AND cp.deleted_at IS NULL
		ORDER BY c.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []models.Chat
	for rows.Next() {
		var chat models.Chat
		if err := rows.Scan(&chat.ID, &chat.Type, &chat.CreatedAt); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}
	return chats, nil
}

// GetUserChatIDs retrieves all chat IDs for a user
func (r *ChatRepository) GetUserChatIDs(userID string) ([]string, error) {
	rows, err := r.DB.Query(`
		SELECT chat_id FROM chat_participants
		WHERE user_id = ? AND deleted_at IS NULL`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chatIDs []string
	for rows.Next() {
		var chatID string
		if err := rows.Scan(&chatID); err != nil {
			return nil, err
		}
		chatIDs = append(chatIDs, chatID)
	}
	return chatIDs, nil
}

// IsUserInChat checks if a user is a participant in a chat
func (r *ChatRepository) IsUserInChat(chatID, userID string) (bool, error) {
	var exists bool
	err := r.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM chat_participants
			WHERE chat_id = ? AND user_id = ? AND deleted_at IS NULL
		)`, chatID, userID).Scan(&exists)
	return exists, err
}

// GetChatInfo retrieves basic chat information
func (r *ChatRepository) GetChatInfo(chatID string) (*models.Chat, error) {
	var chat models.Chat
	err := r.DB.QueryRow(`
		SELECT id, type, created_at
		FROM chats
		WHERE id = ? AND deleted_at IS NULL`, chatID).Scan(
		&chat.ID, &chat.Type, &chat.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &chat, nil
}
