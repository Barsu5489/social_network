package models

import (
	"context"
	"database/sql"
)

// GetCommentLikes retrieves all likes for a specific comment
func GetCommentLikes(db *sql.DB, ctx context.Context, commentID string) ([]Like, error) {
	stmt := `
		SELECT id, user_id, likeable_id, created_at 
		FROM likes 
		WHERE likeable_type = 'comment' AND likeable_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := db.QueryContext(ctx, stmt, commentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var likes []Like
	for rows.Next() {
		var like Like
		like.LikeableType = "comment"

		err := rows.Scan(&like.ID, &like.UserID, &like.LikeableID, &like.CreatedAt)
		if err != nil {
			return nil, err
		}

		likes = append(likes, like)
	}

	if likes == nil {
		return []Like{}, nil
	}

	return likes, nil
}

// HasUserLikedComment checks if a user has liked a specific comment
func HasUserLikedComment(db *sql.DB, ctx context.Context, userID, commentID string) (bool, error) {
	stmt := `
		SELECT 1 FROM likes 
		WHERE user_id = ? AND likeable_type = 'comment' AND likeable_id = ? AND deleted_at IS NULL
		LIMIT 1
	`

	var exists int
	err := db.QueryRowContext(ctx, stmt, userID, commentID).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
