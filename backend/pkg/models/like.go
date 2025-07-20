package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

func CreateLike(db *sql.DB, ctx context.Context, notificationModel *NotificationModel, userID, likeableType, likeableID string) (*Like, error) {
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
	// Create notification for the content owner if it's a post like
	if likeableType == "post" {
		// Get post owner
		var postOwnerID string
		postOwnerStmt := `SELECT user_id FROM posts WHERE id = ?`
		err = db.QueryRowContext(ctx, postOwnerStmt, likeableID).Scan(&postOwnerID)
		if err == nil && postOwnerID != userID {
			log.Printf("DEBUG: Creating like notification - PostOwner: %s, Liker: %s, PostID: %s", postOwnerID, userID, likeableID)
			// Don't notify for self-likes
			notification := Notification{
				ID:          uuid.New().String(),
				UserID:      postOwnerID,
				Type:        "new_like",
				ReferenceID: likeableID, // Use post ID instead of like ID
				IsRead:      false,
				CreatedAt:   time.Now(),
			}
			_, err = notificationModel.Insert(ctx, notification)
			if err != nil {
				log.Printf("ERROR: Failed to create notification for like: %v", err)
			} else {
				log.Printf("SUCCESS: Like notification created for post owner: %s", postOwnerID)
			}
		} else if err != nil {
			log.Printf("ERROR: Failed to get post owner for like notification: %v", err)
		} else {
			log.Printf("DEBUG: Not creating like notification - self-like detected")
		}
	}

	// Create notification for comment likes
	if likeableType == "comment" {
		// Get comment owner and post ID
		var commentOwnerID, postID string
		commentOwnerStmt := `SELECT user_id, post_id FROM comments WHERE id = ?`
		err = db.QueryRowContext(ctx, commentOwnerStmt, likeableID).Scan(&commentOwnerID, &postID)
		if err == nil && commentOwnerID != userID {
			// Don't notify for self-likes
			notification := Notification{
				ID:          uuid.New().String(),
				UserID:      commentOwnerID,
				Type:        "new_like", // Use same type as post likes
				ReferenceID: postID,     // Use post ID for navigation
				IsRead:      false,
				CreatedAt:   time.Now(),
			}
			_, err = notificationModel.Insert(ctx, notification)
			if err != nil {
				fmt.Printf("Failed to create notification for comment like: %v\n", err)
			} else {
				log.Printf("SUCCESS: Comment like notification created for comment owner: %s", commentOwnerID)
			}
		}
	}

	return &Like{
		ID:           id,
		UserID:       userID,
		LikeableType: likeableType,
		LikeableID:   likeableID,
		CreatedAt:    now,
		DeletedAt:    nil,
	}, nil
}

// UnlikeContent removes a like from a post or comment
func UnlikeContent(db *sql.DB, ctx context.Context, userID, likeableType, likeableID string) error {
	if likeableType != "post" && likeableType != "comment" {
		return errors.New("invalid likeable type")
	}

	now := time.Now().Unix()

	// Soft delete the like by setting deleted_at
	stmt := `
		UPDATE likes 
		SET deleted_at = ? 
		WHERE user_id = ? AND likeable_type = ? AND likeable_id = ? AND deleted_at IS NULL
	`

	result, err := db.ExecContext(ctx, stmt, now, userID, likeableType, likeableID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("like not found or already removed")
	}

	return nil
}

// GetPostLikes retrieves all likes for a specific post
func GetPostLikes(db *sql.DB, ctx context.Context, postID string) ([]Like, error) {
	stmt := `
		  SELECT id, user_id, likeable_id, created_at 
        FROM likes 
        WHERE likeable_type = 'post' AND likeable_id = ? AND deleted_at IS NULL
        ORDER BY created_at DESC
	`

	rows, err := db.QueryContext(ctx, stmt, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var likes []Like
	for rows.Next() {
		var like Like
		like.LikeableType = "post"

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

// HasUserLikedPost checks if a user has liked a specific post
func HasUserLikedPost(db *sql.DB, ctx context.Context, userID, postID string) (bool, error) {
	stmt := `
		SELECT 1 FROM likes 
		WHERE user_id = ? AND likeable_type = 'post' AND likeable_id = ? AND deleted_at IS NULL
		LIMIT 1
	`

	var exists int
	err := db.QueryRowContext(ctx, stmt, userID, postID).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

// GetLikeCount returns the count of likes for a post or comment
func GetLikeCount(db *sql.DB, ctx context.Context, likeableType, likeableID string) (int, error) {
	if likeableType != "post" && likeableType != "comment" {
		return 0, errors.New("invalid likeable type")
	}

	stmt := `
		SELECT COUNT(*) FROM likes 
		WHERE likeable_type = ? AND likeable_id = ? AND deleted_at IS NULL
	`

	var count int
	err := db.QueryRowContext(ctx, stmt, likeableType, likeableID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
