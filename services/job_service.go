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
)

// JobService manages job tracking operations
type JobService struct {
	firestoreClient FirestoreClient
}

// NewJobService initializes a new JobService
func NewJobService(client FirestoreClient) *JobService {
	return &JobService{firestoreClient: client}
}

// CreateJob adds a new job application to Firestore
func (s *JobService) CreateJob(ctx context.Context, job *models.Job) (*models.Job, error) {
	job.ID = uuid.New().String()
	job.CreatedAt = time.Now()
	job.UpdatedAt = job.CreatedAt

	if job.UserID == "" {
		return nil, fmt.Errorf("user ID is required. Please log in and try again")
	}

	if job.Company == "" {
		return nil, fmt.Errorf("company name is required")
	}

	if job.Position == "" {
		return nil, fmt.Errorf("position is required")
	}

	// Convert job struct to Firestore format
	firestoreData := mappers.MapJobGoToFirestore(*job)

	// Store job in Firestore
	_, err := s.firestoreClient.Collection("jobs").Doc(job.ID).Set(ctx, firestoreData)
	if err != nil {
		return nil, fmt.Errorf("failed to create job: %v", err)
	}

	return job, nil
}

// GetJob retrieves a job application by ID
func (s *JobService) GetJob(ctx context.Context, jobID string, userID string) (*models.Job, error) {
	doc, err := s.firestoreClient.Collection("jobs").Doc(jobID).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get job: %v", err)
	}

	// Convert Firestore data to Go struct
	job := mappers.MapJobFirestoreToGo(doc.Data())

	// Ensure the job belongs to the user
	if job.UserID != userID {
		return nil, fmt.Errorf("unauthorized access to job")
	}

	return &job, nil
}

// GetJobsByUser retrieves all job applications for a user
func (s *JobService) GetJobsByUser(ctx context.Context, userID string) ([]models.Job, error) {
	iter := s.firestoreClient.Collection("jobs").Where("user_id", "==", userID).Documents(ctx)

	var jobs []models.Job
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to fetch jobs: %v", err)
		}
		job := mappers.MapJobFirestoreToGo(doc.Data())
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// UpdateJob updates an existing job application
func (s *JobService) UpdateJob(ctx context.Context, jobID string, userID string, jobData map[string]interface{}) error {
	docRef := s.firestoreClient.Collection("jobs").Doc(jobID)

	// Fetch existing job to check ownership
	doc, err := docRef.Get(ctx)
	if err != nil {
		return fmt.Errorf("job not found: %v", err)
	}
	job := mappers.MapJobFirestoreToGo(doc.Data())

	// Ensure the user is the owner
	if job.UserID != userID {
		return fmt.Errorf("unauthorized: cannot update job")
	}

	// Add updated timestamp
	jobData["updated_at"] = time.Now()

	// Update Firestore document
	_, err = docRef.Update(ctx, []firestore.Update{
		{Path: "updated_at", Value: jobData["updated_at"]},
		{Path: "company", Value: jobData["company"]},
		{Path: "position", Value: jobData["position"]},
		{Path: "location", Value: jobData["location"]},
		{Path: "status", Value: jobData["status"]},
		{Path: "salary_range", Value: jobData["salary_range"]},
		{Path: "job_type", Value: jobData["job_type"]},
		{Path: "job_description", Value: jobData["job_description"]},
		{Path: "offer_details", Value: jobData["offer_details"]},
		{Path: "rejection_reason", Value: jobData["rejection_reason"]},
		{Path: "follow_up_date", Value: jobData["follow_up_date"]},
		{Path: "company_rating", Value: jobData["company_rating"]},
	})
	if err != nil {
		return fmt.Errorf("failed to update job: %v", err)
	}

	return nil
}

// DeleteJob removes a job application from Firestore
func (s *JobService) DeleteJob(ctx context.Context, jobID string, userID string) error {
	docRef := s.firestoreClient.Collection("jobs").Doc(jobID)

	// Fetch job to check ownership
	doc, err := docRef.Get(ctx)
	if err != nil {
		return fmt.Errorf("job not found: %v", err)
	}
	job := mappers.MapJobFirestoreToGo(doc.Data())

	// Ensure the user is the owner
	if job.UserID != userID {
		return fmt.Errorf("unauthorized: cannot delete job")
	}

	_, err = docRef.Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete job: %v", err)
	}

	return nil
}
