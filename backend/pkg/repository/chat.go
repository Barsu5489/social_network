package repository

import (
	"database/sql"
	"fmt"
	"log"
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
	now := time.Now().Unix()
	
	log.Printf("CreateChat: Inserting chat %s with Unix timestamp: %d", chatID, now)

	_, err := r.DB.Exec(`
		INSERT INTO chats (id, type, created_at)
		VALUES (?, ?, ?)`,
		chatID, chatType, now)
	if err != nil {
		log.Printf("CreateChat: Insert error: %v", err)
		return nil, err
	}

	// Verify what was actually inserted
	var storedCreatedAt interface{}
	var storedType string
	err = r.DB.QueryRow(`
		SELECT created_at, typeof(created_at) 
		FROM chats WHERE id = ?`, chatID).Scan(&storedCreatedAt, &storedType)
	if err != nil {
		log.Printf("CreateChat: Verification query failed: %v", err)
	} else {
		log.Printf("CreateChat: Stored created_at: %v (type: %s)", storedCreatedAt, storedType)
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
		CreatedAt: time.Unix(now, 0),
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
	// First, let's see what's actually in the database
	debugRows, err := r.DB.Query(`
		SELECT c.id, c.type, c.created_at, typeof(c.created_at) as type_info
		FROM chats c
		JOIN chat_participants cp ON c.id = cp.chat_id
		WHERE cp.user_id = ? AND c.deleted_at IS NULL AND cp.deleted_at IS NULL
		ORDER BY c.created_at DESC LIMIT 3`, userID)
	if err != nil {
		log.Printf("GetUserChats: Debug query failed: %v", err)
	} else {
		defer debugRows.Close()
		log.Printf("GetUserChats: Debug - checking what's in database:")
		for debugRows.Next() {
			var id, chatType, createdAt, typeInfo string
			if err := debugRows.Scan(&id, &chatType, &createdAt, &typeInfo); err != nil {
				log.Printf("GetUserChats: Debug scan error: %v", err)
			} else {
				log.Printf("GetUserChats: Debug - Chat %s: type=%s, created_at=%s, sql_type=%s", id, chatType, createdAt, typeInfo)
			}
		}
	}

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
		var createdAtUnix int64
		if err := rows.Scan(&chat.ID, &chat.Type, &createdAtUnix); err != nil {
			log.Printf("GetUserChats: Scan error for chat %s: %v", chat.ID, err)
			return nil, err
		}
		chat.CreatedAt = time.Unix(createdAtUnix, 0)
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
	var createdAtUnix int64
	err := r.DB.QueryRow(`
		SELECT id, type, created_at
		FROM chats
		WHERE id = ? AND deleted_at IS NULL`, chatID).Scan(
		&chat.ID, &chat.Type, &createdAtUnix)
	if err != nil {
		return nil, err
	}
	chat.CreatedAt = time.Unix(createdAtUnix, 0)
	return &chat, nil
}


// CreateGroupChat creates a chat specifically for a group
func (r *ChatRepository) CreateGroupChat(groupID, creatorID string) (string, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	chatID := uuid.New().String()
	now := time.Now().Unix()

	// Create the chat
	_, err = tx.Exec(`
		INSERT INTO chats (id, type, created_at)
		VALUES (?, 'group', ?)`,
		chatID, now)
	if err != nil {
		return "", err
	}

	// Add creator as participant
	_, err = tx.Exec(`
		INSERT INTO chat_participants (id, chat_id, user_id, joined_at)
		VALUES (?, ?, ?, ?)`,
		uuid.New().String(), chatID, creatorID, now)
	if err != nil {
		return "", err
	}

	// Link group to chat
	_, err = tx.Exec(`
		INSERT INTO group_chats (group_id, chat_id, created_at)
		VALUES (?, ?, ?)`,
		groupID, chatID, now)
	if err != nil {
		return "", err
	}

	err = tx.Commit()
	return chatID, err
}

// AddMemberToGroupChat adds a member to group chat when they join the group
func (r *ChatRepository) AddMemberToGroupChat(groupID, userID string) error {
	// Get the group's chat ID
	var chatID string
	err := r.DB.QueryRow(`
		SELECT chat_id FROM group_chats 
		WHERE group_id = ?`, groupID).Scan(&chatID)
	if err != nil {
		return err
	}

	// Add user to chat
	return r.AddParticipant(chatID, userID)
}

// RemoveMemberFromGroupChat removes a member from group chat when they leave the group
func (r *ChatRepository) RemoveMemberFromGroupChat(groupID, userID string) error {
	// Get the group's chat ID
	var chatID string
	err := r.DB.QueryRow(`
		SELECT chat_id FROM group_chats 
		WHERE group_id = ?`, groupID).Scan(&chatID)
	if err != nil {
		return err
	}

	// Remove user from chat
	return r.RemoveParticipant(chatID, userID)
}

// GetGroupChatMembers gets all group members for a group chat
func (r *ChatRepository) GetGroupChatMembers(groupID string) ([]string, error) {
	rows, err := r.DB.Query(`
		SELECT gm.user_id 
		FROM group_members gm
		WHERE gm.group_id = ? AND gm.deleted_at IS NULL`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		members = append(members, userID)
	}
	return members, nil
}

// GetDirectChatBetweenUsers finds existing direct chat between two users
func (r *ChatRepository) GetDirectChatBetweenUsers(user1, user2 string) (string, error) {
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
	return chatID, err
}

// SoftDeleteChat marks a chat as deleted
func (r *ChatRepository) SoftDeleteChat(chatID string) error {
	_, err := r.DB.Exec(`
		UPDATE chats 
		SET deleted_at = ? 
		WHERE id = ? AND deleted_at IS NULL`,
		time.Now(), chatID)
	return err
}



// GetChatType returns the type of a chat
func (r *ChatRepository) GetChatType(chatID string) (string, error) {
	var chatType string
	err := r.DB.QueryRow(`
		SELECT type FROM chats 
		WHERE id = ? AND deleted_at IS NULL`, chatID).Scan(&chatType)
	return chatType, err
}

// GetChatCount returns the number of chats a user is in
func (r *ChatRepository) GetChatCount(userID string) (int, error) {
	var count int
	err := r.DB.QueryRow(`
		SELECT COUNT(DISTINCT cp.chat_id)
		FROM chat_participants cp
		JOIN chats c ON cp.chat_id = c.id
		WHERE cp.user_id = ? AND cp.deleted_at IS NULL AND c.deleted_at IS NULL`,
		userID).Scan(&count)
	return count, err
}
