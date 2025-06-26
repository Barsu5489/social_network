package groups

import (
	"database/sql"

	"social-nework/pkg/repository"
	"social-nework/pkg/websocket"
)

type GroupHandler struct {
	db        *sql.DB
	groupRepo *repository.GroupRepository
	chatRepo  *repository.ChatRepository
	h         *websocket.Hub
}

func NewGroupHandler(
	db *sql.DB,
	groupRepo *repository.GroupRepository,
	chatRepo *repository.ChatRepository,
	hub *websocket.Hub,
) *GroupHandler {
	return &GroupHandler{
		db:        db,
		groupRepo: groupRepo,
		chatRepo:  chatRepo,
		h:         hub,
	}
}
