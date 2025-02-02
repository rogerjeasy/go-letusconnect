package services

type ProjectCoreService struct {
	firestoreClient FirestoreClient
}

func NewProjectCoreService(client FirestoreClient) *ProjectCoreService {
	return &ProjectCoreService{
		firestoreClient: client,
	}
}

func (s *ProjectCoreService) GetGroupMembers(projectID string, groupID *string) ([]string, error) {
	// Logic to fetch group members from Firestore or another data source
	return []string{}, nil
}
