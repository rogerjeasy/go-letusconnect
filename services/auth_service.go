package services

import (
	"context"
	"errors"

	"cloud.google.com/go/firestore"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"google.golang.org/api/iterator"
)

type AuthService struct {
	firestoreClient *firestore.Client
}

func NewAuthService(client *firestore.Client) *AuthService {
	return &AuthService{
		firestoreClient: client,
	}
}

// GetUserByEmail fetches user data by their email and maps it to models.User
func (s *AuthService) GetUserByEmail(email string) (*models.User, error) {
	ctx := context.Background()

	// Query the Firestore collection to find a user document with the given email
	query := s.firestoreClient.Collection("users").Where("email", "==", email).Limit(1).Documents(ctx)
	defer query.Stop()

	doc, err := query.Next()
	if err == iterator.Done {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, errors.New("failed to fetch user data")
	}

	// Map the Firestore document data to models.User
	data := doc.Data()
	user := mappers.MapFrontendToUser(data)
	return &user, nil
}
