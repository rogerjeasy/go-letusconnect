package services

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserConnectionService struct {
	firestoreClient *firestore.Client
	UserService     *UserService
}

func NewUserConnectionService(client *firestore.Client, userService *UserService) *UserConnectionService {
	return &UserConnectionService{
		firestoreClient: client,
		UserService:     userService,
	}
}

func (s *UserConnectionService) CreateUserConnections(ctx context.Context, uid string) (*models.UserConnections, error) {
	connections := models.UserConnections{
		ID:              uuid.New().String(),
		UID:             uid,
		Connections:     make(map[string]models.Connection),
		PendingRequests: make(map[string]models.ConnectionRequest),
		SentRequests:    make(map[string]models.SentRequest),
	}

	_, err := s.firestoreClient.Collection("user_connections").Doc(connections.ID).Set(ctx, mappers.MapConnectionsGoToFirestore(connections))
	if err != nil {
		return nil, fmt.Errorf("failed to create user connections: %v", err)
	}

	return &connections, nil
}

func (s *UserConnectionService) CheckUserConnectionsExist(ctx context.Context, uid string) (bool, string, error) {
	query := s.firestoreClient.Collection("user_connections").Where("uid", "==", uid).Limit(1)
	iter := query.Documents(ctx)
	doc, err := iter.Next()

	if err != nil {
		if status.Code(err) == codes.NotFound || err == iterator.Done {
			return false, "", nil
		}
		return false, "", fmt.Errorf("failed to check user connections: %v", err)
	}

	return true, doc.Ref.ID, nil
}

// Update GetUserConnections to use this check
func (s *UserConnectionService) GetUserConnections(ctx context.Context, uid string) (*models.UserConnections, error) {
	if s.firestoreClient == nil {
		return nil, fmt.Errorf("firestore client not initialized")
	}

	exists, docID, err := s.CheckUserConnectionsExist(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to check user connections: %v", err)
	}

	if !exists {
		// Create new connections if none exist
		return s.CreateUserConnections(ctx, uid)
	}

	// Get existing connections
	doc, err := s.firestoreClient.Collection("user_connections").Doc(docID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user connections document: %v", err)
	}

	data := doc.Data()
	if data == nil {
		// If document exists but is empty, create new connections
		return s.CreateUserConnections(ctx, uid)
	}

	connections := mappers.MapConnectionsFirestoreToGo(data)
	return &connections, nil
}

func (s *UserConnectionService) SendConnectionRequest(ctx context.Context, fromUID, toUID, message string) error {
	// First check if users exist and get their connections
	fromConnections, err := s.GetUserConnections(ctx, fromUID)
	if err != nil {
		return fmt.Errorf("failed to get sender's connections: %v", err)
	}

	toConnections, err := s.GetUserConnections(ctx, toUID)
	if err != nil {
		return fmt.Errorf("failed to get recipient's connections: %v", err)
	}

	// Check if users are already connected
	if _, exists := fromConnections.Connections[toUID]; exists {
		return fmt.Errorf("you are already connected with this user")
	}

	// Check if there's an existing sent request
	if sentReq, exists := fromConnections.SentRequests[toUID]; exists {
		if sentReq.Status == "pending" {
			return fmt.Errorf("you already have a pending connection request to this user")
		}
	}

	// Check if there's an existing pending request from the target user
	if pendingReq, exists := fromConnections.PendingRequests[toUID]; exists {
		if pendingReq.Status == "pending" {
			return fmt.Errorf("this user has already sent you a connection request")
		}
	}

	// Start a transaction for the actual request sending
	err = s.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		fromUsername, err := s.UserService.GetUsernameByUID(fromUID)
		if err != nil {
			return fmt.Errorf("failed to get username: %v", err)
		}

		// Create pending request for recipient
		request := models.ConnectionRequest{
			FromUID:  fromUID,
			FromName: fromUsername,
			ToUID:    toUID,
			Message:  message,
			SentAt:   time.Now(),
			Status:   "pending",
		}

		// Create sent request for sender
		sentRequest := models.SentRequest{
			ToUID:   toUID,
			SentAt:  request.SentAt,
			Message: message,
			Status:  "pending",
		}

		// Update both users' records
		toConnections.PendingRequests[fromUID] = request
		fromConnections.SentRequests[toUID] = sentRequest

		// Save both updates in transaction
		err = tx.Set(s.firestoreClient.Collection("user_connections").Doc(fromConnections.ID),
			mappers.MapConnectionsGoToFirestore(*fromConnections))
		if err != nil {
			return fmt.Errorf("failed to update sender's connections: %v", err)
		}

		err = tx.Set(s.firestoreClient.Collection("user_connections").Doc(toConnections.ID),
			mappers.MapConnectionsGoToFirestore(*toConnections))
		if err != nil {
			return fmt.Errorf("failed to update recipient's connections: %v", err)
		}

		if err := SendConnectionRequestNotification(ctx, fromUID, fromUsername, message, toUID); err != nil {
			fmt.Printf("Failed to send connection request notification: %v\n", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to process connection request: %v", err)
	}

	return nil
}

func (s *UserConnectionService) AcceptConnectionRequest(ctx context.Context, fromUID, toUID string) error {
	return s.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		fromConnections, err := s.GetUserConnections(ctx, fromUID)
		if err != nil {
			return err
		}
		toConnections, err := s.GetUserConnections(ctx, toUID)
		if err != nil {
			return err
		}

		uidUsername, err := s.UserService.GetUsernameByUID(toUID)
		if err != nil {
			return fmt.Errorf("failed to get username: %v", err)
		}

		now := time.Now()
		connection1 := models.Connection{
			TargetUID:  toUID,
			TargetName: uidUsername,
			SentAt:     fromConnections.SentRequests[toUID].SentAt,
			AcceptedAt: now,
			Status:     "active",
		}

		fromUsername, err := s.UserService.GetUsernameByUID(fromUID)
		if err != nil {
			return fmt.Errorf("failed to get username: %v", err)
		}

		connection2 := models.Connection{
			TargetUID:  fromUID,
			TargetName: fromUsername,
			SentAt:     toConnections.PendingRequests[fromUID].SentAt,
			AcceptedAt: now,
			Status:     "active",
		}

		// Update connections and manage requests
		fromConnections.Connections[toUID] = connection1
		toConnections.Connections[fromUID] = connection2

		// Update sent request status
		if sentReq, exists := fromConnections.SentRequests[toUID]; exists {
			sentReq.Status = "accepted"
			sentReq.Accepted = now
			fromConnections.SentRequests[toUID] = sentReq
		}

		delete(toConnections.PendingRequests, fromUID)
		delete(fromConnections.SentRequests, toUID)

		// Save updates
		err = tx.Set(s.firestoreClient.Collection("user_connections").Doc(fromConnections.ID),
			mappers.MapConnectionsGoToFirestore(*fromConnections))
		if err != nil {
			return err
		}

		err = tx.Set(s.firestoreClient.Collection("user_connections").Doc(toConnections.ID),
			mappers.MapConnectionsGoToFirestore(*toConnections))
		if err != nil {
			return err
		}

		if err := SendConnectionAcceptedNotification(ctx, toUID, fromUsername, uidUsername, fromUID); err != nil {
			fmt.Printf("Failed to send connection accepted notification: %v\n", err)
		}

		return nil
	})
}

func (s *UserConnectionService) RejectConnectionRequest(ctx context.Context, fromUID, toUID string) error {
	return s.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		fromConnections, err := s.GetUserConnections(ctx, fromUID)
		if err != nil {
			return err
		}

		toConnections, err := s.GetUserConnections(ctx, toUID)
		if err != nil {
			return err
		}

		// Update sent request status
		if sentReq, exists := fromConnections.SentRequests[toUID]; exists {
			sentReq.Status = "rejected"
			fromConnections.SentRequests[toUID] = sentReq
		}

		delete(toConnections.PendingRequests, fromUID)

		// Save updates
		err = tx.Set(s.firestoreClient.Collection("user_connections").Doc(fromConnections.ID),
			mappers.MapConnectionsGoToFirestore(*fromConnections))
		if err != nil {
			return err
		}

		err = tx.Set(s.firestoreClient.Collection("user_connections").Doc(toConnections.ID),
			mappers.MapConnectionsGoToFirestore(*toConnections))
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *UserConnectionService) RemoveConnection(ctx context.Context, uid1, uid2 string) error {
	err := s.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// Get both users' connections
		conn1, err := s.GetUserConnections(ctx, uid1)
		if err != nil {
			return err
		}
		conn2, err := s.GetUserConnections(ctx, uid2)
		if err != nil {
			return err
		}

		// Remove connections
		delete(conn1.Connections, uid2)
		delete(conn2.Connections, uid1)

		// Update both documents
		err = tx.Set(s.firestoreClient.Collection("user_connections").Doc(conn1.ID),
			mappers.MapConnectionsGoToFirestore(*conn1))
		if err != nil {
			return err
		}

		err = tx.Set(s.firestoreClient.Collection("user_connections").Doc(conn2.ID),
			mappers.MapConnectionsGoToFirestore(*conn2))
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (s *UserConnectionService) CancelSentRequest(ctx context.Context, fromUID, toUID string) error {
	return s.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// Get both users' connections
		fromConnections, err := s.GetUserConnections(ctx, fromUID)
		if err != nil {
			return err
		}

		toConnections, err := s.GetUserConnections(ctx, toUID)
		if err != nil {
			return err
		}

		// Remove the sent request from the sender
		delete(fromConnections.SentRequests, toUID)

		// Remove the pending request from the recipient
		delete(toConnections.PendingRequests, fromUID)

		// Save updates for both users
		err = tx.Set(s.firestoreClient.Collection("user_connections").Doc(fromConnections.ID),
			mappers.MapConnectionsGoToFirestore(*fromConnections))
		if err != nil {
			return err
		}

		err = tx.Set(s.firestoreClient.Collection("user_connections").Doc(toConnections.ID),
			mappers.MapConnectionsGoToFirestore(*toConnections))
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *UserConnectionService) GetUserConnectionCount(ctx context.Context, uid string) (*int, error) {
	userConnections, error := s.GetUserConnections(ctx, uid)
	if error != nil {
		return nil, error
	}

	connections := len(userConnections.Connections)

	return &connections, nil
}
