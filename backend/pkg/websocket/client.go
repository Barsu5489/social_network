package websocket

import (

	"log"
	"net/http"
	"sync"
	"time"



	"github.com/gorilla/websocket"
)

type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan MessagePayload
	UserID string
	Chats  map[string]bool
	mu     sync.RWMutex
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Configure properly for production
	},
}

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

// pkg/websocket/client.go
func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// Get userID from context (set by WebSocketAuth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		log.Println("WebSocket connection missing user ID")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

// Get user's existing chats
	chatIDs, err := hub.chatRepo.GetUserChatIDs(userID)
	if err != nil {
		log.Printf("Failed to get user chats: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	client := &Client{
		Hub:    hub,
		Conn:   conn,
		Send:   make(chan MessagePayload, 256),
		UserID: userID,
		Chats:  make(map[string]bool),
	}

	for _, chatID := range chatIDs {
		client.Chats[chatID] = true
	}

	hub.Register <- client


}

