package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rogerjeasy/go-letusconnect/models"
)

type LinkedInJobsService struct {
	firestoreClient FirestoreClient
}

func NewLinkedInJobsService(client FirestoreClient) *LinkedInJobsService {
	return &LinkedInJobsService{
		firestoreClient: client,
	}
}

func (s *LinkedInJobsService) GetAppliedJobs(accessToken string) ([]models.LinkedInJobApplication, error) {
	client := &http.Client{}
	url := "https://api.linkedin.com/v2/applicationHistory"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch applied jobs: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LinkedIn API returned status code %d", resp.StatusCode)
	}

	var result struct {
		Elements []models.LinkedInJobApplication `json:"elements"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return result.Elements, nil
}

func (s *LinkedInJobsService) StoreJobApplication(ctx context.Context, userID string, jobApp models.LinkedInJobApplication) error {
	jobApp.ApplyTimestamp = time.Now().Unix()
	jobApp.JobID = uuid.New().String()

	_, err := s.firestoreClient.Collection("jobs").Doc(jobApp.JobID).Set(ctx, jobApp)
	if err != nil {
		return fmt.Errorf("failed to store job application: %v", err)
	}

	return nil
}

func (s *LinkedInJobsService) GetUserApplications(ctx context.Context, userID string) ([]models.LinkedInJobApplication, error) {
	docs, err := s.firestoreClient.Collection("job_applications").Doc(userID).Collection("applications").Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch job applications: %v", err)
	}

	var applications []models.LinkedInJobApplication
	for _, doc := range docs {
		var app models.LinkedInJobApplication
		if err := doc.DataTo(&app); err != nil {
			return nil, fmt.Errorf("failed to parse job application: %v", err)
		}
		applications = append(applications, app)
	}

	return applications, nil
}
