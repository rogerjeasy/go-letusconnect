package services

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FAQService struct {
	firestoreClient *firestore.Client
}

func NewFAQService(client *firestore.Client) *FAQService {
	return &FAQService{
		firestoreClient: client,
	}
}

// CreateFAQ creates a new FAQ document in Firestore
func (s *FAQService) CreateFAQ(ctx context.Context, faq models.FAQ, username string, uid string) (*models.FAQ, error) {
	// Validate required fields
	if faq.Question == "" {
		return nil, fmt.Errorf("question is required")
	}
	if faq.Response == "" {
		return nil, fmt.Errorf("response is required")
	}
	if faq.Category == "" {
		return nil, fmt.Errorf("category is required")
	}

	// Set metadata fields
	now := time.Now()
	faq.CreatedAt = now
	faq.UpdatedAt = now
	faq.CreatedBy = uid
	faq.Username = username

	// Convert to Firestore format
	firestoreData := mappers.MapFAQGoToFirestore(faq)

	// Add to Firestore
	docRef, _, err := s.firestoreClient.Collection("faqs").Add(ctx, firestoreData)
	if err != nil {
		return nil, fmt.Errorf("failed to add FAQ to Firestore: %v", err)
	}

	// Set the document ID
	faq.ID = docRef.ID
	return &faq, nil
}

// GetFAQByID retrieves a single FAQ by its ID
func (s *FAQService) GetFAQByID(ctx context.Context, id string) (*models.FAQ, error) {
	docRef := s.firestoreClient.Collection("faqs").Doc(id)
	doc, err := docRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("FAQ not found")
		}
		return nil, fmt.Errorf("failed to get FAQ: %v", err)
	}

	// Get the Firestore data and add the document ID
	data := doc.Data()
	data["id"] = doc.Ref.ID

	// Convert to FAQ struct
	faq := mappers.MapFAQFirestoreToGo(data)
	return &faq, nil
}
