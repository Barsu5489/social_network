package repository

import (
	"database/sql"
	"fmt"
	"time"

	"social-nework/pkg/models"

	"github.com/google/uuid"
)

type GroupRepository struct {
	DB *sql.DB
}

// IsUserMember checks if a user is a member of a group
func (r *GroupRepository) IsUserMember(groupID, userID string) (bool, error) {
	var exists bool
	err := r.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM group_members 
			WHERE group_id = ? AND user_id = ? AND deleted_at IS NULL
		)`, groupID, userID).Scan(&exists)
	return exists, err
}

// AddMember adds a user to a group with the specified role
func (r *GroupRepository) AddMember(groupID, userID, role string) error {
	// Check if user is already a member (including soft-deleted members)
	var existingID string
	var deletedAt sql.NullTime

	err := r.DB.QueryRow(`
		SELECT id, deleted_at FROM group_members 
		WHERE group_id = ? AND user_id = ?`,
		groupID, userID).Scan(&existingID, &deletedAt)

	if err == nil {
		// User exists, check if they're soft-deleted
		if deletedAt.Valid {
			// Restore the membership
			_, err = r.DB.Exec(`
				UPDATE group_members 
				SET deleted_at = NULL, role = ?, joined_at = ?
				WHERE id = ?`,
				role, time.Now(), existingID)
			return err
		}
		// User is already an active member
		return fmt.Errorf("user is already a member of this group")
	}

	if err != sql.ErrNoRows {
		return err
	}

	// Create new membership
	_, err = r.DB.Exec(`
		INSERT INTO group_members (id, group_id, user_id, role, joined_at)
		VALUES (?, ?, ?, ?, ?)`,
		uuid.New().String(), groupID, userID, role, time.Now())
	return err
}

// RemoveMember removes a user from a group (soft delete)
func (r *GroupRepository) RemoveMember(groupID, userID string) error {
	_, err := r.DB.Exec(`
		UPDATE group_members 
		SET deleted_at = ? 
		WHERE group_id = ? AND user_id = ? AND deleted_at IS NULL`,
		time.Now(), groupID, userID)
	return err
}

// GetGroupMembers retrieves all active members of a group
func (r *GroupRepository) GetGroupMembers(groupID string) ([]models.GroupMember, error) {
	rows, err := r.DB.Query(`
		SELECT gm.id, gm.group_id, gm.user_id, gm.role, gm.joined_at,
			   u.first_name, u.last_name, u.email, u.avatar_url
		FROM group_members gm
		JOIN users u ON gm.user_id = u.id
		WHERE gm.group_id = ? AND gm.deleted_at IS NULL
		ORDER BY gm.joined_at ASC`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.GroupMember
	for rows.Next() {
		var member models.GroupMember
		var user models.User
		err := rows.Scan(
			&member.ID, &member.GroupID, &member.UserID, &member.Role, &member.JoinedAt,
			&user.FirstName, &user.LastName, &user.Email, &user.AvatarURL,
		)
		if err != nil {
			return nil, err
		}
		member.User = user
		members = append(members, member)
	}
	return members, nil
}

// GetUserGroups retrieves all groups a user is a member of
func (r *GroupRepository) GetUserGroups(userID string) ([]models.Group, error) {
	rows, err := r.DB.Query(`
		SELECT g.id, g.name, g.description, g.creator_id, g.is_private, g.created_at, g.updated_at
		FROM groups g
		JOIN group_members gm ON g.id = gm.group_id
		WHERE gm.user_id = ? AND gm.deleted_at IS NULL AND g.deleted_at IS NULL
		ORDER BY g.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []models.Group
	for rows.Next() {
		var group models.Group
		err := rows.Scan(
			&group.ID, &group.Name, &group.Description, &group.CreatorID,
			&group.IsPrivate, &group.CreatedAt, &group.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}


// CreateGroupWithChat creates a group and its associated chat in a single transaction
func (r *GroupRepository) CreateGroupWithChat(group *models.Group) error {
	
	// Check if repository is properly initialized
	if r == nil {
		return fmt.Errorf("GroupRepository is nil")
	}
	if r.DB == nil {
		return fmt.Errorf("Database connection is nil")
	}
	if group == nil {
		return fmt.Errorf("Group model is nil")
	}

	// Start transaction
	tx, err := r.DB.Begin()
	if err != nil {
		fmt.Println("error starting transaction:", err)
		return err
	}

	// Use a flag to track if we should rollback
	committed := false
	defer func() {
		if !committed {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("error rolling back transaction: %v\n", rbErr)
			}
		}
	}()

	// Create the group
	_, err = tx.Exec(`
		INSERT INTO groups (id, name, description, creator_id, is_private, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		group.ID, group.Name, group.Description, group.CreatorID, group.IsPrivate,
		group.CreatedAt, group.UpdatedAt)
	if err != nil {
		fmt.Printf("error creating group: %v\n", err)
		return err
	}

	// Add creator as admin member
	_, err = tx.Exec(`
		INSERT INTO group_members (id, group_id, user_id, role, joined_at)
		VALUES (?, ?, ?, 'admin', ?)`,
		uuid.New().String(), group.ID, group.CreatorID, time.Now())
	if err != nil {
		fmt.Printf("error adding creator as admin: %v\n", err)
		return err
	}

	// Create group chat using the same transaction
	chatID := uuid.New().String()
	_, err = tx.Exec(`
		INSERT INTO chats (id, type, created_at)
		VALUES (?, 'group', ?)`,
		chatID, time.Now())
	if err != nil {
		fmt.Printf("error creating group chat: %v\n", err)
		return err
	}

	// Add creator as chat participant
	_, err = tx.Exec(`
		INSERT INTO chat_participants (id, chat_id, user_id, joined_at)
		VALUES (?, ?, ?, ?)`,
		uuid.New().String(), chatID, group.CreatorID, time.Now())
	if err != nil {
		fmt.Printf("error adding creator as chat participant: %v\n", err)
		return err
	}

	// Link the group to its chat
	_, err = tx.Exec(`
		INSERT INTO group_chats (group_id, chat_id, created_at)
		VALUES (?, ?, ?)`,
		group.ID, chatID, time.Now())
	if err != nil {
		fmt.Printf("error linking group to chat: %v\n", err)
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		fmt.Printf("error committing transaction: %v\n", err)
		return err
	}

	committed = true
	return nil
}
