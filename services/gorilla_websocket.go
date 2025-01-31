package services

import (
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// WebSocketClient represents a WebSocket client
type WebSocketClient struct {
	UserID string
	Conn   *websocket.Conn
}

// WebSocketService manages WebSocket connections
type GorillaWebSocketService struct {
	clients    map[string]*WebSocketClient
	register   chan *WebSocketClient
	unregister chan *WebSocketClient
	broadcast  chan Message
	mutex      sync.RWMutex
}

// Message represents a chat message
type Message struct {
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Content    string `json:"content"`
	Timestamp  int64  `json:"timestamp"`
}

// NewWebSocketService creates a new WebSocket service
func NewGorillaWebSocketService() *GorillaWebSocketService {
	service := &GorillaWebSocketService{
		clients:    make(map[string]*WebSocketClient),
		register:   make(chan *WebSocketClient),
		unregister: make(chan *WebSocketClient),
		broadcast:  make(chan Message),
	}
	go service.run()
	return service
}

// Run manages WebSocket client connections
func (s *GorillaWebSocketService) run() {
	for {
		select {
		case client := <-s.register:
			s.mutex.Lock()
			s.clients[client.UserID] = client
			s.mutex.Unlock()

		case client := <-s.unregister:
			s.mutex.Lock()
			delete(s.clients, client.UserID)
			s.mutex.Unlock()

		case message := <-s.broadcast:
			s.mutex.RLock()
			if receiverClient, exists := s.clients[message.ReceiverID]; exists {
				if err := receiverClient.Conn.WriteJSON(message); err != nil {
					log.Printf("Error sending message: %v", err)
				}
			}
			s.mutex.RUnlock()
		}
	}
}

// HandleWebSocket sets up WebSocket connection handler
func (s *GorillaWebSocketService) HandleWebSocket(app *fiber.App) {
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:userID", websocket.New(func(conn *websocket.Conn) {
		userID := conn.Params("userID")
		client := &WebSocketClient{
			UserID: userID,
			Conn:   conn,
		}

		// Register client
		s.register <- client

		// Unregister on connection close
		defer func() {
			s.unregister <- client
			conn.Close()
		}()

		// Message handling loop
		for {
			var msg Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Printf("WebSocket read error: %v", err)
				break
			}

			// Set sender ID and broadcast
			msg.SenderID = userID
			s.broadcast <- msg
		}
	}))
}

// BroadcastMessage allows sending messages from server-side
func (s *GorillaWebSocketService) GorillaBroadcastMessage(message Message) {
	s.broadcast <- message
}
