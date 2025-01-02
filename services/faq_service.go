package services

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	// "github.com/rogerjeasy/go-letusconnect/mappers"
	// "github.com/rogerjeasy/go-letusconnect/models"
)

type FAQService struct {
	firestoreClient *firestore.Client
}

func NewFAQService(client *firestore.Client) *FAQService {
	return &FAQService{
		firestoreClient: client,
	}
}

func (s *FAQService) CreateFAQ(ctx context.Context, data map[string]interface{}) error {
	// Map the data to Go struct format
	faq := data

	// Add the FAQ to Firestore
	_, _, err := s.firestoreClient.Collection("faqs").Add(ctx, faq)
	if err != nil {
		return fmt.Errorf("failed to add FAQ to Firestore: %v", err)
	}

	return nil
}
