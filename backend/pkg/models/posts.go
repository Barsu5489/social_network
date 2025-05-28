package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
		fmt.Println(err)
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

		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&groupID,
			&p.Content,
			&p.Privacy,
			&p.CreatedAt,
			&p.UpdatedAt,
			&deletedAt,
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

		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("Rows iteration error: %v\n", err)
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	fmt.Printf("Successfully retrieved %d posts\n", len(posts))
	return posts, nil
}
