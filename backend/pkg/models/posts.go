package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func CreatePost(db *sql.DB, ctx context.Context, userID, content, privacy string, groupID *string, allowedUserIDs []string, imageURL *string) (string, error) {
	if content == "" && imageURL == nil {
		return "", errors.New("content or image is required")
	}

	// Validate privacy setting
	isValidPrivacy := false
	switch privacy {
	case "public", "almost_private", "private":
		isValidPrivacy = true
	}
	if !isValidPrivacy {
		return "", errors.New("invalid privacy setting")
	}

	// If privacy is private, the allowed users list must not be empty
	if privacy == "private" && len(allowedUserIDs) == 0 {
		return "", errors.New("private posts must specify at least one allowed user")
	}

	// --- Start Transaction ---
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Defer a rollback in case anything fails. It will be ignored if we commit.
	defer tx.Rollback()

	postID := uuid.New().String()
	now := time.Now().Unix()

	// --- 1. Insert into posts table ---
	postStm := `INSERT INTO posts (id, user_id, group_id, content, privacy, image_url, created_at, updated_at)
                VALUES(?,?,?,?,?,?,?,?)`
	_, err = tx.ExecContext(ctx, postStm, postID, userID, groupID, content, privacy, imageURL, now, now)
	if err != nil {
		return "", fmt.Errorf("failed to insert post: %w", err)
	}
	
	// --- 2. If the post is private, insert into post_allowed_users ---
	if privacy == "private" {
		allowedUsersStm, err := tx.PrepareContext(ctx, `INSERT INTO post_allowed_users (post_id, user_id) VALUES (?, ?)`)
		if err != nil {
			return "", fmt.Errorf("failed to prepare statement for allowed users: %w", err)
		}
		defer allowedUsersStm.Close()

		for _, allowedID := range allowedUserIDs {
			// It's good practice to ensure the creator doesn't need to add themselves, 
			// but for explicit control, we'll allow whatever is passed.
			if _, err := allowedUsersStm.ExecContext(ctx, postID, allowedID); err != nil {
				return "", fmt.Errorf("failed to insert allowed user %s: %w", allowedID, err)
			}
		}
	}

	// --- Commit Transaction ---
	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return postID, nil
}

// GetFeed retrieves posts from users the given user follows
func GetFollowingPosts(db *sql.DB, userID string) ([]Post, error) {
	stm := `
		SELECT p.id, p.user_id, p.content, p.privacy, p.created_at,
		       (SELECT COUNT(*) FROM likes 
		        WHERE likeable_type = 'post' AND likeable_id = p.id AND deleted_at IS NULL) as likes_count,
		       EXISTS(SELECT 1 FROM likes 
		              WHERE likeable_type = 'post' AND likeable_id = p.id 
		              AND user_id = ? AND deleted_at IS NULL) as user_liked
	FROM posts p
		JOIN follows f ON p.user_id = f.followed_id
		WHERE f.follower_id = ? AND f.status = 'accepted'
		  AND (p.privacy = 'public' 
		       OR (p.privacy = 'almost_private' AND EXISTS(
		             SELECT 1 FROM follows 
		             WHERE follower_id = p.user_id AND followed_id = ? AND status = 'accepted'
		          ))
		       OR p.user_id = ?)
		ORDER BY p.created_at DESC;
	`
	rows, err := db.Query(stm, userID, userID, userID, userID)
	if err != nil {
		return []Post{}, err
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.Privacy, &post.CreatedAt, &post.LikesCount, &post.UserLiked)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if posts == nil {
		return []Post{}, nil
	}
	return posts, nil
}

func DeletePost(db *sql.DB, postID, userID string) error {
	stmt := `
		DELETE FROM posts
		WHERE id = ? AND user_id = ?;
`
	result, err := db.Exec(stmt, postID, userID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("post not found or not owned by user")
	}

	return nil
}

func GetAllPosts(db *sql.DB, userID string) ([]Post, error) {
	query := `
        SELECT 
            p.id, 
            p.user_id, 
            p.group_id, 
            p.content, 
            p.privacy, 
            p.created_at, 
            p.updated_at, 
            p.deleted_at,
            COUNT(l.id) AS likes_count,
            MAX(CASE WHEN l.user_id = $1 AND l.deleted_at IS NULL THEN 1 ELSE 0 END) AS user_liked
        FROM posts p
        LEFT JOIN likes l ON 
            p.id = l.likeable_id AND 
            l.likeable_type = 'post'
        WHERE p.deleted_at IS NULL
        GROUP BY 
            p.id, p.user_id, p.group_id, p.content, p.privacy,
            p.created_at, p.updated_at, p.deleted_at
        ORDER BY p.created_at DESC
    `

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("database query error: %v", err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.UserID, &post.GroupID, &post.Content, &post.Privacy, &post.CreatedAt, &post.UpdatedAt, &post.DeletedAt, &post.LikesCount, &post.UserLiked)
		if err != nil {
			return nil, fmt.Errorf("error scanning post: %v", err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}
