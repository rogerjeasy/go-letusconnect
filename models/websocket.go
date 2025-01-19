package models

import (
	"log"
	"sync"

	"github.com/coder/websocket"
)

type WebSocketMessage struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
	From    string      `json:"from"`
	To      string      `json:"to"`
	Time    int64       `json:"time"`
}

type Client struct {
	ID       string
	Conn     *websocket.Conn
	Manager  *Manager
	Send     chan []byte
	UserID   string
	IsActive bool
}

type Manager struct {
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	mutex      sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte, 256), // Buffered channel for scalability
	}
}

func (m *Manager) Run() {
	for {
		select {
		case client := <-m.Register:
			m.registerClient(client)

		case client := <-m.Unregister:
			m.unregisterClient(client)

		case message := <-m.Broadcast:
			m.broadcastMessage(message)
		}
	}
}

// func (m *Manager) registerClient(client *Client) {
// 	m.mutex.Lock()
// 	defer m.mutex.Unlock()

// 	if _, exists := m.Clients[client.UserID]; exists {
// 		log.Printf("Client already registered: %s", client.UserID)
// 		return
// 	}

// 	m.Clients[client.UserID] = client
// 	log.Printf("Client registered: %s", client.UserID)
// }

func (m *Manager) unregisterClient(client *Client) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.Clients[client.UserID]; exists {
		log.Printf("Unregistering client: %s", client.UserID)
		delete(m.Clients, client.UserID)
		close(client.Send) // Close the send channel to release resources
	}
}

func (m *Manager) broadcastMessage(message []byte) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for userID, client := range m.Clients {
		if client.IsActive {
			select {
			case client.Send <- message:
				log.Printf("Message sent to client: %s", userID)
			default:
				log.Printf("Message send failed; closing client: %s", userID)
				m.unregisterClient(client)
			}
		}
	}
}

func (m *Manager) registerClient(client *Client) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Update existing client if it exists
	if existingClient, exists := m.Clients[client.UserID]; exists {
		existingClient.IsActive = false // Mark old connection as inactive
		close(existingClient.Send)      // Close old send channel
	}

	client.IsActive = true
	m.Clients[client.UserID] = client
	log.Printf("Client registered and active: %s", client.UserID)
}

func (m *Manager) GetClientStatus(userID string) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if client, exists := m.Clients[userID]; exists {
		return client.IsActive
	}
	return false
}
