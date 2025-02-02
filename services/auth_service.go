package services

import (
	"context"
	"errors"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"google.golang.org/api/iterator"
)

type AuthService struct {
	firestoreClient FirestoreClient
}

func NewAuthService(client FirestoreClient) *AuthService {
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

// GetUserByUsername fetches user data by their username and maps it to models.User
func (s *AuthService) GetUserByUsername(username string) (*models.User, error) {
	ctx := context.Background()

	// Query the Firestore collection to find a user document with the given username
	query := s.firestoreClient.Collection("users").Where("username", "==", username).Limit(1).Documents(ctx)
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

// CreateUser creates a new user in Firestore
func (s *AuthService) CreateUser(ctx context.Context, user *models.User) error {
	backendUser := mappers.MapUserFrontendToBackend(user)
	_, _, err := s.firestoreClient.Collection("users").Add(ctx, backendUser)
	return err
}

// UpdateUser updates an existing user in Firestore
func (s *AuthService) UpdateUser(ctx context.Context, userID string, updates map[string]interface{}) error {
	docRef := s.firestoreClient.Collection("users").Doc(userID)
	_, err := docRef.Set(ctx, updates, firestore.MergeAll) // Use firestore.MergeAll
	return err
}

// GetUserByUID fetches user data by their UID
func (s *AuthService) GetUserByUID(uid string) (*models.User, error) {
	ctx := context.Background()

	query := s.firestoreClient.Collection("users").Where("uid", "==", uid).Limit(1).Documents(ctx)
	defer query.Stop()

	doc, err := query.Next()
	if err == iterator.Done {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, errors.New("failed to fetch user data")
	}

	data := doc.Data()
	user := mappers.MapFrontendToUser(data)
	return &user, nil
}

// CheckExistingUser checks if a user already exists with the given email or username
func (s *AuthService) CheckExistingUser(ctx context.Context, email, username string) error {
	// Check for existing email
	emailQuery := s.firestoreClient.Collection("users").Where("email", "==", email).Limit(1).Documents(ctx)
	defer emailQuery.Stop()

	if _, err := emailQuery.Next(); err != iterator.Done {
		return errors.New("email already exists")
	}

	// Check for existing username
	usernameQuery := s.firestoreClient.Collection("users").Where("username", "==", username).Limit(1).Documents(ctx)
	defer usernameQuery.Stop()

	if _, err := usernameQuery.Next(); err != iterator.Done {
		return errors.New("username already exists")
	}

	return nil
}

func (s *AuthService) CheckEmailorUsername(ctx context.Context, emailOrUsername string) (firestore.Query, error) {
	if strings.Contains(emailOrUsername, "@") {
		return s.firestoreClient.Collection("users").Where("email", "==", emailOrUsername), nil
	} else {
		return s.firestoreClient.Collection("users").Where("username", "==", emailOrUsername), nil
	}
}
