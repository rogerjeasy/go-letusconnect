package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"google.golang.org/api/iterator"
)

type UserService struct {
	firestoreClient FirestoreClient
}

func NewUserService(client FirestoreClient) *UserService {
	return &UserService{
		firestoreClient: client,
	}
}

// GetUserByEmail fetches user data by their email and maps it to models.User
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
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
func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
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

// GetUserRole retrieves the roles of a user by their UID from the "users" collection in Firestore.
func (s *UserService) GetUserRole(uid string) ([]string, error) {
	ctx := context.Background()

	// Query the Firestore collection to find a user document with the given UID
	query := s.firestoreClient.Collection("users").Where("uid", "==", uid).Limit(1).Documents(ctx)
	defer query.Stop()

	doc, err := query.Next()
	if err == iterator.Done {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, errors.New("failed to fetch user data")
	}

	// Extract the role from the document data
	data := doc.Data()
	roleField, exists := data["role"]

	if !exists {
		return nil, errors.New("role not found for the user")
	}

	// Ensure the role field is a slice of strings
	roles, ok := roleField.([]interface{})
	if !ok {
		return nil, errors.New("invalid role format")
	}

	// Convert []interface{} to []string
	var roleList []string
	for _, role := range roles {
		if roleStr, ok := role.(string); ok {
			roleList = append(roleList, roleStr)
		}
	}

	return roleList, nil
}

// GetUsernameByUID fetches the username of a user by their UID
func (s *UserService) GetUsernameByUID(uid string) (string, error) {
	ctx := context.Background()

	// Query the Firestore collection to find a user document with the given UID
	query := s.firestoreClient.Collection("users").Where("uid", "==", uid).Limit(1).Documents(ctx)
	defer query.Stop()

	// Get the next document matching the query
	doc, err := query.Next()
	if err == iterator.Done {
		return "", errors.New("user not found")
	}
	if err != nil {
		return "", errors.New("failed to fetch user data")
	}

	// Extract the username from the document data
	data := doc.Data()
	username, exists := data["username"]
	if !exists {
		return "", errors.New("username not found for the user")
	}

	// Ensure the username is a string
	usernameStr, ok := username.(string)
	if !ok || usernameStr == "" {
		return "", errors.New("invalid username format")
	}

	return usernameStr, nil
}

// GetUserByUID fetches user data by their UID
func (s *UserService) GetUserByUID(uid string) (map[string]interface{}, error) {
	ctx := context.Background()

	// Query the Firestore collection to find a user document with the given UID
	query := s.firestoreClient.Collection("users").Where("uid", "==", uid).Limit(1).Documents(ctx)
	defer query.Stop()

	// Get the next document matching the query
	doc, err := query.Next()
	if err == iterator.Done {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, errors.New("failed to fetch user data")
	}

	// Extract the user data from the document
	data := doc.Data()
	return data, nil
}

func (s *UserService) GetUserByUIDinGoStruct(uid string) (*models.User, error) {
	ctx := context.Background()

	query := s.firestoreClient.Collection("users").Where("uid", "==", uid).Limit(1)
	userDoc, err := query.Documents(ctx).Next()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user data: %v", err)
	}

	userData := userDoc.Data()
	user := mappers.MapBackendToUser(userData)
	return &user, nil
}

// check uniqueness of username
func (s *UserService) CheckUsernameUniqueness(username string) (bool, error) {
	ctx := context.Background()

	query := s.firestoreClient.Collection("users").Where("username", "==", username).Limit(1).Documents(ctx)
	defer query.Stop()

	doc, err := query.Next()
	if err == iterator.Done {
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check username uniqueness: %v", err)
	}

	if doc != nil {
		return false, nil
	}

	return true, nil
}

// check uniqueness of email
func (s *UserService) CheckEmailUniqueness(email string) (bool, error) {
	ctx := context.Background()

	query := s.firestoreClient.Collection("users").Where("email", "==", email).Limit(1).Documents(ctx)
	defer query.Stop()

	doc, err := query.Next()
	if err == iterator.Done {
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check email uniqueness: %v", err)
	}

	if doc != nil {
		return false, nil
	}

	return true, nil
}

// function to get all users
func (s *UserService) GetAllUsers(ctx context.Context) ([]map[string]interface{}, error) {
	iter := s.firestoreClient.Collection("users").Documents(ctx)
	defer iter.Stop()

	var users []map[string]interface{}

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to fetch users: %v", err)
		}

		// Log raw data to check the Firestore document
		data := doc.Data()

		// Map the backend user data to frontend format
		frontendUser := mappers.MapUserBackendToFrontend(data)
		users = append(users, frontendUser)
	}

	return users, nil
}
