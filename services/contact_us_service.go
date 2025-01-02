package services

import (
	"context"

	"cloud.google.com/go/firestore"
)

type ContactUsService struct {
	firestoreClient *firestore.Client
}

func NewContactUsService(client *firestore.Client) *ContactUsService {
	return &ContactUsService{
		firestoreClient: client,
	}
}

func (s *ContactUsService) CreateContactUs(ctx context.Context) {
	// do nothing
}
