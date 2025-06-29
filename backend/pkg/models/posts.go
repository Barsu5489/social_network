package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func CreatePost(db *sql.DB, ctx context.Context, userID, content, privacy string, groupID *string, allowedUserIDs []string) (string, error) {
	if content == "" {
		return "", errors.New("content cannot be empty")
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
	postStm := `INSERT INTO posts (id, user_id, group_id, content, privacy, created_at, updated_at)
                VALUES(?,?,?,?,?,?,?)`
	_, err = tx.ExecContext(ctx, postStm, postID, userID, groupID, content, privacy, now, now)
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
            u.nickname AS author_nickname,
            u.avatar_url AS author_avatar_url,
            COUNT(l.id) AS likes_count,
            MAX(CASE WHEN l.user_id = $1 AND l.deleted_at IS NULL THEN 1 ELSE 0 END) AS user_liked
        FROM posts p
        JOIN users u ON p.user_id = u.id
        LEFT JOIN likes l ON 
            p.id = l.likeable_id AND 
            l.likeable_type = 'post'
        WHERE p.deleted_at IS NULL
        GROUP BY 
            p.id, p.user_id, p.group_id, p.content, p.privacy,
            p.created_at, p.updated_at, p.deleted_at
        ORDER BY p.created_at DESC
    `

	fmt.Printf("Executing query:\n%s\nWith userID: %s\n", query, userID)

	rows, err := db.Query(query, userID)
	if err != nil {
		fmt.Printf("Query error: %v\n", err)
		return nil, fmt.Errorf("database query error: %v", err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		var groupID sql.NullString
		var deletedAt sql.NullInt64
		var authorAvatarURL sql.NullString

		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&groupID,
			&p.Content,
			&p.Privacy,
			&p.CreatedAt,
			&p.UpdatedAt,
			&deletedAt,
			&p.AuthorNickname,
			&authorAvatarURL,
			&p.LikesCount,
			&p.UserLiked,
		)
		if err != nil {
			fmt.Printf("Row scan error: %v\n", err)
			return nil, fmt.Errorf("row scan error: %v", err)
		}

		// Handle nullable fields
		if groupID.Valid {
			p.GroupID = &groupID.String
		} else {
			p.GroupID = nil
		}

		if deletedAt.Valid {
			p.DeletedAt = &deletedAt.Int64
		} else {
			p.DeletedAt = nil
		}

		if authorAvatarURL.Valid {
			p.AuthorAvatarURL = authorAvatarURL.String
		} else {
			p.AuthorAvatarURL = "" // Default to empty string if NULL
		}

		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("Rows iteration error: %v\n", err)
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	fmt.Printf("Successfully retrieved %d posts\n", len(posts))
	return posts, nil
}
