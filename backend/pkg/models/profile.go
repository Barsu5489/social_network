package models

import "database/sql"

// If private Only see avater, following and folowers and requst btn
// if public see, posts, following and followers, if your following, avatar,
// If personal Profile - Be able to edit everything and see all details. Have a settings section to maniplate profile

func GetUser(db sql.DB, viewerID, targetID string) (*User, []Post, []User, []User, error) {
	stmt := `
	SELECT id, email, ...
FROM users
WHERE id = ? AND deleted_at IS NULL

`
	result, _ := db.Query(stmt, targetID)

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
	result,	err := db.Query(veiwStmt,viewerID, targetID)
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
	
}
