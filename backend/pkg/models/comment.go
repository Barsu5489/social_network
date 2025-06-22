package models

import (
	"context"
	"database/sql"
	"time"
    
    "github.com/google/uuid"
)

// CreateComment creates a new comment on a post
func CreateComment(db *sql.DB, ctx context.Context, postID, userID, content string, imageURL *string) (*Comment, error) {
    commentID := uuid.New().String()
    now := time.Now().Unix()
    
    stmt := `
        INSERT INTO comments (id, post_id, user_id, content, image_url, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `
    
    _, err := db.ExecContext(ctx, stmt, commentID, postID, userID, content, imageURL, now, now)
    if err != nil {
        return nil, err
    }
    
    // Return the created comment
    return GetCommentByID(db, ctx, commentID, userID)
}

// GetCommentByID retrieves a single comment by ID
func GetCommentByID(db *sql.DB, ctx context.Context, commentID, requestingUserID string) (*Comment, error) {
    stmt := `
        SELECT 
            c.id, c.post_id, c.user_id, c.content, c.image_url, 
            c.created_at, c.updated_at, c.deleted_at,
            u.nickname, u.avatar_url,
            COUNT(l.id) as likes_count,
            EXISTS(
                SELECT 1 FROM likes 
                WHERE user_id = ? AND likeable_type = 'comment' AND likeable_id = c.id AND deleted_at IS NULL
            ) as user_liked
        FROM comments c
        JOIN users u ON c.user_id = u.id
        LEFT JOIN likes l ON l.likeable_id = c.id AND l.likeable_type = 'comment' AND l.deleted_at IS NULL
        WHERE c.id = ? AND c.deleted_at IS NULL
        GROUP BY c.id
    `
    
    var comment Comment
    var imageURL sql.NullString
    var deletedAt sql.NullInt64
    
    err := db.QueryRowContext(ctx, stmt, requestingUserID, commentID).Scan(
        &comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &imageURL,
        &comment.CreatedAt, &comment.UpdatedAt, &deletedAt,
        &comment.UserNickname, &comment.UserAvatar,
        &comment.LikesCount, &comment.UserLiked,
    )
    
    if err != nil {
        return nil, err
    }
    
    if imageURL.Valid {
        comment.ImageURL = &imageURL.String
    }
    
    if deletedAt.Valid {
        comment.DeletedAt = &deletedAt.Int64
    }
    
    return &comment, nil
}

// GetPostComments retrieves all comments for a specific post
func GetPostComments(db *sql.DB, ctx context.Context, postID, requestingUserID string) ([]Comment, error) {
    stmt := `
        SELECT 
            c.id, c.post_id, c.user_id, c.content, c.image_url, 
            c.created_at, c.updated_at, c.deleted_at,
            u.nickname, u.avatar_url,
            COUNT(l.id) as likes_count,
            EXISTS(
                SELECT 1 FROM likes 
                WHERE user_id = ? AND likeable_type = 'comment' AND likeable_id = c.id AND deleted_at IS NULL
            ) as user_liked
        FROM comments c
        JOIN users u ON c.user_id = u.id
        LEFT JOIN likes l ON l.likeable_id = c.id AND l.likeable_type = 'comment' AND l.deleted_at IS NULL
        WHERE c.post_id = ? AND c.deleted_at IS NULL
        GROUP BY c.id
        ORDER BY c.created_at ASC
    `
    
    rows, err := db.QueryContext(ctx, stmt, requestingUserID, postID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var comments []Comment
    for rows.Next() {
        var comment Comment
        var imageURL sql.NullString
        var deletedAt sql.NullInt64
        
        err := rows.Scan(
            &comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &imageURL,
            &comment.CreatedAt, &comment.UpdatedAt, &deletedAt,
            &comment.UserNickname, &comment.UserAvatar,
            &comment.LikesCount, &comment.UserLiked,
        )
        
        if err != nil {
            return nil, err
        }
        
        if imageURL.Valid {
            comment.ImageURL = &imageURL.String
        }
        
        if deletedAt.Valid {
            comment.DeletedAt = &deletedAt.Int64
        }
        
        comments = append(comments, comment)
    }
    
    if comments == nil {
        return []Comment{}, nil
    }
    
    return comments, nil
}


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
