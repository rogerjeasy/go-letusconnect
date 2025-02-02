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
	"github.com/rogerjeasy/go-letusconnect/config"
	"google.golang.org/api/option"
)

var (
	FirebaseApp  *firebase.App
	Firestore    FirestoreClient
	FirebaseAuth *auth.Client
)

// FirestoreClient defines the methods from firestore.Client that are used in UserService.
type FirestoreClient interface {
	Collection(path string) *firestore.CollectionRef
	Doc(path string) *firestore.DocumentRef
	GetAll(ctx context.Context, docRefs []*firestore.DocumentRef) ([]*firestore.DocumentSnapshot, error)
	Batch() *firestore.WriteBatch
	RunTransaction(ctx context.Context, f func(context.Context, *firestore.Transaction) error, opts ...firestore.TransactionOption) error
	Close() error
}

func InitializeFirebase() error {
	// opt := option.WithCredentialsFile(os.Getenv("FIREBASE_SERVICE_ACCOUNT"))

	// Retrieve the base64-encoded service account key from the environment variable
	base64EncodedKey := os.Getenv("FIREBASE_SERVICE_ACCOUNT")
	var opt option.ClientOption

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

	if config.EnvServiceAccount == "local" {
		// Use the hardcoded local credentials file path
		opt = option.WithCredentialsFile(config.JsonServiceAccountPath)
		log.Println("Using local Firebase credentials file")
	} else {
		// Use the temporary file as the credentials file
		opt = option.WithCredentialsFile(tempFile.Name())
	}

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

	Firestore = client
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
