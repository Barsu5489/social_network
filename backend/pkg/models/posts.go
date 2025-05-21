package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

func CreatePost(db *sql.DB, ctx context.Context, userID, content, privacy string, groupID *string) (string, error) {
	if content == "" {
		return "", errors.New("content cannot be empty")
	}
	if privacy != "public" && privacy != "almost_private" && privacy != "private" {
		return "", errors.New("invalid privacy setting")
	}
	id := uuid.New().String()
	now := time.Now().Unix()

	stm := `INSERT INTO posts (id, user_id, group_id, content, privacy, created_at, updated_at)
	VALUES(?,?,?,?,?,?,?)
	`
	result, err := db.ExecContext(ctx, stm, id, userID, groupID, content, privacy, now, now)
	if err != nil {
		return "", err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return "", errors.New("post insert failed")
	}
	return id, nil
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
	rows, err := db.Query(stm, userID)
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
