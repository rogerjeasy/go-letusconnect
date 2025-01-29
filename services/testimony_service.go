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
)

type TestimonialService struct {
	firestoreClient *firestore.Client
	userService     *UserService
}

func NewTestimonialService(fClient *firestore.Client, uService *UserService) *TestimonialService {
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
