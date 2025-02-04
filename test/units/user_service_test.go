package test

import (
	"context"
	"errors"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/rogerjeasy/go-letusconnect/models"
	"github.com/rogerjeasy/go-letusconnect/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/api/iterator"
)

// MockFirestoreClient is a mock implementation of the Firestore client
type MockFirestoreClient struct {
	mock.Mock
}

func (m *MockFirestoreClient) Collection(path string) *firestore.CollectionRef {
	args := m.Called(path)
	return args.Get(0).(*firestore.CollectionRef)
}

// MockCollectionRef is a mock implementation of the Firestore CollectionRef
type MockCollectionRef struct {
	mock.Mock
}

func (m *MockCollectionRef) Where(path, op string, value interface{}) *firestore.Query {
	args := m.Called(path, op, value)
	return args.Get(0).(*firestore.Query)
}

func (m *MockCollectionRef) Documents(ctx context.Context) *firestore.DocumentIterator {
	args := m.Called(ctx)
	return args.Get(0).(*firestore.DocumentIterator)
}

// MockDocumentIterator is a mock implementation of the Firestore DocumentIterator
type MockDocumentIterator struct {
	mock.Mock
}

func (m *MockDocumentIterator) Next() (*firestore.DocumentSnapshot, error) {
	args := m.Called()
	return args.Get(0).(*firestore.DocumentSnapshot), args.Error(1)
}

func (m *MockDocumentIterator) Stop() {
	m.Called()
}

// MockDocumentSnapshot is a mock implementation of the Firestore DocumentSnapshot
type MockDocumentSnapshot struct {
	mock.Mock
}

func (m *MockDocumentSnapshot) Data() map[string]interface{} {
	args := m.Called()
	return args.Get(0).(map[string]interface{})
}

func TestUserService_GetUserByEmail(t *testing.T) {
	mockFirestoreClient := new(MockFirestoreClient)
	mockCollectionRef := new(MockCollectionRef)
	mockDocumentIterator := new(MockDocumentIterator)
	mockDocumentSnapshot := new(MockDocumentSnapshot)

	userService := services.NewUserService(mockFirestoreClient)

	// Test case: User found
	t.Run("User found", func(t *testing.T) {
		expectedUser := &models.User{
			Email: "test@example.com",
		}

		mockFirestoreClient.On("Collection", "users").Return(mockCollectionRef)
		mockCollectionRef.On("Where", "email", "==", "test@example.com").Return(&firestore.Query{})
		mockCollectionRef.On("Documents", mock.Anything).Return(mockDocumentIterator)
		mockDocumentIterator.On("Next").Return(mockDocumentSnapshot, nil)
		mockDocumentSnapshot.On("Data").Return(map[string]interface{}{
			"email": "test@example.com",
		})

		user, err := userService.GetUserByEmail("test@example.com")
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.Email, user.Email)
	})

	// Test case: User not found
	t.Run("User not found", func(t *testing.T) {
		mockFirestoreClient.On("Collection", "users").Return(mockCollectionRef)
		mockCollectionRef.On("Where", "email", "==", "notfound@example.com").Return(&firestore.Query{})
		mockCollectionRef.On("Documents", mock.Anything).Return(mockDocumentIterator)
		mockDocumentIterator.On("Next").Return(nil, iterator.Done)

		_, err := userService.GetUserByEmail("notfound@example.com")
		assert.EqualError(t, err, "user not found")
	})

	// Test case: Firestore error
	t.Run("Firestore error", func(t *testing.T) {
		mockFirestoreClient.On("Collection", "users").Return(mockCollectionRef)
		mockCollectionRef.On("Where", "email", "==", "error@example.com").Return(&firestore.Query{})
		mockCollectionRef.On("Documents", mock.Anything).Return(mockDocumentIterator)
		mockDocumentIterator.On("Next").Return(nil, errors.New("firestore error"))

		_, err := userService.GetUserByEmail("error@example.com")
		assert.EqualError(t, err, "failed to fetch user data")
	})
}

func TestUserService_GetUserByUsername(t *testing.T) {
	mockFirestoreClient := new(MockFirestoreClient)
	mockCollectionRef := new(MockCollectionRef)
	mockDocumentIterator := new(MockDocumentIterator)
	mockDocumentSnapshot := new(MockDocumentSnapshot)

	userService := NewUserService(mockFirestoreClient)

	// Test case: User found
	t.Run("User found", func(t *testing.T) {
		expectedUser := &models.User{
			Username: "testuser",
		}

		mockFirestoreClient.On("Collection", "users").Return(mockCollectionRef)
		mockCollectionRef.On("Where", "username", "==", "testuser").Return(&firestore.Query{})
		mockCollectionRef.On("Documents", mock.Anything).Return(mockDocumentIterator)
		mockDocumentIterator.On("Next").Return(mockDocumentSnapshot, nil)
		mockDocumentSnapshot.On("Data").Return(map[string]interface{}{
			"username": "testuser",
		})

		user, err := userService.GetUserByUsername("testuser")
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.Username, user.Username)
	})

	// Test case: User not found
	t.Run("User not found", func(t *testing.T) {
		mockFirestoreClient.On("Collection", "users").Return(mockCollectionRef)
		mockCollectionRef.On("Where", "username", "==", "notfound").Return(&firestore.Query{})
		mockCollectionRef.On("Documents", mock.Anything).Return(mockDocumentIterator)
		mockDocumentIterator.On("Next").Return(nil, iterator.Done)

		_, err := userService.GetUserByUsername("notfound")
		assert.EqualError(t, err, "user not found")
	})

	// Test case: Firestore error
	t.Run("Firestore error", func(t *testing.T) {
		mockFirestoreClient.On("Collection", "users").Return(mockCollectionRef)
		mockCollectionRef.On("Where", "username", "==", "error").Return(&firestore.Query{})
		mockCollectionRef.On("Documents", mock.Anything).Return(mockDocumentIterator)
		mockDocumentIterator.On("Next").Return(nil, errors.New("firestore error"))

		_, err := userService.GetUserByUsername("error")
		assert.EqualError(t, err, "failed to fetch user data")
	})
}

func TestUserService_GetUserRole(t *testing.T) {
	mockFirestoreClient := new(MockFirestoreClient)
	mockCollectionRef := new(MockCollectionRef)
	mockDocumentIterator := new(MockDocumentIterator)
	mockDocumentSnapshot := new(MockDocumentSnapshot)

	userService := NewUserService(mockFirestoreClient)

	// Test case: Role found
	t.Run("Role found", func(t *testing.T) {
		expectedRoles := []string{"admin", "user"}

		mockFirestoreClient.On("Collection", "users").Return(mockCollectionRef)
		mockCollectionRef.On("Where", "uid", "==", "123").Return(&firestore.Query{})
		mockCollectionRef.On("Documents", mock.Anything).Return(mockDocumentIterator)
		mockDocumentIterator.On("Next").Return(mockDocumentSnapshot, nil)
		mockDocumentSnapshot.On("Data").Return(map[string]interface{}{
			"role": []interface{}{"admin", "user"},
		})

		roles, err := userService.GetUserRole("123")
		assert.NoError(t, err)
		assert.Equal(t, expectedRoles, roles)
	})

	// Test case: Role not found
	t.Run("Role not found", func(t *testing.T) {
		mockFirestoreClient.On("Collection", "users").Return(mockCollectionRef)
		mockCollectionRef.On("Where", "uid", "==", "456").Return(&firestore.Query{})
		mockCollectionRef.On("Documents", mock.Anything).Return(mockDocumentIterator)
		mockDocumentIterator.On("Next").Return(mockDocumentSnapshot, nil)
		mockDocumentSnapshot.On("Data").Return(map[string]interface{}{})

		_, err := userService.GetUserRole("456")
		assert.EqualError(t, err, "role not found for the user")
	})

	// Test case: Invalid role format
	t.Run("Invalid role format", func(t *testing.T) {
		mockFirestoreClient.On("Collection", "users").Return(mockCollectionRef)
		mockCollectionRef.On("Where", "uid", "==", "789").Return(&firestore.Query{})
		mockCollectionRef.On("Documents", mock.Anything).Return(mockDocumentIterator)
		mockDocumentIterator.On("Next").Return(mockDocumentSnapshot, nil)
		mockDocumentSnapshot.On("Data").Return(map[string]interface{}{
			"role": "admin",
		})

		_, err := userService.GetUserRole("789")
		assert.EqualError(t, err, "invalid role format")
	})
}

func TestUserService_GetUsernameByUID(t *testing.T) {
	mockFirestoreClient := new(MockFirestoreClient)
	mockCollectionRef := new(MockCollectionRef)
	mockDocumentIterator := new(MockDocumentIterator)
	mockDocumentSnapshot := new(MockDocumentSnapshot)

	userService := NewUserService(mockFirestoreClient)

	// Test case: Username found
	t.Run("Username found", func(t *testing.T) {
		expectedUsername := "testuser"

		mockFirestoreClient.On("Collection", "users").Return(mockCollectionRef)
		mockCollectionRef.On("Where", "uid", "==", "123").Return(&firestore.Query{})
		mockCollectionRef.On("Documents", mock.Anything).Return(mockDocumentIterator)
		mockDocumentIterator.On("Next").Return(mockDocumentSnapshot, nil)
		mockDocumentSnapshot.On("Data").Return(map[string]interface{}{
			"username": "testuser",
		})

		username, err := userService.GetUsernameByUID("123")
		assert.NoError(t, err)
		assert.Equal(t, expectedUsername, username)
	})

	// Test case: Username not found
	t.Run("Username not found", func(t *testing.T) {
		mockFirestoreClient.On("Collection", "users").Return(mockCollectionRef)
		mockCollectionRef.On("Where", "uid", "==", "456").Return(&firestore.Query{})
		mockCollectionRef.On("Documents", mock.Anything).Return(mockDocumentIterator)
		mockDocumentIterator.On("Next").Return(mockDocumentSnapshot, nil)
		mockDocumentSnapshot.On("Data").Return(map[string]interface{}{})

		_, err := userService.GetUsernameByUID("456")
		assert.EqualError(t, err, "username not found for the user")
	})

	// Test case: Invalid username format
	t.Run("Invalid username format", func(t *testing.T) {
		mockFirestoreClient.On("Collection", "users").Return(mockCollectionRef)
		mockCollectionRef.On("Where", "uid", "==", "789").Return(&firestore.Query{})
		mockCollectionRef.On("Documents", mock.Anything).Return(mockDocumentIterator)
		mockDocumentIterator.On("Next").Return(mockDocumentSnapshot, nil)
		mockDocumentSnapshot.On("Data").Return(map[string]interface{}{
			"username": 123,
		})

		_, err := userService.GetUsernameByUID("789")
		assert.EqualError(t, err, "invalid username format")
	})
}

func TestUserService_GetUserByUID(t *testing.T) {
	mockFirestoreClient := new(MockFirestoreClient)
	mockCollectionRef := new(MockCollectionRef)
	mockDocumentIterator := new(MockDocumentIterator)
	mockDocumentSnapshot := new(MockDocumentSnapshot)

	userService := NewUserService(mockFirestoreClient)

	// Test case: User found
	t.Run("User found", func(t *testing.T) {
		expectedData := map[string]interface{}{
			"uid": "123",
		}

		mockFirestoreClient.On("Collection", "users").Return(mockCollectionRef)
		mockCollectionRef.On("Where", "uid", "==", "123").Return(&firestore.Query{})
		mockCollectionRef.On("Documents", mock.Anything).Return(mockDocumentIterator)
		mockDocumentIterator.On("Next").Return(mockDocumentSnapshot, nil)
		mockDocumentSnapshot.On("Data").Return(expectedData)

		data, err := userService.GetUserByUID("123")
		assert.NoError(t, err)
		assert.Equal(t, expectedData, data)
	})

	// Test case: User not found
	t.Run("User not found", func(t *testing.T) {
		mockFirestoreClient.On("Collection", "users").Return(mockCollectionRef)
		mockCollectionRef.On("Where", "uid", "==", "456").Return(&firestore.Query{})
		mockCollectionRef.On("Documents", mock.Anything).Return(mockDocumentIterator)
		mockDocumentIterator.On("Next").Return(nil, iterator.Done)

		_, err := userService.GetUserByUID("456")
		assert.EqualError(t, err, "user not found")
	})

	// Test case: Firestore error
	t.Run("Firestore error", func(t *testing.T) {
		mockFirestoreClient.On("Collection", "users").Return(mockCollectionRef)
		mockCollectionRef.On("Where", "uid", "==", "789").Return(&firestore.Query{})
		mockCollectionRef.On("Documents", mock.Anything).Return(mockDocumentIterator)
		mockDocumentIterator.On("Next").Return(nil, errors.New("firestore error"))

		_, err := userService.GetUserByUID("789")
		assert.EqualError(t, err, "failed to fetch user data")
	})
}

func TestUserService_GetUserByUIDinGoStruct(t *testing.T) {
	mockFirestoreClient := new(MockFirestoreClient)
	mockCollectionRef := new(MockCollectionRef)
	mockDocumentIterator := new(MockDocumentIterator)
	mockDocumentSnapshot := new(MockDocumentSnapshot)

	userService := NewUserService(mockFirestoreClient)

	// Test case: User found
	t.Run("User found", func(t *testing.T) {
		expectedUser := &models.User{
			UID: "123",
		}

		mockFirestoreClient.On("Collection", "users").Return(mockCollectionRef)
		mockCollectionRef.On("Where", "uid", "==", "123").Return(&firestore.Query{})
		mockCollectionRef.On("Documents", mock.Anything).Return(mockDocumentIterator)
		mockDocumentIterator.On("Next").Return(mockDocumentSnapshot, nil)
		mockDocumentSnapshot.On("Data").Return(map[string]interface{}{
			"uid": "123",
		})

		user, err := userService.GetUserByUIDinGoStruct("123")
		assert.NoError(t, err)
		assert.Equal(t, expectedUser.UID, user.UID)
	})

	// Test case: User not found
	t.Run("User not found", func(t *testing.T) {
		mockFirestoreClient.On("Collection", "users").Return(mockCollectionRef)
		mockCollectionRef.On("Where", "uid", "==", "456").Return(&firestore.Query{})
		mockCollectionRef.On("Documents", mock.Anything).Return(mockDocumentIterator)
		mockDocumentIterator.On("Next").Return(nil, iterator.Done)

		_, err := userService.GetUserByUIDinGoStruct("456")
		assert.EqualError(t, err, "failed to fetch user data: user not found")
	})

	// Test case: Firestore error
	t.Run("Firestore error", func(t *testing.T) {
		mockFirestoreClient.On("Collection", "users").Return(mockCollectionRef)
		mockCollectionRef.On("Where", "uid", "==", "789").Return(&firestore.Query{})
		mockCollectionRef.On("Documents", mock.Anything).Return(mockDocumentIterator)
		mockDocumentIterator.On("Next").Return(nil, errors.New("firestore error"))

		_, err := userService.GetUserByUIDinGoStruct("789")
		assert.EqualError(t, err, "failed to fetch user data: firestore error")
	})
}
