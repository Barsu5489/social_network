// pkg/websocket/hub.go
package websocket

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"social-nework/pkg/repository"
	"social-nework/pkg/models"
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
	db               *sql.DB
	messageRepo      *repository.MessageRepository
	chatRepo         *repository.ChatRepository
	notificationModel *models.NotificationModel
}
type ChatRoom struct {
	ID           string
	Type         string // "direct" or "group"
	Members      map[string]*Client
	Participants map[string]*Client
	CreatedAt    time.Time
}

type MessagePayload struct {
	Type      string      `json:"type"` // "message", "history_request"
	ChatID    string      `json:"chat_id"`
	SenderID  string      `json:"sender_id"`
	Content   string      `json:"content,omitempty"`
	Timestamp int64   `json:"timestamp,omitempty"`
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

// InitializeChatRoom creates or updates a chat room with participants
func (h *Hub) InitializeChatRoom(chatID, chatType string, participantIDs []string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Create or get chat room
	chatRoom, exists := h.ChatRooms[chatID]
	if !exists {
		chatRoom = &ChatRoom{
			ID:        chatID,
			Type:      chatType,
			Members:   make(map[string]*Client),
			CreatedAt: time.Now(),
		}
		h.ChatRooms[chatID] = chatRoom
	}

	// Add connected participants to the room
	for _, userID := range participantIDs {
		if client, ok := h.Clients[userID]; ok {
			chatRoom.Members[userID] = client

			// Update client's chat list
			client.mu.Lock()
			client.Chats[chatID] = true
			client.mu.Unlock()
		}
	}
}

// AddUserToChatRoom adds a user to an existing chat room
func (h *Hub) AddUserToChatRoom(chatID, userID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	chatRoom, exists := h.ChatRooms[chatID]
	if !exists {
		return
	}

	// Add user if they're connected
	if client, ok := h.Clients[userID]; ok {
		chatRoom.Members[userID] = client

		// Update client's chat list
		client.mu.Lock()
		client.Chats[chatID] = true
		client.mu.Unlock()
	}
}

// RemoveUserFromChatRoom removes a user from a chat room
func (h *Hub) RemoveUserFromChatRoom(chatID, userID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	chatRoom, exists := h.ChatRooms[chatID]
	if !exists {
		return
	}

	delete(chatRoom.Members, userID)

	// Update client's chat list if they're connected
	if client, ok := h.Clients[userID]; ok {
		client.mu.Lock()
		delete(client.Chats, chatID)
		client.mu.Unlock()
	}

	// Clean up empty chat room
	if len(chatRoom.Members) == 0 {
		delete(h.ChatRooms, chatID)
	}
}

// BroadcastToChatRoom sends a message to all members of a chat room except the sender
func (h *Hub) BroadcastToChatRoom(chatID string, message MessagePayload, excludeUserID string) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	log.Printf("BroadcastToChatRoom: Broadcasting to chat %s, excluding %s", chatID, excludeUserID)

	chatRoom, exists := h.ChatRooms[chatID]
	if !exists {
		log.Printf("BroadcastToChatRoom: Chat room %s not found for broadcast", chatID)
		return
	}

	log.Printf("BroadcastToChatRoom: Chat room %s has %d members", chatID, len(chatRoom.Members))

	for userID, client := range chatRoom.Members {
		if userID == excludeUserID {
			log.Printf("BroadcastToChatRoom: Skipping sender %s", userID)
			continue
		}

		log.Printf("BroadcastToChatRoom: Sending to user %s", userID)
		select {
		case client.Send <- message:
			log.Printf("BroadcastToChatRoom: Message sent to user %s", userID)
		default:
			log.Printf("BroadcastToChatRoom: Client %s send buffer full during broadcast", userID)
		}
	}
}

// Complete the initializeUserChatRooms method that was stubbed
func (h *Hub) initializeUserChatRooms(client *Client) {
	log.Printf("initializeUserChatRooms: Initializing chat rooms for user %s", client.UserID)
	
	// Get all chat IDs for the user from database
	chatIDs, err := h.chatRepo.GetUserChatIDs(client.UserID)
	if err != nil {
		log.Printf("initializeUserChatRooms: Error getting user chat IDs: %v", err)
		return
	}

	log.Printf("initializeUserChatRooms: User %s has %d chats", client.UserID, len(chatIDs))

	// Initialize or join each chat room
	for _, chatID := range chatIDs {
		log.Printf("initializeUserChatRooms: Processing chat %s", chatID)
		
		// Get chat participants
		participants, err := h.chatRepo.GetChatParticipants(chatID)
		if err != nil {
			log.Printf("initializeUserChatRooms: Error getting chat participants for %s: %v", chatID, err)
			continue
		}

		log.Printf("initializeUserChatRooms: Chat %s has participants: %v", chatID, participants)

		// Create chat room if it doesn't exist
		chatRoom, exists := h.ChatRooms[chatID]
		if !exists {
			log.Printf("initializeUserChatRooms: Creating new chat room %s", chatID)
			chatRoom = &ChatRoom{
				ID:        chatID,
				Type:      "direct", // You might want to query this from database
				Members:   make(map[string]*Client),
				CreatedAt: time.Now(),
			}
			h.ChatRooms[chatID] = chatRoom
		}

		// Add this client to the chat room
		chatRoom.Members[client.UserID] = client
		client.Chats[chatID] = true
		
		log.Printf("initializeUserChatRooms: Added user %s to chat room %s. Room now has %d members", 
			client.UserID, chatID, len(chatRoom.Members))
	}
}

// GetChatRoomInfo returns information about a chat room
func (h *Hub) GetChatRoomInfo(chatID string) (*ChatRoom, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	chatRoom, exists := h.ChatRooms[chatID]
	return chatRoom, exists
}

// GetConnectedUsers returns a list of connected user IDs
func (h *Hub) GetConnectedUsers() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users := make([]string, 0, len(h.Clients))
	for userID := range h.Clients {
		users = append(users, userID)
	}
	return users
}

// IsUserOnline checks if a user is currently connected
func (h *Hub) IsUserOnline(userID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	_, exists := h.Clients[userID]
	return exists
}

// SendDirectMessage sends a message directly to a specific user
func (h *Hub) SendDirectMessage(userID string, message MessagePayload) bool {
	h.mu.RLock()
	client, exists := h.Clients[userID]
	h.mu.RUnlock()

	if !exists {
		return false
	}

	select {
	case client.Send <- message:
		return true
	default:
		log.Printf("Client %s send buffer full for direct message", userID)
		return false
	}
}

func (h *Hub) cleanupDisconnectedClient(client *Client) {
	for _, chatRoom := range h.ChatRooms {
		delete(chatRoom.Members, client.UserID)
		if len(chatRoom.Members) == 0 {
			delete(h.ChatRooms, chatRoom.ID)
		}
	}
}

// Add method to set notification model
func (h *Hub) SetNotificationModel(notificationModel *models.NotificationModel) {
	h.notificationModel = notificationModel
}
