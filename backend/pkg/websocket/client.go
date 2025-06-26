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
		return true // Allow all origins for simplicity; adjust as needed
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

	// Fetch user chat IDs from the repository
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

	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		var msg MessagePayload
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		msg.SenderID = c.UserID
		c.Hub.MessageQueue <- msg
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteJSON(message); err != nil {
				log.Println("Write error:", err)
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

