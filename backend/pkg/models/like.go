package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

func CreateLike(db *sql.DB, ctx context.Context, userID, likeableType, likeableID string) (*Like, error) {
	if likeableType != "post" && likeableType != "comment" {
		return nil, errors.New("invalid likeable type")
	}
	// Check if the user already liked this content
	var existingID string
	var deletedAt *int64
	checkStmt := `
		SELECT id, deleted_at FROM likes 
		WHERE user_id = ? AND likeable_type = ? AND likeable_id = ?
	`
	err := db.QueryRowContext(ctx, checkStmt, userID, likeableType, likeableID).Scan(&existingID, &deletedAt)

	now := time.Now().Unix()
	// If like exists but was deleted, reactivate it
	if err == nil && deletedAt != nil {
		// Reactivate the like by setting deleted_at to NULL
		updateStmt := `UPDATE likes SET deleted_at = NULL WHERE id = ?`
		_, err := db.ExecContext(ctx, updateStmt, existingID)
		if err != nil {
			return nil, err
		}
		return &Like{
			ID:           existingID,
			UserID:       userID,
			LikeableType: likeableType,
			LikeableID:   likeableID,
			CreatedAt:    now,
			DeletedAt:    nil,
		}, nil
	} else if err == nil {
		// Like already exists and is active
		return nil, errors.New("user already liked this content")
	} else if err != sql.ErrNoRows {
		// Some other database error
		return nil, err
	}
	// Create new like
	id := uuid.New().String()
	insertStmt := `
		INSERT INTO likes (id, user_id, likeable_type, likeable_id, created_at) 
		VALUES (?, ?, ?, ?, ?)
	`

	_, err = db.ExecContext(ctx, insertStmt, id, userID, likeableType, likeableID, now)
	if err != nil {
		return nil, err
	}
}
