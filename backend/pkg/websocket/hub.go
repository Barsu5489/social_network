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
type ChatRoom struct {
	ID        string
	Type      string // "direct" or "group"
	Members   map[string]*Client
	CreatedAt time.Time
}

type MessagePayload struct {
	Type      string      `json:"type"` // "message", "history_request"
	ChatID    string      `json:"chat_id"`
	SenderID  string      `json:"sender_id"`
	Content   string      `json:"content,omitempty"`
	Timestamp time.Time   `json:"timestamp,omitempty"`
	Data      interface{} `json:"data,omitempty"` // For additional payload
}

// NewHub creates a new Hub instance with all required dependencies
func NewHub(db *sql.DB, messageRepo *repository.MessageRepository, chatRepo *repository.ChatRepository) *Hub {
	return &Hub{
		Clients:      make(map[string]*Client),
		ChatRooms:    make(map[string]*ChatRoom),
		Register:     make(chan *Client, 100),
		Unregister:   make(chan *Client, 100),
		MessageQueue: make(chan MessagePayload, 1000),
		db:           db,
		messageRepo:  messageRepo,
		chatRepo:     chatRepo,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.UserID] = client
			h.initializeUserChatRooms(client)
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			delete(h.Clients, client.UserID)
			h.cleanupDisconnectedClient(client)
			h.mu.Unlock()

		case msg := <-h.MessageQueue:
			h.mu.RLock()
			switch msg.Type {
			case "message":
				h.handleNewMessage(msg)
			case "history_request":
				h.handleHistoryRequest(msg)
			}
			h.mu.RUnlock()
		}
	}
}

// pkg/websocket/hub_extensions.go - Add these methods to your existing Hub struct

// InitializeChatRoom creates or updates a chat room with participants






// Complete the initializeUserChatRooms method that was stubbed
func (h *Hub) initializeUserChatRooms(client *Client) {
	// Get all chat IDs for the user from database
	chatIDs, err := h.chatRepo.GetUserChatIDs(client.UserID)
	if err != nil {
		log.Printf("Error getting user chat IDs: %v", err)
		return
	}

	// Initialize or join each chat room
	for _, chatID := range chatIDs {
		// Get chat participants
		participants, err := h.chatRepo.GetChatParticipants(chatID)
		if err != nil {
			log.Printf("Error getting chat participants for %s: %v", chatID, err)
			continue
		}

		// Create chat room if it doesn't exist
		chatRoom, exists := h.ChatRooms[chatID]
		if !exists {
			chatRoom = &ChatRoom{
				ID:        chatID,
				Type:      "direct", // You might want to query this from database
				Members:   make(map[string]*Client),
				CreatedAt: time.Now(),
			}
			h.ChatRooms[chatID] = chatRoom
		}

		// Add client to chat room
		chatRoom.Members[client.UserID] = client

		// Update client's chat list
		client.Chats[chatID] = true

		// Add other connected participants to the room
		for _, participantID := range participants {
			if participantID != client.UserID {
				if otherClient, ok := h.Clients[participantID]; ok {
					chatRoom.Members[participantID] = otherClient
				}
			}
		}
	}
}
