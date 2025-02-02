package services

import (
	"context"
	"fmt"
	// "github.com/rogerjeasy/go-letusconnect/mappers"
	// "github.com/rogerjeasy/go-letusconnect/models"
)

type MessageService struct {
	firestoreClient FirestoreClient
}

func NewMessageService(client FirestoreClient) *MessageService {
	return &MessageService{
		firestoreClient: client,
	}
}

func (s *MessageService) CreateMessage(ctx context.Context, data map[string]interface{}) error {
	// Map the data to Go struct format
	message := data

	// Add the message to Firestore
	_, _, err := s.firestoreClient.Collection("messages").Add(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to add message to Firestore: %v", err)
	}

	return nil
}
