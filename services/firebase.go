package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"

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
	// opt := option.WithCredentialsFile(os.Getenv("FIREBASE_SERVICE_ACCOUNT"))

	// Retrieve the base64-encoded service account key from the environment variable
	base64EncodedKey := os.Getenv("FIREBASE_SERVICE_ACCOUNT")

	log.Println("base64EncodedKey: ", base64EncodedKey)

	// Decode the base64-encoded key
	decodedKey, err := base64.StdEncoding.DecodeString(base64EncodedKey)
	if err != nil {
		return fmt.Errorf("failed to decode service account key: %v", err)
	}

	// Create a temporary file to store the decoded key
	tempFile, err := ioutil.TempFile("", "firebase-service-account")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up the temporary file

	if _, err := tempFile.Write(decodedKey); err != nil {
		return fmt.Errorf("failed to write to temporary file: %v", err)
	}

	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("failed to close temporary file: %v", err)
	}

	// Use the temporary file as the credentials file
	opt := option.WithCredentialsFile(tempFile.Name())

	// Initialize Firebase app
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return fmt.Errorf("error initializing Firebase app: %v", err)
	}

	FirebaseApp = app

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
