package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type FollowModel struct {
	DB *sql.DB
}

func (m *FollowModel) Follow(ctx context.Context, followerID, followedID string) error {
	if followerID == followedID {
		return errors.New("cannot follow yourself")
	}

	now := time.Now().Unix() // 8:35 PM EAT, May 15, 2025 = 1744732500
	id := uuid.New().String()

	stmt := `
        INSERT INTO follows (id, follower_id, followed_id, status, created_at)
        VALUES (?, ?, ?, ?, ?)
        ON CONFLICT (follower_id, followed_id) DO NOTHING;
		`
	result, err := m.DB.ExecContext(ctx, stmt, id, followerID, followedID, "accepted", now)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("follow already exists")
	}

	return nil
}

func (m *FollowModel) Unfollow(ctx context.Context, followerID, followedID string) error {
	stmt := `
        DELETE FROM follows
        WHERE follower_id = ? AND followed_id = ?;
    `
	result, err := m.DB.ExecContext(ctx, stmt, followerID, followedID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("follow does not exist")
	}

	return nil
}

func (m *FollowModel) GetFollowers(ctx context.Context, userID string) ([]User, error) {
	stmt := `
        SELECT u.id, u.email, u.first_name, u.last_name, u.nickname, u.date_of_birth,
               u.about_me, u.avatar_url, u.is_private, u.created_at, u.updated_at
        FROM users u
        JOIN follows f ON u.id = f.follower_id
        WHERE f.followed_id = ?;
    `
	rows, err := m.DB.QueryContext(ctx, stmt, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var followers []User
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Nickname,
			&user.DateOfBirth,
			&user.AboutMe,
			&user.AvatarURL,
			&user.IsPrivate,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		followers = append(followers, user)
	}
	if followers == nil {
		followers = []User{}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return followers, nil
}

func (m *FollowModel) GetFollowing(ctx context.Context, userID string) ([]User, error) {
	stmt := `
        SELECT u.id, u.email, u.first_name, u.last_name, u.nickname, u.date_of_birth,
               u.about_me, u.avatar_url, u.is_private, u.created_at, u.updated_at
        FROM users u
        JOIN follows f ON u.id = f.followed_id
        WHERE f.follower_id = ?;
    `
	rows, err := m.DB.QueryContext(ctx, stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var following []User
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Nickname,
			&user.DateOfBirth,
			&user.AboutMe,
			&user.AvatarURL,
			&user.IsPrivate,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		following = append(following, user)
	}
	// Initialize empty slice if no rows returned
	if following == nil {
		following = []User{}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return following, nil
}
