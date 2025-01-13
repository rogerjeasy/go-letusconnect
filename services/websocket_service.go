package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gofiber/websocket/v2"
	"github.com/rogerjeasy/go-letusconnect/models"
)

type WebSocketService struct {
	connections     sync.Map
	firestoreClient *firestore.Client
}

func NewWebSocketService(firestoreClient *firestore.Client) *WebSocketService {
	if firestoreClient == nil {
		log.Fatal("firestoreClient cannot be nil")
	}
	return &WebSocketService{
		firestoreClient: firestoreClient,
	}
}

func (s *WebSocketService) RegisterConnection(userID string, conn *websocket.Conn) error {
	if userID == "" {
		return fmt.Errorf("userID cannot be empty")
	}
	if conn == nil {
		return fmt.Errorf("websocket connection cannot be nil")
	}

	s.connections.Store(userID, &models.WebSocketConnection{
		UserID: userID,
		Conn:   conn,
	})
	log.Printf("User connected: %s", userID)
	return nil
}

func (s *WebSocketService) UnregisterConnection(userID string) error {
	if userID == "" {
		return fmt.Errorf("userID cannot be empty")
	}

	s.connections.Delete(userID)
	log.Printf("User disconnected: %s", userID)
	return nil
}

func (s *WebSocketService) HandleDirectMessage(message models.WebSocketMessage) error {
	// Validate message
	if err := s.validateMessage(message); err != nil {
		return fmt.Errorf("invalid message: %w", err)
	}

	// Add timestamp if not present
	if message.Time == 0 {
		message.Time = time.Now().Unix()
	}

	// Store message in Firestore
	docRef, writeResult, err := s.firestoreClient.Collection("messages").Add(context.Background(), message)
	if err != nil {
		return fmt.Errorf("failed to store message: %w", err)
	}

	log.Printf("Message stored with ID: %s at time: %v", docRef.ID, writeResult.UpdateTime)

	// Send to recipient if online
	if conn, ok := s.getConnection(message.To); ok {
		if err := s.sendMessage(conn, message); err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
	}

	return nil
}

func (s *WebSocketService) HandleConnectionRequest(message models.WebSocketMessage) error {
	if err := s.validateMessage(message); err != nil {
		return fmt.Errorf("invalid connection request: %w", err)
	}

	// Store connection request
	docRef, writeResult, err := s.firestoreClient.Collection("connection_requests").Add(context.Background(), message)
	if err != nil {
		return fmt.Errorf("failed to store connection request: %w", err)
	}

	log.Printf("Connection request stored with ID: %s at time: %v", docRef.ID, writeResult.UpdateTime)

	// Notify recipient if online
	if conn, ok := s.getConnection(message.To); ok {
		if err := s.sendMessage(conn, message); err != nil {
			return fmt.Errorf("failed to send connection request: %w", err)
		}
	}

	return nil
}

func (s *WebSocketService) HandleNotification(message models.WebSocketMessage) error {
	if err := s.validateMessage(message); err != nil {
		return fmt.Errorf("invalid notification: %w", err)
	}

	// Store notification
	docRef, writeResult, err := s.firestoreClient.Collection("notifications").Add(context.Background(), message)
	if err != nil {
		return fmt.Errorf("failed to store notification: %w", err)
	}

	log.Printf("Notification stored with ID: %s at time: %v", docRef.ID, writeResult.UpdateTime)

	// Send to recipient if online
	if conn, ok := s.getConnection(message.To); ok {
		if err := s.sendMessage(conn, message); err != nil {
			return fmt.Errorf("failed to send notification: %w", err)
		}
	}

	return nil
}

func (s *WebSocketService) sendMessage(conn *models.WebSocketConnection, message models.WebSocketMessage) error {
	if conn == nil {
		return fmt.Errorf("connection cannot be nil")
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	if err := conn.Conn.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

func (s *WebSocketService) getConnection(userID string) (*models.WebSocketConnection, bool) {
	if userID == "" {
		return nil, false
	}

	if conn, ok := s.connections.Load(userID); ok {
		return conn.(*models.WebSocketConnection), true
	}
	return nil, false
}

// Helper methods
func (s *WebSocketService) validateMessage(message models.WebSocketMessage) error {
	if message.From == "" {
		return fmt.Errorf("message sender (From) cannot be empty")
	}
	if message.To == "" {
		return fmt.Errorf("message recipient (To) cannot be empty")
	}
	if message.Type == "" {
		return fmt.Errorf("message type cannot be empty")
	}
	return nil
}

// Additional helper method to check if a user is online
func (s *WebSocketService) IsUserOnline(userID string) bool {
	_, ok := s.getConnection(userID)
	return ok
}

// Method to broadcast a message to multiple users
func (s *WebSocketService) BroadcastMessage(message models.WebSocketMessage, recipients []string) error {
	if len(recipients) == 0 {
		return fmt.Errorf("recipients list cannot be empty")
	}

	var errs []error
	for _, recipientID := range recipients {
		if conn, ok := s.getConnection(recipientID); ok {
			if err := s.sendMessage(conn, message); err != nil {
				errs = append(errs, fmt.Errorf("failed to send to %s: %w", recipientID, err))
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("broadcast errors: %v", errs)
	}
	return nil
}
