package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"
)

// If private Only see avater, following and folowers and requst btn
// if public see, posts, following and followers, if your following, avatar,
// If personal Profile - Be able to edit everything and see all details. Have a settings section to maniplate profile

func GetUser(db *sql.DB, viewerID, targetID string) (*User, []Post, []User, []User, error) {
	stmt := `
	 SELECT id, email, first_name, last_name, nickname, date_of_birth, about_me, avatar_url, is_private, created_at
        FROM users
        WHERE id = ? AND deleted_at IS NULL
`
	result := db.QueryRow(stmt, targetID)

	user := &User{}
	err := result.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Nickname, &user.DateOfBirth, &user.AboutMe, &user.AvatarURL, &user.IsPrivate, &user.CreatedAt)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	canViewFull := viewerID == targetID || !user.IsPrivate
	veiwStmt := `
SELECT COUNT(*)
		FROM follows
		WHERE follower_id = ? AND followed_id = ? AND status = 'accepted' AND deleted_at IS NULL
`
	if !canViewFull {
		var isFollower int
		result, err := db.Query(veiwStmt, viewerID, targetID)
		result.Scan(&isFollower)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		if isFollower == 0 {
			// Return limited user (id, nickname only)
			limitedUser := &User{ID: user.ID, Nickname: user.Nickname}
			return limitedUser, nil, nil, nil, nil
		}
	}
	// Get posts
	var posts []Post
	if canViewFull {

		rows, err := db.Query(`
			   SELECT id, user_id, group_id, content, privacy, created_at
			   FROM posts
			   WHERE user_id = ? AND deleted_at IS NULL`, targetID)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		defer rows.Close()
		for rows.Next() {
			post := Post{}
			var groupID sql.NullString
			if err := rows.Scan(&post.ID, &post.UserID, &groupID, &post.Content, &post.Privacy, &post.CreatedAt); err != nil {
				return nil, nil, nil, nil, err
			}
			if groupID.Valid {
				post.GroupID = &groupID.String
			}
			posts = append(posts, post)
		}
	}
	// Get followers
	var followers []User
	rows, err := db.Query(`
        SELECT u.id, u.first_name, u.last_name, u.nickname
        FROM users u
        JOIN follows f ON u.id = f.follower_id
        WHERE f.followed_id = ? AND f.status = 'accepted' AND u.deleted_at IS NULL AND f.deleted_at IS NULL`, targetID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer rows.Close()
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Nickname); err != nil {
			return nil, nil, nil, nil, err
		}
		followers = append(followers, user)
	}

	// Get following
	var following []User
	rows, err = db.Query(`
        SELECT u.id, u.first_name, u.last_name, u.nickname
        FROM users u
        JOIN follows f ON u.id = f.followed_id
        WHERE f.follower_id = ? AND f.status = 'accepted' AND u.deleted_at IS NULL AND f.deleted_at IS NULL`, targetID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer rows.Close()
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Nickname); err != nil {
			return nil, nil, nil, nil, err
		}
		following = append(following, user)
	}

	return user, posts, followers, following, nil
}

// UpdateProfile updates user fields based on provided updates

func UpdateProfile(db *sql.DB, userID string, updates map[string]interface{}) error {
	query := `UPDATE users SET updated_at = ?`
	args := []interface{}{time.Now().Unix()}
	var fields []string

	if nickname, ok := updates["nickname"].(string); ok && nickname != "" {
		fields = append(fields, "nickname = ?")
		args = append(args, nickname)
	}
	if aboutMe, ok := updates["about_me"].(string); ok && aboutMe != "" {
		fields = append(fields, "about_me = ?")
		args = append(args, aboutMe)
	}
	if avatarURL, ok := updates["avatar_url"].(string); ok && avatarURL != "" {
		fields = append(fields, "avatar_url = ?")
		args = append(args, avatarURL)
	}
	if isPrivate, ok := updates["is_private"].(bool); ok {
		fields = append(fields, "is_private = ?")
		args = append(args, isPrivate)
	}

	if len(fields) == 0 {
		return errors.New("no valid fields to update")
	}

	query += ", " + strings.Join(fields, ", ") + " WHERE id = ? AND deleted_at IS NULL"
	args = append(args, userID)

	result, err := db.Exec(query, args...)
	if err != nil {
		return err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return errors.New("user not found")
	}
	return nil
}
