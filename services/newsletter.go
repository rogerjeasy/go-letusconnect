package services

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type NewsletterService struct {
	FirestoreClient *firestore.Client
}

func NewNewsletterService(client *firestore.Client) *NewsletterService {
	return &NewsletterService{
		FirestoreClient: client,
	}
}

func (s *NewsletterService) GetTotalSubscribers(ctx context.Context) (int, error) {
	if s.FirestoreClient == nil {
		return 0, errors.New("firestore client is not initialized")
	}

	iter := s.FirestoreClient.Collection("newsletters").Documents(ctx)
	defer iter.Stop()

	count := 0
	for {
		_, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, fmt.Errorf("error counting subscribers: %w", err)
		}
		count++
	}

	return count, nil
}

func (s *NewsletterService) CreateNewsletter(ctx context.Context) {

}
