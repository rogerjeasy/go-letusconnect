package services

import (
	"context"
)

type ContactUsService struct {
	firestoreClient FirestoreClient
}

func NewContactUsService(client FirestoreClient) *ContactUsService {
	return &ContactUsService{
		firestoreClient: client,
	}
}

func (s *ContactUsService) CreateContactUs(ctx context.Context) {
	// do nothing
}
