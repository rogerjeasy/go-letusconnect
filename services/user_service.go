package services

import (
	"context"
	"errors"

	"google.golang.org/api/iterator"
)

// GetUserRole retrieves the roles of a user by their UID from the "users" collection in Firestore.
func GetUserRole(uid string) ([]string, error) {
	ctx := context.Background()

	// Query the Firestore collection to find a user document with the given UID
	query := FirestoreClient.Collection("users").Where("uid", "==", uid).Limit(1).Documents(ctx)
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
func GetUsernameByUID(uid string) (string, error) {
	ctx := context.Background()

	// Query the Firestore collection to find a user document with the given UID
	query := FirestoreClient.Collection("users").Where("uid", "==", uid).Limit(1).Documents(ctx)
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
