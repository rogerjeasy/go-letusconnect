package test

import (
	"context"
	"errors"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/api/iterator"

	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"github.com/rogerjeasy/go-letusconnect/services"
)

func (m *MockFirestoreClient) Collection(name string) services.FirestoreClient {
	args := m.Called(name)
	return args.Get(0).(services.FirestoreClient)
}

type MockFirestoreCollection struct {
	mock.Mock
}

func (m *MockFirestoreCollection) Doc(docID string) services.FirestoreClient {
	args := m.Called(docID)
	return args.Get(0).(services.FirestoreClient)
}

func (m *MockFirestoreCollection) Where(path, op string, value interface{}) services.FirestoreClient {
	args := m.Called(path, op, value)
	return args.Get(0).(services.FirestoreClient)
}

type MockFirestoreDocument struct {
	mock.Mock
}

func (m *MockFirestoreDocument) Get(ctx context.Context) (*firestore.DocumentSnapshot, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).(*firestore.DocumentSnapshot), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockFirestoreDocument) Set(ctx context.Context, data interface{}, opts ...firestore.SetOption) (*firestore.WriteResult, error) {
	args := m.Called(ctx, data, opts)
	return args.Get(0).(*firestore.WriteResult), args.Error(1)
}

func TestUserConnectionService_CreateUserConnections(t *testing.T) {
	mockFirestore := new(MockFirestoreClient)
	mockCollection := new(MockFirestoreCollection)
	mockDocument := new(MockFirestoreDocument)

	mockFirestore.On("Collection", "user_connections").Return(mockCollection)
	mockCollection.On("Doc", mock.Anything).Return(mockDocument)
	mockDocument.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(&firestore.WriteResult{}, nil)

	userConnectionService := services.NewUserConnectionService(mockFirestore, nil)

	connections, err := userConnectionService.CreateUserConnections(context.Background(), "12345")
	assert.NoError(t, err)
	assert.NotNil(t, connections)
	assert.Equal(t, "12345", connections.UID)
}

func TestUserConnectionService_CheckUserConnectionsExist(t *testing.T) {
	mockFirestore := new(MockFirestoreClient)
	mockCollection := new(MockFirestoreCollection)
	mockQuery := new(MockFirestoreQuery)
	mockDocumentIterator := new(MockFirestoreDocument)

	mockFirestore.On("Collection", "user_connections").Return(mockCollection)
	mockCollection.On("Where", "uid", "==", "12345").Return(mockQuery)
	mockQuery.On("Documents", mock.Anything).Return(mockDocumentIterator)
	mockDocumentIterator.On("Next").Return(nil, iterator.Done)

	userConnectionService := services.NewUserConnectionService(mockFirestore, nil)

	exists, docID, err := userConnectionService.CheckUserConnectionsExist(context.Background(), "12345")
	assert.NoError(t, err)
	assert.False(t, exists)
	assert.Equal(t, "", docID)
}

func TestUserConnectionService_GetUserConnections(t *testing.T) {
	mockFirestore := new(MockFirestoreClient)
	mockCollection := new(MockFirestoreCollection)
	mockDocument := new(MockFirestoreDocument)

	data := map[string]interface{}{
		"uid":              "12345",
		"connections":      map[string]interface{}{},
		"pending_requests": map[string]interface{}{},
		"sent_requests":    map[string]interface{}{},
	}

	mockFirestore.On("Collection", "user_connections").Return(mockCollection)
	mockCollection.On("Where", "uid", "==", "12345").Return(mockCollection)
	mockCollection.On("Documents", mock.Anything).Return(mockDocument)
	mockDocument.On("Next").Return(data, nil)

	userConnectionService := services.NewUserConnectionService(mockFirestore, nil)

	connections, err := userConnectionService.GetUserConnections(context.Background(), "12345")
	assert.NoError(t, err)
	assert.NotNil(t, connections)
	assert.Equal(t, "12345", connections.UID)
}

func TestUserConnectionService_GetUserConnections_Fail(t *testing.T) {
	mockFirestore := new(MockFirestoreClient)
	mockCollection := new(MockFirestoreCollection)
	mockDocument := new(MockFirestoreDocument)

	mockFirestore.On("Collection", "user_connections").Return(mockCollection)
	mockCollection.On("Where", "uid", "==", "12345").Return(mockCollection)
	mockCollection.On("Documents", mock.Anything).Return(mockDocument)
	mockDocument.On("Next").Return(nil, errors.New("firestore error"))

	userConnectionService := services.NewUserConnectionService(mockFirestore, nil)

	_, err := userConnectionService.GetUserConnections(context.Background(), "12345")
	assert.Error(t, err)
	assert.EqualError(t, err, "failed to check user connections: firestore error")
}

func TestUserConnectionService_SendConnectionRequest(t *testing.T) {
	mockFirestoreClient := new(MockFirestoreClient)
	mockCollection := new(MockFirestoreCollection)
	mockDocument := new(MockFirestoreDocument)
	mockDocumentSnapshot := new(MockDocumentSnapshot)
	mockQuery := new(MockFirestoreClient)
	mockDocumentIterator := new(MockDocumentIterator)

	userService := new(MockUserService)
	userConnectionService := services.NewUserConnectionService(mockFirestoreClient, userService)

	ctx := context.Background()
	fromUID := "fromUID"
	toUID := "toUID"
	message := "Hello, let's connect!"

	// Helper function to mock GetUserConnections
	mockGetUserConnections := func(uid string, connections *models.UserConnections, err error) {
		mockFirestoreClient.On("Collection", "user_connections").Return(mockCollection)
		mockCollection.On("Where", "uid", "==", uid).Return(mockQuery)
		mockQuery.On("Documents", ctx).Return(mockDocumentIterator)
		mockDocumentIterator.On("Next").Return(mockDocumentSnapshot, nil)
		mockDocumentSnapshot.On("Ref").Return(&firestore.DocumentRef{ID: "docID"})
		mockCollection.On("Doc", "docID").Return(mockDocument)
		mockDocument.On("Get", ctx).Return(mockDocumentSnapshot, nil)
		mockDocumentSnapshot.On("Data").Return(mappers.MapConnectionsGoToFirestore(*connections))
	}

	t.Run("Successfully send connection request", testSendConnectionRequestSuccess(mockGetUserConnections, userService, userConnectionService, ctx, fromUID, toUID, message))
	t.Run("Failed to get sender's connections", testSendConnectionRequestSenderConnectionsError(mockFirestoreClient, mockCollection, mockQuery, mockDocumentIterator, userConnectionService, ctx, fromUID, toUID, message))
	t.Run("Failed to get recipient's connections", testSendConnectionRequestRecipientConnectionsError(mockGetUserConnections, mockFirestoreClient, mockCollection, mockQuery, mockDocumentIterator, userConnectionService, ctx, fromUID, toUID, message))
	t.Run("Users are already connected", testSendConnectionRequestAlreadyConnected(mockGetUserConnections, userConnectionService, ctx, fromUID, toUID, message))
	t.Run("Sender already has a pending request to the recipient", testSendConnectionRequestSenderPendingRequest(mockGetUserConnections, userConnectionService, ctx, fromUID, toUID, message))
	t.Run("Recipient already has a pending request to the sender", testSendConnectionRequestRecipientPendingRequest(mockGetUserConnections, userConnectionService, ctx, fromUID, toUID, message))
	t.Run("Firestore transaction fails", testSendConnectionRequestTransactionError(mockGetUserConnections, userService, mockFirestoreClient, userConnectionService, ctx, fromUID, toUID, message))
}

func testSendConnectionRequestSuccess(mockGetUserConnections func(string, *models.UserConnections, error), userService *MockUserService, userConnectionService *UserConnectionService, ctx context.Context, fromUID, toUID, message string) func(*testing.T) {
	return func(t *testing.T) {
		fromConnections := &models.UserConnections{
			ID:              "fromDocID",
			UID:             fromUID,
			Connections:     make(map[string]models.Connection),
			PendingRequests: make(map[string]models.ConnectionRequest),
			SentRequests:    make(map[string]models.SentRequest),
		}

		toConnections := &models.UserConnections{
			ID:              "toDocID",
			UID:             toUID,
			Connections:     make(map[string]models.Connection),
			PendingRequests: make(map[string]models.ConnectionRequest),
			SentRequests:    make(map[string]models.SentRequest),
		}

		mockGetUserConnections(fromUID, fromConnections, nil)
		mockGetUserConnections(toUID, toConnections, nil)

		userService.On("GetUsernameByUID", fromUID).Return("fromUsername", nil)
		mockFirestoreClient.On("RunTransaction", ctx, mock.Anything).Return(nil)

		err := userConnectionService.SendConnectionRequest(ctx, fromUID, toUID, message)
		assert.NoError(t, err)
	}
}

func testSendConnectionRequestSenderConnectionsError(mockFirestoreClient *MockFirestoreClient, mockCollection *MockFirestoreCollection, mockQuery *MockFirestoreClient, mockDocumentIterator *MockDocumentIterator, userConnectionService *UserConnectionService, ctx context.Context, fromUID, toUID, message string) func(*testing.T) {
	return func(t *testing.T) {
		mockFirestoreClient.On("Collection", "user_connections").Return(mockCollection)
		mockCollection.On("Where", "uid", "==", fromUID).Return(mockQuery)
		mockQuery.On("Documents", ctx).Return(mockDocumentIterator)
		mockDocumentIterator.On("Next").Return(nil, errors.New("firestore error"))

		err := userConnectionService.SendConnectionRequest(ctx, fromUID, toUID, message)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get sender's connections")
	}
}

func testSendConnectionRequestRecipientConnectionsError(mockGetUserConnections func(string, *models.UserConnections, error), mockFirestoreClient *MockFirestoreClient, mockCollection *MockFirestoreCollection, mockQuery *MockFirestoreQuery, mockDocumentIterator *MockDocumentIterator, userConnectionService *UserConnectionService, ctx context.Context, fromUID, toUID, message string) func(*testing.T) {
	return func(t *testing.T) {
		fromConnections := &models.UserConnections{
			ID:              "fromDocID",
			UID:             fromUID,
			Connections:     make(map[string]models.Connection),
			PendingRequests: make(map[string]models.ConnectionRequest),
			SentRequests:    make(map[string]models.SentRequest),
		}

		mockGetUserConnections(fromUID, fromConnections, nil)
		mockFirestoreClient.On("Collection", "user_connections").Return(mockCollection)
		mockCollection.On("Where", "uid", "==", toUID).Return(mockQuery)
		mockQuery.On("Documents", ctx).Return(mockDocumentIterator)
		mockDocumentIterator.On("Next").Return(nil, errors.New("firestore error"))

		err := userConnectionService.SendConnectionRequest(ctx, fromUID, toUID, message)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get recipient's connections")
	}
}

func testSendConnectionRequestAlreadyConnected(mockGetUserConnections func(string, *models.UserConnections, error), userConnectionService *UserConnectionService, ctx context.Context, fromUID, toUID, message string) func(*testing.T) {
	return func(t *testing.T) {
		fromConnections := &models.UserConnections{
			ID:  "fromDocID",
			UID: fromUID,
			Connections: map[string]models.Connection{
				toUID: {UID: toUID},
			},
			PendingRequests: make(map[string]models.ConnectionRequest),
			SentRequests:    make(map[string]models.SentRequest),
		}

		toConnections := &models.UserConnections{
			ID:              "toDocID",
			UID:             toUID,
			Connections:     make(map[string]models.Connection),
			PendingRequests: make(map[string]models.ConnectionRequest),
			SentRequests:    make(map[string]models.SentRequest),
		}

		mockGetUserConnections(fromUID, fromConnections, nil)
		mockGetUserConnections(toUID, toConnections, nil)

		err := userConnectionService.SendConnectionRequest(ctx, fromUID, toUID, message)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "you are already connected with this user")
	}
}

func testSendConnectionRequestSenderPendingRequest(mockGetUserConnections func(string, *models.UserConnections, error), userConnectionService *UserConnectionService, ctx context.Context, fromUID, toUID, message string) func(*testing.T) {
	return func(t *testing.T) {
		fromConnections := &models.UserConnections{
			ID:              "fromDocID",
			UID:             fromUID,
			Connections:     make(map[string]models.Connection),
			PendingRequests: make(map[string]models.ConnectionRequest),
			SentRequests: map[string]models.SentRequest{
				toUID: {ToUID: toUID, Status: "pending"},
			},
		}

		toConnections := &models.UserConnections{
			ID:              "toDocID",
			UID:             toUID,
			Connections:     make(map[string]models.Connection),
			PendingRequests: make(map[string]models.ConnectionRequest),
			SentRequests:    make(map[string]models.SentRequest),
		}

		mockGetUserConnections(fromUID, fromConnections, nil)
		mockGetUserConnections(toUID, toConnections, nil)

		err := userConnectionService.SendConnectionRequest(ctx, fromUID, toUID, message)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "you already have a pending connection request to this user")
	}
}

func testSendConnectionRequestRecipientPendingRequest(mockGetUserConnections func(string, *models.UserConnections, error), userConnectionService *UserConnectionService, ctx context.Context, fromUID, toUID, message string) func(*testing.T) {
	return func(t *testing.T) {
		fromConnections := &models.UserConnections{
			ID:          "fromDocID",
			UID:         fromUID,
			Connections: make(map[string]models.Connection),
			PendingRequests: map[string]models.ConnectionRequest{
				toUID: {FromUID: toUID, Status: "pending"},
			},
			SentRequests: make(map[string]models.SentRequest),
		}

		toConnections := &models.UserConnections{
			ID:              "toDocID",
			UID:             toUID,
			Connections:     make(map[string]models.Connection),
			PendingRequests: make(map[string]models.ConnectionRequest),
			SentRequests:    make(map[string]models.SentRequest),
		}

		mockGetUserConnections(fromUID, fromConnections, nil)
		mockGetUserConnections(toUID, toConnections, nil)

		err := userConnectionService.SendConnectionRequest(ctx, fromUID, toUID, message)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "this user has already sent you a connection request")
	}
}

func testSendConnectionRequestTransactionError(mockGetUserConnections func(string, *models.UserConnections, error), userService *MockUserService, mockFirestoreClient *MockFirestoreClient, userConnectionService *UserConnectionService, ctx context.Context, fromUID, toUID, message string) func(*testing.T) {
	return func(t *testing.T) {
		fromConnections := &models.UserConnections{
			ID:              "fromDocID",
			UID:             fromUID,
			Connections:     make(map[string]models.Connection),
			PendingRequests: make(map[string]models.ConnectionRequest),
			SentRequests:    make(map[string]models.SentRequest),
		}

		toConnections := &models.UserConnections{
			ID:              "toDocID",
			UID:             toUID,
			Connections:     make(map[string]models.Connection),
			PendingRequests: make(map[string]models.ConnectionRequest),
			SentRequests:    make(map[string]models.SentRequest),
		}

		mockGetUserConnections(fromUID, fromConnections, nil)
		mockGetUserConnections(toUID, toConnections, nil)

		userService.On("GetUsernameByUID", fromUID).Return("fromUsername", nil)
		mockFirestoreClient.On("RunTransaction", ctx, mock.Anything).Return(errors.New("transaction error"))

		err := userConnectionService.SendConnectionRequest(ctx, fromUID, toUID, message)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to process connection request")
	}
}

func TestUserConnectionService_AcceptConnectionRequest(t *testing.T) {
	mockFirestoreClient := new(MockFirestoreClient)
	mockCollection := new(MockFirestoreCollection)
	mockDocument := new(MockFirestoreDocument)
	mockDocumentSnapshot := new(MockDocumentSnapshot)
	mockQuery := new(MockFirestoreQuery)
	mockDocumentIterator := new(MockDocumentIterator)

	userService := new(MockUserService)
	userConnectionService := NewUserConnectionService(mockFirestoreClient, userService)

	ctx := context.Background()
	fromUID := "fromUID"
	toUID := "toUID"

	// Helper function to mock GetUserConnections
	mockGetUserConnections := func(uid string, connections *models.UserConnections, err error) {
		mockFirestoreClient.On("Collection", "user_connections").Return(mockCollection)
		mockCollection.On("Where", "uid", "==", uid).Return(mockQuery)
		mockQuery.On("Documents", ctx).Return(mockDocumentIterator)
		mockDocumentIterator.On("Next").Return(mockDocumentSnapshot, nil)
		mockDocumentSnapshot.On("Ref").Return(&firestore.DocumentRef{ID: "docID"})
		mockCollection.On("Doc", "docID").Return(mockDocument)
		mockDocument.On("Get", ctx).Return(mockDocumentSnapshot, nil)
		mockDocumentSnapshot.On("Data").Return(mappers.MapConnectionsGoToFirestore(*connections))
	}

	t.Run("Successfully accept connection request", testAcceptConnectionRequestSuccess(mockGetUserConnections, userService, mockFirestoreClient, userConnectionService, ctx, fromUID, toUID))
	t.Run("Failed to get sender's connections", testAcceptConnectionRequestSenderConnectionsError(mockFirestoreClient, mockCollection, mockQuery, mockDocumentIterator, userConnectionService, ctx, fromUID, toUID))
	t.Run("Failed to get recipient's connections", testAcceptConnectionRequestRecipientConnectionsError(mockGetUserConnections, mockFirestoreClient, mockCollection, mockQuery, mockDocumentIterator, userConnectionService, ctx, fromUID, toUID))
	t.Run("Failed to get recipient's username", testAcceptConnectionRequestRecipientUsernameError(mockGetUserConnections, userService, userConnectionService, ctx, fromUID, toUID))
	t.Run("Failed to get sender's username", testAcceptConnectionRequestSenderUsernameError(mockGetUserConnections, userService, userConnectionService, ctx, fromUID, toUID))
	t.Run("Firestore transaction fails", testAcceptConnectionRequestTransactionError(mockGetUserConnections, userService, mockFirestoreClient, userConnectionService, ctx, fromUID, toUID))
}
