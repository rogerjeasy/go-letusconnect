package services

import (
	"context"

	"cloud.google.com/go/firestore"
)

type NewsletterService struct {
	firestoreClient *firestore.Client
}

func NewNewsletterService(client *firestore.Client) *NewsletterService {
	return &NewsletterService{
		firestoreClient: client,
	}
}

func (s *NewsletterService) CreateNewsletter(ctx context.Context) {

}
