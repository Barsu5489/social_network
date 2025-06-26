// pkg/handlers/chat_handlers.go
package handlers

import (
	"social-nework/pkg/repository"
	"social-nework/pkg/websocket"
)

type ChatHandler struct {
	chatRepo    *repository.ChatRepository
	messageRepo *repository.MessageRepository
	groupRepo   *repository.GroupRepository // Add this for group chat integration
	hub         *websocket.Hub
}

func NewChatHandler(chatRepo *repository.ChatRepository, messageRepo *repository.MessageRepository, groupRepo *repository.GroupRepository, hub *websocket.Hub) *ChatHandler {
	return &ChatHandler{
		chatRepo:    chatRepo,
		messageRepo: messageRepo,
		groupRepo:   groupRepo,
		hub:         hub,
	}
}
