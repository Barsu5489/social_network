package groups

import (
	"database/sql"

	"social-nework/pkg/models"
	"social-nework/pkg/repository"
	"social-nework/pkg/websocket"
)

type GroupHandler struct {
	db                *sql.DB
	groupRepo         *repository.GroupRepository
	chatRepo          *repository.ChatRepository
	h                 *websocket.Hub
	NotificationModel *models.NotificationModel
}

func NewGroupHandler(
	db *sql.DB,
	groupRepo *repository.GroupRepository,
	chatRepo *repository.ChatRepository,
	hub *websocket.Hub,
	notificationModel *models.NotificationModel,
) *GroupHandler {
	return &GroupHandler{
		db:                db,
		groupRepo:         groupRepo,
		chatRepo:          chatRepo,
		h:                 hub,
		NotificationModel: notificationModel,
	}
}
