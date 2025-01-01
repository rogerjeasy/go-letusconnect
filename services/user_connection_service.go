package services

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserConnectionService struct {
	firestoreClient *firestore.Client
}

func NewUserConnectionService(client *firestore.Client) *UserConnectionService {
	return &UserConnectionService{
		firestoreClient: client,
	}
}

func (s *UserConnectionService) CreateUserConnections(ctx context.Context, uid string) (*models.UserConnections, error) {
	connections := models.UserConnections{
		ID:              uuid.New().String(),
		UID:             uid,
		Connections:     make(map[string]models.Connection),
		PendingRequests: make(map[string]models.ConnectionRequest),
	}

	_, err := s.firestoreClient.Collection("user_connections").Doc(connections.ID).Set(ctx, mappers.MapConnectionsGoToFirestore(connections))
	if err != nil {
		return nil, fmt.Errorf("failed to create user connections: %v", err)
	}

	return &connections, nil
}

func (s *UserConnectionService) GetUserConnections(ctx context.Context, uid string) (*models.UserConnections, error) {
	query := s.firestoreClient.Collection("user_connections").Where("uid", "==", uid).Limit(1)
	iter := query.Documents(ctx)
	doc, err := iter.Next()

	if err != nil {
		if status.Code(err) == codes.NotFound {
			// Create new connections if none exist
			return s.CreateUserConnections(ctx, uid)
		}
		return nil, fmt.Errorf("failed to get user connections: %v", err)
	}

	connections := mappers.MapConnectionsFirestoreToGo(doc.Data())
	return &connections, nil
}

func (s *UserConnectionService) SendConnectionRequest(ctx context.Context, fromUID, toUID, message string) error {
	request := models.ConnectionRequest{
		FromUID: fromUID,
		ToUID:   toUID,
		Message: message,
		SentAt:  time.Now(),
		Status:  "pending",
	}

	targetConnections, err := s.GetUserConnections(ctx, toUID)
	if err != nil {
		targetConnections, err = s.CreateUserConnections(ctx, toUID)
		if err != nil {
			return err
		}
	}

	targetConnections.PendingRequests[fromUID] = request

	_, err = s.firestoreClient.Collection("user_connections").Doc(targetConnections.ID).Set(ctx,
		mappers.MapConnectionsGoToFirestore(*targetConnections))
	if err != nil {
		return err
	}

	fromUsername, err := GetUsernameByUID(fromUID)
	if err != nil {
		return fmt.Errorf("failed to get username: %v", err)
	}
	if err := SendConnectionRequestNotification(ctx, fromUID, fromUsername, message, toUID); err != nil {
		fmt.Printf("Failed to send connection request notification: %v\n", err)
	}
	return nil
}

func (s *UserConnectionService) AcceptConnectionRequest(ctx context.Context, fromUID, toUID string) error {
	// Start a transaction
	err := s.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		// Get both users' connections
		fromConnections, err := s.GetUserConnections(ctx, fromUID)
		if err != nil {
			return err
		}
		toConnections, err := s.GetUserConnections(ctx, toUID)
		if err != nil {
			return err
		}

		// Create connection objects
		now := time.Now()
		connection1 := models.Connection{
			TargetUID:  toUID,
			SentAt:     now,
			AcceptedAt: now,
			Status:     "active",
		}
		connection2 := models.Connection{
			TargetUID:  fromUID,
			SentAt:     now,
			AcceptedAt: now,
			Status:     "active",
		}

		// Update connections and remove request
		fromConnections.Connections[toUID] = connection1
		toConnections.Connections[fromUID] = connection2
		delete(toConnections.PendingRequests, fromUID)

		// Update both documents in transaction
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

	return err
}

func (s *UserConnectionService) RejectConnectionRequest(ctx context.Context, fromUID, toUID string) error {
	targetConnections, err := s.GetUserConnections(ctx, toUID)
	if err != nil {
		return err
	}

	delete(targetConnections.PendingRequests, fromUID)

	_, err = s.firestoreClient.Collection("user_connections").Doc(targetConnections.ID).Set(ctx,
		mappers.MapConnectionsGoToFirestore(*targetConnections))
	return err
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
