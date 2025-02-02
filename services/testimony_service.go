package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TestimonialService struct {
	firestoreClient FirestoreClient
	userService     *UserService
}

func NewTestimonialService(fClient FirestoreClient, uService *UserService) *TestimonialService {
	return &TestimonialService{
		firestoreClient: fClient,
		userService:     uService,
	}
}

func (s *TestimonialService) CreateTestimonial(ctx context.Context, input models.Testimonial, userID string) (*models.Testimonial, error) {
	if input.Title == "" {
		return nil, errors.New("testimonial title is required")
	}

	if input.Content == "" {
		return nil, errors.New("testimonial content is required")
	}

	if input.ID == "" {
		input.ID = uuid.New().String()
	}

	input.UserID = userID
	input.CreatedAt = time.Now()
	input.UpdatedAt = time.Now()
	input.IsPublished = false
	input.Likes = 0

	testimonialData := mappers.MapTestimonialGoToFirestore(input)

	_, err := s.firestoreClient.Collection("testimonials").Doc(input.ID).Set(ctx, testimonialData)
	if err != nil {
		return nil, fmt.Errorf("failed to create testimonial: %v", err)
	}

	return &input, nil
}

// CreateAlumniTestimonial creates a new alumni testimonial
func (s *TestimonialService) CreateAlumniTestimonial(ctx context.Context, input models.AlumniTestimonial, userID string) (*models.AlumniTestimonial, error) {
	if input.GraduationYear == 0 {
		return nil, errors.New("graduation year is required")
	}

	testimonial, err := s.CreateTestimonial(ctx, input.Testimonial, userID)
	if err != nil {
		return nil, err
	}

	input.Testimonial = *testimonial
	alumniData := mappers.MapAlumniTestimonialGoToFirestore(input)

	_, err = s.firestoreClient.Collection("alumni_testimonials").Doc(input.ID).Set(ctx, alumniData)
	if err != nil {
		return nil, fmt.Errorf("failed to create alumni testimonial: %v", err)
	}

	return &input, nil
}

// CreateStudentSpotlight creates a new student spotlight
func (s *TestimonialService) CreateStudentSpotlight(ctx context.Context, input models.StudentSpotlight, userID string) (*models.StudentSpotlight, error) {
	if input.CurrentSemester == 0 {
		return nil, errors.New("current semester is required")
	}

	testimonial, err := s.CreateTestimonial(ctx, input.Testimonial, userID)
	if err != nil {
		return nil, err
	}

	input.Testimonial = *testimonial
	spotlightData := mappers.MapStudentSpotlightGoToFirestore(input)

	_, err = s.firestoreClient.Collection("student_spotlights").Doc(input.ID).Set(ctx, spotlightData)
	if err != nil {
		return nil, fmt.Errorf("failed to create student spotlight: %v", err)
	}

	return &input, nil
}

// GetTestimonial retrieves a testimonial by ID
func (s *TestimonialService) GetTestimonial(ctx context.Context, testimonialID string) (*models.Testimonial, error) {
	doc, err := s.firestoreClient.Collection("testimonials").Doc(testimonialID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get testimonial: %v", err)
	}

	testimonial := mappers.MapTestimonialFirestoreToGo(doc.Data())
	return &testimonial, nil
}

// GetAlumniTestimonial retrieves an alumni testimonial by ID
func (s *TestimonialService) GetAlumniTestimonial(ctx context.Context, testimonialID string) (*models.AlumniTestimonial, error) {
	doc, err := s.firestoreClient.Collection("alumni_testimonials").Doc(testimonialID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get alumni testimonial: %v", err)
	}

	testimonial := mappers.MapAlumniTestimonialFirestoreToGo(doc.Data())
	return &testimonial, nil
}

// GetStudentSpotlight retrieves a student spotlight by ID
func (s *TestimonialService) GetStudentSpotlight(ctx context.Context, testimonialID string) (*models.StudentSpotlight, error) {
	doc, err := s.firestoreClient.Collection("student_spotlights").Doc(testimonialID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get student spotlight: %v", err)
	}

	spotlight := mappers.MapStudentSpotlightFirestoreToGo(doc.Data())
	return &spotlight, nil
}

// UpdateTestimonial updates an existing testimonial
func (s *TestimonialService) UpdateTestimonial(ctx context.Context, testimonialID string, updates models.Testimonial) error {
	updates.UpdatedAt = time.Now()
	updateData := mappers.MapTestimonialGoToFirestore(updates)

	_, err := s.firestoreClient.Collection("testimonials").Doc(testimonialID).Set(ctx, updateData, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("failed to update testimonial: %v", err)
	}

	return nil
}

// PublishTestimonial publishes a testimonial
func (s *TestimonialService) PublishTestimonial(ctx context.Context, testimonialID string) error {
	_, err := s.firestoreClient.Collection("testimonials").Doc(testimonialID).Update(ctx, []firestore.Update{
		{Path: "is_published", Value: true},
		{Path: "updated_at", Value: time.Now()},
	})
	return err
}

// ListTestimonials retrieves all published testimonials
func (s *TestimonialService) ListTestimonials(ctx context.Context) ([]models.Testimonial, error) {
	query := s.firestoreClient.Collection("testimonials").Where("is_published", "==", true)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil && status.Code(err) != codes.NotFound {
		return nil, fmt.Errorf("failed to get testimonials: %v", err)
	}

	var testimonials []models.Testimonial
	for _, doc := range docs {
		testimonial := mappers.MapTestimonialFirestoreToGo(doc.Data())
		testimonials = append(testimonials, testimonial)
	}

	return testimonials, nil
}

// AddLike increments the like count for a testimonial
func (s *TestimonialService) AddLike(ctx context.Context, testimonialID string) error {
	_, err := s.firestoreClient.Collection("testimonials").Doc(testimonialID).Update(ctx, []firestore.Update{
		{Path: "likes", Value: firestore.Increment(1)},
		{Path: "updated_at", Value: time.Now()},
	})
	return err
}

// DeleteTestimonial deletes a testimonial
func (s *TestimonialService) DeleteTestimonial(ctx context.Context, testimonialID string, userID string) error {
	testimonial, err := s.GetTestimonial(ctx, testimonialID)
	if err != nil {
		return err
	}

	if testimonial.UserID != userID {
		return errors.New("unauthorized: only the testimonial author can delete it")
	}

	_, err = s.firestoreClient.Collection("testimonials").Doc(testimonialID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete testimonial: %v", err)
	}

	return nil
}
