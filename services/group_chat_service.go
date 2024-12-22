package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"google.golang.org/api/iterator"
)

type GroupChatInput struct {
	ProjectID      string
	Name           string
	Description    string
	CreatedByUID   string
	CreatedByName  string
	Email          string
	ProfilePicture string
}

func CreateGroupChatService(ctx context.Context, input GroupChatInput) (*models.GroupChat, error) {
	chatID := uuid.New().String()

	// Create the group chat with defaults and provided values
	groupChat := models.GroupChat{
		ID:            chatID,
		ProjectID:     input.ProjectID, // Optional, can be empty
		CreatedByUID:  input.CreatedByUID,
		CreatedByName: input.CreatedByName,
		Name:          input.Name,
		Description:   input.Description,
		Participants: []models.Participant{
			{
				UserID:         input.CreatedByUID,
				Username:       input.CreatedByName,
				Role:           "owner",
				JoinedAt:       time.Now(),
				Email:          input.Email,
				ProfilePicture: input.ProfilePicture,
			},
		},
		Messages:       []models.BaseMessage{},
		PinnedMessages: []string{},
		IsArchived:     false,
		Notifications:  map[string]bool{input.CreatedByUID: true},
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		ReadStatus:     map[string]bool{input.CreatedByUID: true},
		GroupSettings: models.GroupSettings{
			AllowFileSharing:  true,
			AllowPinning:      true,
			AllowReactions:    true,
			AllowReplies:      true,
			MuteNotifications: false,
			OnlyAdminsCanPost: false,
		},
	}

	// Save to Firestore
	_, _, err := FirestoreClient.Collection("group_chats").Add(ctx, mappers.MapGroupChatGoToFirestore(groupChat))
	if err != nil {
		return nil, err
	}

	return &groupChat, nil
}

// GetGroupChatService fetches a group chat by project ID and returns it in frontend format
func GetGroupChatService(ctx context.Context, projectId string) (map[string]interface{}, error) {
	if projectId == "" {
		return nil, fmt.Errorf("project ID is required")
	}

	// Query Firestore for the group chat where projectId matches
	query := FirestoreClient.Collection("group_chats").Where("project_id", "==", projectId).Limit(1)
	iter := query.Documents(ctx)
	docSnap, err := iter.Next()
	if err == iterator.Done {
		return nil, fmt.Errorf("no group chat found for the given project ID")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch group chat: %v", err)
	}

	// Convert to map[string]interface{}
	data := docSnap.Data()
	if data == nil {
		return nil, fmt.Errorf("group chat not found")
	}

	// Ensure ID is in the data
	data["id"] = docSnap.Ref.ID

	// Convert Firestore data to frontend format using the mapper
	frontendData := mappers.MapGroupChatFirestoreToFrontend(data)

	return frontendData, nil
}

// GetGroupChatsByProjectService fetches all group chats for a project
func GetGroupChatsByProjectService(ctx context.Context, projectId string) ([]map[string]interface{}, error) {
	if projectId == "" {
		return nil, fmt.Errorf("project ID is required")
	}

	// Query Firestore for group chats with matching project ID
	query := FirestoreClient.Collection("group_chats").Where("project_id", "==", projectId)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch group chats: %v", err)
	}

	var groupChats []map[string]interface{}
	for _, doc := range docs {
		data := doc.Data()
		// Ensure ID is in the data
		data["id"] = doc.Ref.ID

		// Convert each group chat to frontend format
		frontendData := mappers.MapGroupChatFirestoreToFrontend(data)
		groupChats = append(groupChats, frontendData)
	}

	return groupChats, nil
}

// GetGroupChatsByUserService fetches all group chats where the user is a participant
func GetGroupChatsByUserService(ctx context.Context, userId string) ([]map[string]interface{}, error) {
	if userId == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	// Query Firestore for group chats where user is a participant
	// This assumes participants array contains objects with userId field
	query := FirestoreClient.Collection("group_chats").Where("participants", "array-contains", map[string]interface{}{
		"user_id": userId,
	})

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch group chats: %v", err)
	}

	var groupChats []map[string]interface{}
	for _, doc := range docs {
		data := doc.Data()
		// Ensure ID is in the data
		data["id"] = doc.Ref.ID

		// Convert each group chat to frontend format
		frontendData := mappers.MapGroupChatFirestoreToFrontend(data)
		groupChats = append(groupChats, frontendData)
	}

	return groupChats, nil
}
