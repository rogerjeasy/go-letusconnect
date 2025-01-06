package services

// import (
//     "context"
//     "testing"
//     "cloud.google.com/go/firestore"
//     "github.com/rogerjeasy/go-letusconnect/models"
//     "github.com/stretchr/testify/assert"
//     "github.com/stretchr/testify/mock"
//     "google.golang.org/api/iterator"

// )

// func setupMockAddressService() *AddressService {
// 	mockClient := &MockFirestoreClient{
// 		addresses: make(map[string]map[string]interface{}),
// 	}
// 	return &AddressService{
// 		FirestoreClient: mockClient,
// 	}
// }

// // MockFirestoreClient implements necessary firestore.Client methods
// type MockFirestoreClient struct {
//     mock.Mock
// }

// // MockCollectionRef mocks firestore.CollectionRef
// type MockCollectionRef struct {
//     mock.Mock
// }

// // MockDocumentRef mocks firestore.DocumentRef
// type MockDocumentRef struct {
//     mock.Mock
//     id string
// }

// // MockDocumentIterator mocks firestore.DocumentIterator
// type MockDocumentIterator struct {
//     mock.Mock
//     docs    []map[string]interface{}
//     currIdx int
// }

// // MockQuery implementation
// type MockQuery struct {
//     mock.Mock
// }

// // Implement necessary methods for MockFirestoreClient
// func (m *MockFirestoreClient) Collection(path string) *firestore.CollectionRef {
//     args := m.Called(path)
//     return args.Get(0).(*firestore.CollectionRef)
// }

// func (m *MockFirestoreClient) Close() error {
//     return nil
// }

// // Implement necessary methods for MockCollectionRef
// func (m *MockCollectionRef) Doc(path string) *firestore.DocumentRef {
//     args := m.Called(path)
//     return args.Get(0).(*firestore.DocumentRef)
// }

// func (m *MockCollectionRef) Add(ctx context.Context, data interface{}) (*firestore.DocumentRef, *firestore.WriteResult, error) {
//     args := m.Called(ctx, data)
//     return args.Get(0).(*firestore.DocumentRef), nil, args.Error(2)
// }

// func (m *MockCollectionRef) Where(path, op string, value interface{}) *firestore.Query {
//     args := m.Called(path, op, value)
//     return args.Get(0).(*firestore.Query)
// }

// // Implement necessary methods for MockDocumentRef
// func (m *MockDocumentRef) Set(ctx context.Context, data interface{}, opts ...firestore.SetOption) (*firestore.WriteResult, error) {
//     args := m.Called(ctx, data, opts)
//     return nil, args.Error(0)
// }

// func (m *MockDocumentRef) Get(ctx context.Context) (*firestore.DocumentSnapshot, error) {
//     args := m.Called(ctx)
//     return args.Get(0).(*firestore.DocumentSnapshot), args.Error(1)
// }

// func (m *MockQuery) Documents(ctx context.Context) *firestore.DocumentIterator {
//     args := m.Called(ctx)
//     return args.Get(0).(*firestore.DocumentIterator)
// }

// // MockDocumentIterator implementation
// func (m *MockDocumentIterator) Next() (*firestore.DocumentSnapshot, error) {
//     if m.currIdx >= len(m.docs) {
//         return nil, iterator.Done
//     }

//     doc := &firestore.DocumentSnapshot{
//         Ref: &firestore.DocumentRef{},
//     }
//     m.currIdx++
//     return doc, nil
// }

// func (m *MockDocumentIterator) Stop() {}

// // Test cases
// func TestCreateUserAddress_Success(t *testing.T) {
//     // Setup
//     mockClient := new(MockFirestoreClient)
//     mockCollRef := new(MockCollectionRef)
//     mockDocRef := new(MockDocumentRef)

//     service := &AddressService{
//         FirestoreClient: &firestore.Client{},  // Initialize with empty client
//     }
//     // Replace the client with our mock
//     service.FirestoreClient = mockClient

//     uid := "test-uid"

//     // Set up expectations
//     mockClient.On("Collection", "user_addresses").Return(mockCollRef)
//     mockCollRef.On("Add", mock.Anything, mock.Anything).Return(mockDocRef, nil, nil)
//     mockDocRef.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil)

//     // Execute
//     result, err := service.CreateUserAddress(uid)

//     // Assert
//     assert.NoError(t, err)
//     assert.Equal(t, uid, result.UID)
//     mockClient.AssertExpectations(t)
//     mockCollRef.AssertExpectations(t)
// }

// func TestGetUserAddresses_Success(t *testing.T) {
// 	// Setup
// 	mockClient := newMockFirestoreClient(t)
// 	mockCollRef := &MockCollectionRef{}
// 	mockQuery := &firestore.Query{}
// 	// mockIter := &MockDocumentIterator{
// 	// 	docs: []map[string]interface{}{
// 	// 		{
// 	// 			"id":  "test-id",
// 	// 			"uid": "test-uid",
// 	// 		},
// 	// 	},
// 	// }

// 	service := &AddressService{
// 		FirestoreClient: mockClient.Client,
// 	}

// 	uid := "test-uid"

// 	// Expectations
// 	mockClient.On("Collection", "user_addresses").Return(mockCollRef)
// 	mockCollRef.On("Where", "uid", "==", uid).Return(mockQuery)

// 	// Execute
// 	addresses, err := service.GetUserAddresses(uid)

// 	// Assert
// 	assert.NoError(t, err)
// 	assert.Len(t, addresses, 1)
// 	mockClient.AssertExpectations(t)
// 	mockCollRef.AssertExpectations(t)
// }

// func TestUpdateUserAddress_Success(t *testing.T) {
// 	// Setup
// 	mockClient := newMockFirestoreClient(t)
// 	mockCollRef := &MockCollectionRef{}
// 	mockDocRef := &MockDocumentRef{}
// 	mockSnapshot := &firestore.DocumentSnapshot{}

// 	service := &AddressService{
// 		FirestoreClient: mockClient.Client,
// 	}

// 	addressID := "test-id"
// 	uid := "test-uid"
// 	updatedAddress := models.UserAddress{
// 		ID:     addressID,
// 		UID:    uid,
// 		Street: "Updated Street",
// 	}

// 	// Setup mock document data
// 	// mockData := map[string]interface{}{
// 	// 	"uid": uid,
// 	// }

// 	// Expectations
// 	mockClient.On("Collection", "user_addresses").Return(mockCollRef)
// 	mockCollRef.On("Doc", addressID).Return(mockDocRef)
// 	mockDocRef.On("Get", mock.Anything).Return(mockSnapshot, nil)
// 	mockDocRef.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil)

// 	// Execute
// 	result, err := service.UpdateUserAddress(addressID, uid, updatedAddress)

// 	// Assert
// 	assert.NoError(t, err)
// 	assert.Equal(t, updatedAddress.Street, result.Street)
// 	mockClient.AssertExpectations(t)
// 	mockCollRef.AssertExpectations(t)
// 	mockDocRef.AssertExpectations(t)
// }

// func TestCreateUserAddress_NilClient(t *testing.T) {
// 	// Setup
// 	service := &AddressService{
// 		FirestoreClient: nil,
// 	}

// 	// Execute
// 	_, err := service.CreateUserAddress("test-uid")

// 	// Assert
// 	assert.Error(t, err)
// 	assert.Contains(t, err.Error(), "firestore client is not initialized")
// }

// func TestUpdateUserAddress_UnauthorizedUpdate(t *testing.T) {
// 	// Setup
// 	mockClient := newMockFirestoreClient(t)
// 	mockCollRef := &MockCollectionRef{}
// 	mockDocRef := &MockDocumentRef{}
// 	mockSnapshot := &firestore.DocumentSnapshot{}

// 	service := &AddressService{
// 		FirestoreClient: mockClient.Client,
// 	}

// 	addressID := "test-id"
// 	uid := "test-uid"
// 	wrongUID := "wrong-uid"
// 	updatedAddress := models.UserAddress{
// 		ID:     addressID,
// 		UID:    wrongUID,
// 		Street: "Updated Street",
// 	}

// 	// Expectations
// 	mockClient.On("Collection", "user_addresses").Return(mockCollRef)
// 	mockCollRef.On("Doc", addressID).Return(mockDocRef)
// 	mockDocRef.On("Get", mock.Anything).Return(mockSnapshot, nil)

// 	// Execute
// 	_, err := service.UpdateUserAddress(addressID, uid, updatedAddress)

// 	// Assert
// 	assert.Error(t, err)
// 	assert.Contains(t, err.Error(), "unauthorized to update this address")
// }

// // MockDocumentIterator implementation
// func (m *MockDocumentIterator) Next() (*firestore.DocumentSnapshot, error) {
// 	if m.currIdx >= len(m.docs) {
// 		return nil, iterator.Done
// 	}

// 	doc := &firestore.DocumentSnapshot{
// 		Ref: &firestore.DocumentRef{},
// 	}
// 	m.currIdx++
// 	return doc, nil
// }

// func (m *MockDocumentIterator) Stop() {}
