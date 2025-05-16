package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) CreatePost(ctx context.Context, userID, content, privacy string, groupID *string) (string, error) {
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
	result, err := m.DB.ExecContext(ctx, stm, id,userID,  groupID, content, privacy,now, now)

	if err != nil{
		return "", err
	}
	rowsAffected, _ := result.RowsAffected()
if rowsAffected == 0 {
    return "", errors.New("post insert failed")
}
	return id, nil
}
