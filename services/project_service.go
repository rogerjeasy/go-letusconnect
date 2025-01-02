package services

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
)

type ProjectService struct {
	firestoreClient *firestore.Client
	userService     *UserService
}

func NewProjectService(client *firestore.Client, userService *UserService) *ProjectService {
	return &ProjectService{
		firestoreClient: client,
		userService:     userService,
	}
}

func GetGroupMembers(projectID string, groupID *string) ([]string, error) {
	// Logic to fetch group members from Firestore or another data source
	return []string{}, nil
}

func (s *ProjectService) JoinProject(ctx context.Context, projectID, userID, message string) error {
	// Get user details
	user, err := s.userService.GetUserByUID(userID)
	if err != nil {
		return fmt.Errorf("failed to fetch user details: %v", err)
	}

	// Create join request
	joinRequest := models.JoinRequest{
		UserID:         userID,
		Username:       user["username"].(string),
		ProfilePicture: user["profile_picture"].(string),
		Email:          user["email"].(string),
		Message:        message,
		RequestedAt:    time.Now(),
		Status:         "pending",
	}

	// Get project
	doc, err := s.firestoreClient.Collection("projects").Doc(projectID).Get(ctx)
	if err != nil {
		return fmt.Errorf("project not found: %v", err)
	}

	projectData := doc.Data()

	// Validate project status
	if projectData["status"] == "completed" {
		return fmt.Errorf("this project has been completed")
	}

	// Check if user is owner
	if projectData["owner_id"] == userID {
		return fmt.Errorf("owners cannot join their own project")
	}

	// Check existing requests
	joinRequests := mappers.GetJoinRequestsArray(projectData, "join_requests")
	for _, jr := range joinRequests {
		if jr.UserID == userID {
			return fmt.Errorf("you have already applied to join this project")
		}
	}

	// Check rejected participants
	if rejectedParticipants, ok := projectData["rejected_participants"].([]interface{}); ok {
		for _, rejectedUID := range rejectedParticipants {
			if rejectedUID == userID {
				return fmt.Errorf("your request was previously rejected")
			}
		}
	}

	// Add join request
	joinRequestMap := mappers.MapJoinRequestGoToFirestore(joinRequest)
	_, err = s.firestoreClient.Collection("projects").Doc(projectID).Update(ctx, []firestore.Update{
		{Path: "join_requests", Value: firestore.ArrayUnion(joinRequestMap)},
	})
	if err != nil {
		return fmt.Errorf("failed to apply to join project: %v", err)
	}

	return nil
}
