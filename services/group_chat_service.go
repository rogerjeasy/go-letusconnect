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
		ProjectID:     input.ProjectID,
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

// AddParticipantToGroupChat adds a participant to a given group chat by groupChatID or projectID
func AddParticipantToGroupChat(ctx context.Context, groupChatID, projectID string, participant models.Participant) error {
	// Validate input
	if groupChatID == "" && projectID == "" {
		return fmt.Errorf("either groupChatID or projectID must be provided")
	}

	if participant.UserID == "" {
		return fmt.Errorf("participant UserID is required")
	}

	var docRef *firestore.DocumentRef
	var docSnap *firestore.DocumentSnapshot
	var err error

	// Fetch the group chat document
	if groupChatID != "" {
		docRef = FirestoreClient.Collection("group_chats").Doc(groupChatID)
		docSnap, err = docRef.Get(ctx)
	} else if projectID != "" {
		query := FirestoreClient.Collection("group_chats").Where("project_id", "==", projectID).Limit(1)
		iter := query.Documents(ctx)
		docSnap, err = iter.Next()
		if err == iterator.Done {
			return fmt.Errorf("group chat not found for the given project ID")
		}
		if err == nil {
			docRef = docSnap.Ref
		}
	}

	if err != nil {
		return fmt.Errorf("failed to fetch group chat: %v", err)
	}

	data := docSnap.Data()
	if data == nil {
		return fmt.Errorf("group chat data is missing")
	}

	// Check if the participant already exists
	existingParticipants := mappers.GetParticipantsArray(data, "participants")
	for _, p := range existingParticipants {
		if p.UserID == participant.UserID {
			return fmt.Errorf("participant with UserID %s already exists in the group chat", participant.UserID)
		}
	}

	// Append the new participant as a models.Participant
	existingParticipants = append(existingParticipants, participant)

	// Map the participants back to Firestore format
	participantsFirestore := mappers.MapParticipantsArrayToFirestore(existingParticipants)

	// Update the Firestore document
	if _, err := docRef.Update(ctx, []firestore.Update{
		{Path: "participants", Value: participantsFirestore},
		{Path: "updated_at", Value: time.Now()},
	}); err != nil {
		return fmt.Errorf("failed to update group chat participants: %v", err)
	}

	return nil
}

func SendMessageService(ctx context.Context, groupChatID string, senderID string, senderName string, content string) (*models.BaseMessage, error) {
	// Validate required parameters
	if groupChatID == "" {
		return nil, fmt.Errorf("groupChatID is required")
	}
	if senderID == "" || senderName == "" {
		return nil, fmt.Errorf("sender information is required")
	}
	if content == "" {
		return nil, fmt.Errorf("message content cannot be empty")
	}

	// Fetch the group chat document
	docRef := FirestoreClient.Collection("group_chats").Doc(groupChatID)
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch group chat: %v", err)
	}

	data := docSnap.Data()
	if data == nil {
		return nil, fmt.Errorf("group chat not found")
	}

	// Retrieve participants to set read statuses
	participants := mappers.GetParticipantsArray(data, "participants")
	if len(participants) == 0 {
		return nil, fmt.Errorf("no participants found in the group chat")
	}

	readStatus := make(map[string]bool)
	for _, participant := range participants {
		if participant.UserID == "" {
			return nil, fmt.Errorf("participant UserID cannot be empty")
		}
		if participant.UserID == senderID {
			readStatus[participant.UserID] = true
		} else {
			readStatus[participant.UserID] = false
		}
	}

	// Create the new message
	message := models.BaseMessage{
		ID:          uuid.New().String(),
		SenderID:    senderID,
		SenderName:  senderName,
		Content:     content,
		CreatedAt:   time.Now().Format(time.RFC3339),
		ReadStatus:  readStatus,
		IsDeleted:   false,
		Attachments: []string{},
		Reactions:   make(map[string]int),
		MessageType: "text",
	}

	// Retrieve existing messages and append the new message
	messages := mappers.GetBaseMessagesArrayFromFirestore(data, "messages")
	if messages == nil {
		messages = []models.BaseMessage{} // Initialize as empty slice if nil
	}
	messages = append(messages, message)

	// Prepare Firestore payload
	firestorePayload := map[string]interface{}{
		"messages":   mappers.MapBaseMessagesArrayToFirestore(messages),
		"updated_at": time.Now(),
	}

	// Update Firestore
	if _, err := docRef.Set(ctx, firestorePayload, firestore.MergeAll); err != nil {
		return nil, fmt.Errorf("failed to update group chat with new message: %v", err)
	}

	// Return the new message
	return &message, nil
}
