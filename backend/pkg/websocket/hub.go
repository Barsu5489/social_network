// pkg/websocket/hub.go
package websocket

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"social-nework/pkg/repository"
)

type Hub struct {
	// Connection management
	Clients      map[string]*Client   // userID -> Client
	ChatRooms    map[string]*ChatRoom // chatID -> ChatRoom
	Register     chan *Client
	Unregister   chan *Client
	MessageQueue chan MessagePayload
	mu           sync.RWMutex

	// Database dependencies
	db          *sql.DB
	messageRepo *repository.MessageRepository
	chatRepo    *repository.ChatRepository
}



// NewHub creates a new Hub instance with all required dependencies
func NewHub(db *sql.DB, messageRepo *repository.MessageRepository, chatRepo *repository.ChatRepository) *Hub {
	return &Hub{
		Clients:      make(map[string]*Client),

		Register:     make(chan *Client, 100),
		Unregister:   make(chan *Client, 100),
	
		db:           db,
		messageRepo:  messageRepo,
		chatRepo:     chatRepo,
	}
}

