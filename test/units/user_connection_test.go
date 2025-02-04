package test

import (
	"context"
	"errors"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/api/iterator"

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
