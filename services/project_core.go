package services

import (
	"cloud.google.com/go/firestore"
)

type ProjectCoreService struct {
	firestoreClient *firestore.Client
}

func NewProjectCoreService(client *firestore.Client) *ProjectCoreService {
	return &ProjectCoreService{
		firestoreClient: client,
	}
}

func (s *ProjectCoreService) GetGroupMembers(projectID string, groupID *string) ([]string, error) {
	// Logic to fetch group members from Firestore or another data source
	return []string{}, nil
}
