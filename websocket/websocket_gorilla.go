package websocket

import (
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
)

// Client represents a WebSocket client
type Client struct {
	ID   string
	Conn *websocket.Conn
}

// WebSocketServer manages WebSocket connections
type WebSocketServer struct {
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
	Mutex      sync.RWMutex
}

// Message represents a chat message
type Message struct {
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Content    string `json:"content"`
}

// NewWebSocketServer creates a new WebSocket server
func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message),
	}
}

// Run starts the WebSocket server
func (server *WebSocketServer) Run() {
	for {
		select {
		case client := <-server.Register:
			server.Mutex.Lock()
			server.Clients[client.ID] = client
			server.Mutex.Unlock()

		case client := <-server.Unregister:
			server.Mutex.Lock()
			delete(server.Clients, client.ID)
			server.Mutex.Unlock()

		case message := <-server.Broadcast:
			server.Mutex.RLock()
			if receiverClient, exists := server.Clients[message.ReceiverID]; exists {
				err := receiverClient.Conn.WriteJSON(message)
				if err != nil {
					log.Printf("Error sending message: %v", err)
				}
			}
			server.Mutex.RUnlock()
		}
	}
}
