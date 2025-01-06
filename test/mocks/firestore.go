package mocks

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/mock"
)

type FirestoreClient interface {
	Collection(name string) FirestoreCollection
}

type FirestoreCollection interface {
	Doc(docID string) FirestoreDocument
	Add(ctx context.Context, data interface{}) (*firestore.DocumentRef, *firestore.WriteResult, error)
	Where(path, op string, value interface{}) FirestoreQuery
}

type FirestoreDocument interface {
	Get(ctx context.Context) (*firestore.DocumentSnapshot, error)
	Set(ctx context.Context, data interface{}, opts ...firestore.SetOption) (*firestore.WriteResult, error)
}

type FirestoreQuery interface {
	Documents(ctx context.Context) *firestore.DocumentIterator
}

type MockFirestoreClient struct {
	mock.Mock
}

func (m *MockFirestoreClient) Collection(name string) FirestoreCollection {
	args := m.Called(name)
	return args.Get(0).(FirestoreCollection)
}

type MockFirestoreCollection struct {
	mock.Mock
}

func (m *MockFirestoreCollection) Doc(docID string) FirestoreDocument {
	args := m.Called(docID)
	return args.Get(0).(FirestoreDocument)
}

func (m *MockFirestoreCollection) Add(ctx context.Context, data interface{}) (*firestore.DocumentRef, *firestore.WriteResult, error) {
	args := m.Called(ctx, data)
	return args.Get(0).(*firestore.DocumentRef), args.Get(1).(*firestore.WriteResult), args.Error(2)
}

func (m *MockFirestoreCollection) Where(path, op string, value interface{}) FirestoreQuery {
	args := m.Called(path, op, value)
	return args.Get(0).(FirestoreQuery)
}

type MockFirestoreDocument struct {
	mock.Mock
}

func (m *MockFirestoreDocument) Get(ctx context.Context) (*firestore.DocumentSnapshot, error) {
	args := m.Called(ctx)
	return args.Get(0).(*firestore.DocumentSnapshot), args.Error(1)
}

func (m *MockFirestoreDocument) Set(ctx context.Context, data interface{}, opts ...firestore.SetOption) (*firestore.WriteResult, error) {
	args := m.Called(ctx, data, opts)
	return args.Get(0).(*firestore.WriteResult), args.Error(1)
}

type MockFirestoreQuery struct {
	mock.Mock
}

func (m *MockFirestoreQuery) Documents(ctx context.Context) *firestore.DocumentIterator {
	args := m.Called(ctx)
	return args.Get(0).(*firestore.DocumentIterator)
}
