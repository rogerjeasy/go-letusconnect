package services

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

var (
	FirebaseApp     *firebase.App
	FirestoreClient *firestore.Client
	FirebaseAuth    *auth.Client
)

func InitializeFirebase() error {
	opt := option.WithCredentialsFile("C:/Users/rogej/Downloads/go-connect-4f9f8-firebase-adminsdk-3qjyg-2a778c1517.json")

	// Initialize Firebase app
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return fmt.Errorf("error initializing Firebase app: %v", err)
	}

	FirebaseApp = app
	log.Println("Firebase initialized successfully")

	// Initialize Firestore client
	client, err := app.Firestore(context.Background())
	if err != nil {
		return fmt.Errorf("error initializing Firestore client: %v", err)
	}

	FirestoreClient = client
	log.Println("Firestore client initialized successfully")

	// Initialize Firebase Auth client
	authClient, err := app.Auth(context.Background())
	if err != nil {
		return fmt.Errorf("error initializing Firebase Auth client: %v", err)
	}

	FirebaseAuth = authClient
	log.Println("Firebase Auth client initialized successfully")

	return nil
}
