package services

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
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
	existingParticipants := mappers.GetParticipantsGoArray(data, "participants")
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
	participants := mappers.GetParticipantsGoArray(data, "participants")
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
		messages = []models.BaseMessage{}
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

func MarkMessagesAsReadService(ctx context.Context, groupChatID, userID string) error {
	// Validate required parameters
	if groupChatID == "" {
		return fmt.Errorf("groupChatID is required")
	}
	if userID == "" {
		return fmt.Errorf("userID is required")
	}

	// Fetch the group chat document
	docRef := FirestoreClient.Collection("group_chats").Doc(groupChatID)
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch group chat: %v", err)
	}

	data := docSnap.Data()
	if data == nil {
		return fmt.Errorf("group chat not found")
	}

	// Retrieve existing messages
	messages := mappers.GetBaseMessagesArrayFromFirestore(data, "messages")
	if messages == nil || len(messages) == 0 {
		return fmt.Errorf("no messages found in the group chat")
	}

	// Update the `read_status` for the given user in each message
	for i := range messages {
		if messages[i].ReadStatus == nil {
			messages[i].ReadStatus = make(map[string]bool)
		}
		messages[i].ReadStatus[userID] = true
	}

	// Map updated messages to Firestore format
	firestoreMessages := mappers.MapBaseMessagesArrayToFirestore(messages)

	// Update Firestore document
	if _, err := docRef.Update(ctx, []firestore.Update{
		{Path: "messages", Value: firestoreMessages},
		{Path: "updated_at", Value: time.Now()},
	}); err != nil {
		return fmt.Errorf("failed to update group chat messages: %v", err)
	}

	return nil
}

func CountUnreadMessagesService(ctx context.Context, groupChatID, projectID, userID string) (int, error) {
	// Validate required parameters
	if groupChatID == "" && projectID == "" {
		return 0, fmt.Errorf("either groupChatID or projectID is required")
	}
	if userID == "" {
		return 0, fmt.Errorf("userID is required")
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
	}

	if err != nil {
		return 0, fmt.Errorf("failed to fetch group chat: %v", err)
	}

	data := docSnap.Data()
	if data == nil {
		return 0, fmt.Errorf("group chat not found")
	}

	// Retrieve existing messages
	messages := mappers.GetBaseMessagesArrayFromFirestore(data, "messages")
	if messages == nil || len(messages) == 0 {
		return 0, nil // No messages means no unread messages
	}

	// Count unread messages for the user
	unreadCount := 0
	for _, message := range messages {
		if read, ok := message.ReadStatus[userID]; !ok || !read {
			unreadCount++
		}
	}

	return unreadCount, nil
}

func RemoveParticipantFromGroupChatService(ctx context.Context, groupChatID, ownerID, participantID string) error {
	// Validate required parameters
	if groupChatID == "" {
		return fmt.Errorf("groupChatID is required")
	}
	if ownerID == "" {
		return fmt.Errorf("ownerID is required")
	}
	if participantID == "" {
		return fmt.Errorf("participantID is required")
	}

	// Fetch the group chat document
	docRef := FirestoreClient.Collection("group_chats").Doc(groupChatID)
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch group chat: %v", err)
	}

	data := docSnap.Data()
	if data == nil {
		return fmt.Errorf("group chat not found")
	}

	// Retrieve existing participants
	existingParticipants := mappers.GetParticipantsGoArray(data, "participants")
	if len(existingParticipants) == 0 {
		return fmt.Errorf("no participants found in the group chat")
	}

	// Check if the owner has the required role
	isOwner := false
	for _, participant := range existingParticipants {
		if participant.UserID == ownerID && participant.Role == "owner" {
			isOwner = true
			break
		}
	}
	if !isOwner {
		return fmt.Errorf("only an owner can remove participants")
	}

	// Check if the participant exists and remove them
	updatedParticipants := []models.Participant{}
	participantFound := false
	for _, participant := range existingParticipants {
		if participant.UserID == participantID {
			participantFound = true
			continue // Skip adding this participant to the updated list
		}
		updatedParticipants = append(updatedParticipants, participant)
	}

	if !participantFound {
		return fmt.Errorf("participant with ID %s not found", participantID)
	}

	// Map the updated participants to Firestore format
	participantsFirestore := mappers.MapParticipantsArrayToFirestore(updatedParticipants)

	// Update the Firestore document
	if _, err := docRef.Update(ctx, []firestore.Update{
		{Path: "participants", Value: participantsFirestore},
		{Path: "updated_at", Value: time.Now()},
	}); err != nil {
		return fmt.Errorf("failed to update group chat participants: %v", err)
	}

	return nil
}

func ReplyToMessageService(ctx context.Context, groupChatID, senderID, senderName, content, messageIDToReply string) (*models.BaseMessage, error) {
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
	if messageIDToReply == "" {
		return nil, fmt.Errorf("messageIDToReply is required")
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

	// Retrieve existing messages
	messages := mappers.GetBaseMessagesArrayFromFirestore(data, "messages")
	if messages == nil || len(messages) == 0 {
		return nil, fmt.Errorf("no messages found in the group chat")
	}

	// Find the message being replied to
	var repliedMessage *models.BaseMessage
	for _, msg := range messages {
		if msg.ID == messageIDToReply {
			repliedMessage = &msg
			break
		}
	}
	if repliedMessage == nil {
		return nil, fmt.Errorf("message with ID %s not found", messageIDToReply)
	}

	// Retrieve participants to set read statuses
	participants := mappers.GetParticipantsGoArray(data, "participants")
	if len(participants) == 0 {
		return nil, fmt.Errorf("no participants found in the group chat")
	}

	readStatus := make(map[string]bool)
	for _, participant := range participants {
		if participant.UserID == senderID {
			readStatus[participant.UserID] = true
		} else {
			readStatus[participant.UserID] = false
		}
	}

	// Create the reply message
	replyMessage := models.BaseMessage{
		ID:          uuid.New().String(),
		SenderID:    senderID,
		SenderName:  senderName,
		Content:     content,
		CreatedAt:   time.Now().Format(time.RFC3339),
		ReadStatus:  readStatus,
		IsDeleted:   false,
		Attachments: []string{},
		Reactions:   make(map[string]int),
		MessageType: "reply",
		ReplyToID:   &messageIDToReply, // Reference to the original message
	}

	// Append the reply message to the group chat
	messages = append(messages, replyMessage)

	// Prepare Firestore payload
	firestorePayload := map[string]interface{}{
		"messages":   mappers.MapBaseMessagesArrayToFirestore(messages),
		"updated_at": time.Now(),
	}

	// Update Firestore
	if _, err := docRef.Set(ctx, firestorePayload, firestore.MergeAll); err != nil {
		return nil, fmt.Errorf("failed to update group chat with the reply: %v", err)
	}

	// Return the reply message
	return &replyMessage, nil
}

func AttachFilesToMessageService(ctx context.Context, groupChatID, senderID, senderName, content string, files []*multipart.FileHeader) (*models.BaseMessage, error) {
	// Validate required parameters
	if groupChatID == "" {
		return nil, fmt.Errorf("groupChatID is required")
	}
	if senderID == "" || senderName == "" {
		return nil, fmt.Errorf("sender information is required")
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("at least one file is required")
	}

	// Initialize Cloudinary client
	cld := CloudinaryClient
	if cld == nil {
		return nil, fmt.Errorf("Cloudinary client not initialized")
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

	// Upload files to Cloudinary
	var attachments []string
	for _, fileHeader := range files {
		// Open the file
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %v", fileHeader.Filename, err)
		}
		defer file.Close()

		// Determine the resource type based on file extension
		resourceType := "auto"
		switch {
		case isVideo(fileHeader.Filename):
			resourceType = "video"
		case isDocument(fileHeader.Filename):
			resourceType = "raw"
		}

		// Upload to Cloudinary
		uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
			PublicID:     uuid.New().String(),
			Folder:       fmt.Sprintf("group_chats/%s/files", groupChatID),
			ResourceType: resourceType,
		})
		if err != nil {
			log.Printf("Error uploading file %s: %v", fileHeader.Filename, err)
			return nil, fmt.Errorf("failed to upload file %s: %v", fileHeader.Filename, err)
		}

		attachments = append(attachments, uploadResult.SecureURL)
	}

	// Retrieve participants to set read statuses
	participants := mappers.GetParticipantsGoArray(data, "participants")
	if len(participants) == 0 {
		return nil, fmt.Errorf("no participants found in the group chat")
	}

	readStatus := make(map[string]bool)
	for _, participant := range participants {
		if participant.UserID == senderID {
			readStatus[participant.UserID] = true
		} else {
			readStatus[participant.UserID] = false
		}
	}

	// Create the message
	message := models.BaseMessage{
		ID:          uuid.New().String(),
		SenderID:    senderID,
		SenderName:  senderName,
		Content:     content,
		CreatedAt:   time.Now().Format(time.RFC3339),
		ReadStatus:  readStatus,
		IsDeleted:   false,
		Attachments: attachments,
		Reactions:   make(map[string]int),
		MessageType: "attachment",
	}

	// Append the message to the group chat
	messages := mappers.GetBaseMessagesArrayFromFirestore(data, "messages")
	if messages == nil {
		messages = []models.BaseMessage{}
	}
	messages = append(messages, message)

	// Update Firestore
	firestorePayload := map[string]interface{}{
		"messages":   mappers.MapBaseMessagesArrayToFirestore(messages),
		"updated_at": time.Now(),
	}

	if _, err := docRef.Set(ctx, firestorePayload, firestore.MergeAll); err != nil {
		return nil, fmt.Errorf("failed to update group chat with new message: %v", err)
	}

	return &message, nil
}

func isVideo(filename string) bool {
	videoExtensions := []string{".mp4", ".mov", ".avi", ".mkv"}
	return hasExtension(filename, videoExtensions)
}

func isDocument(filename string) bool {
	documentExtensions := []string{".pdf", ".doc", ".docx", ".ppt", ".xls"}
	return hasExtension(filename, documentExtensions)
}

func hasExtension(filename string, extensions []string) bool {
	for _, ext := range extensions {
		if strings.HasSuffix(strings.ToLower(filename), ext) {
			return true
		}
	}
	return false
}

func PinMessageService(ctx context.Context, groupChatID, userID, messageID string) error {
	// Validate required parameters
	if groupChatID == "" {
		return fmt.Errorf("groupChatID is required")
	}
	if userID == "" {
		return fmt.Errorf("userID is required")
	}
	if messageID == "" {
		return fmt.Errorf("messageID is required")
	}

	// Fetch the group chat document
	docRef := FirestoreClient.Collection("group_chats").Doc(groupChatID)
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch group chat: %v", err)
	}

	data := docSnap.Data()
	if data == nil {
		return fmt.Errorf("group chat not found")
	}

	// Retrieve participants to check user permissions
	participants := mappers.GetParticipantsGoArray(data, "participants")
	isAuthorized := false
	for _, participant := range participants {
		if participant.UserID == userID && (participant.Role == "owner" || participant.Role == "admin") {
			isAuthorized = true
			break
		}
	}

	if !isAuthorized {
		return fmt.Errorf("only an owner or admin can pin messages")
	}

	// Retrieve existing messages and pinned messages
	messages := mappers.GetBaseMessagesArrayFromFirestore(data, "messages")
	pinnedMessages := mappers.GetStringArray(data, "pinned_messages")

	// Check if the message exists
	var messageExists bool
	for _, msg := range messages {
		if msg.ID == messageID {
			messageExists = true
			break
		}
	}
	if !messageExists {
		return fmt.Errorf("message with ID %s not found", messageID)
	}

	// Check if the message is already pinned
	for _, pinnedMessage := range pinnedMessages {
		if pinnedMessage == messageID {
			return fmt.Errorf("message with ID %s is already pinned", messageID)
		}
	}

	// Pin the message
	pinnedMessages = append(pinnedMessages, messageID)

	// Update Firestore
	if _, err := docRef.Update(ctx, []firestore.Update{
		{Path: "pinned_messages", Value: pinnedMessages},
		{Path: "updated_at", Value: time.Now()},
	}); err != nil {
		return fmt.Errorf("failed to pin message: %v", err)
	}

	return nil
}

func GetPinnedMessagesService(ctx context.Context, groupChatID string) ([]models.BaseMessage, error) {
	// Validate required parameters
	if groupChatID == "" {
		return nil, fmt.Errorf("groupChatID is required")
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

	// Retrieve pinned messages
	pinnedMessageIDs := mappers.GetStringArray(data, "pinned_messages")
	if len(pinnedMessageIDs) == 0 {
		return []models.BaseMessage{}, nil // No pinned messages
	}

	// Retrieve all messages
	messages := mappers.GetBaseMessagesArrayFromFirestore(data, "messages")
	if messages == nil || len(messages) == 0 {
		return nil, fmt.Errorf("no messages found in the group chat")
	}

	// Filter messages to include only pinned messages
	var pinnedMessages []models.BaseMessage
	messageMap := make(map[string]models.BaseMessage)
	for _, message := range messages {
		messageMap[message.ID] = message
	}

	for _, pinnedID := range pinnedMessageIDs {
		if pinnedMessage, exists := messageMap[pinnedID]; exists {
			pinnedMessages = append(pinnedMessages, pinnedMessage)
		}
	}

	return pinnedMessages, nil
}

func UnpinMessageService(ctx context.Context, groupChatID, userID, messageID string) error {
	// Validate required parameters
	if groupChatID == "" {
		return fmt.Errorf("groupChatID is required")
	}
	if userID == "" {
		return fmt.Errorf("userID is required")
	}
	if messageID == "" {
		return fmt.Errorf("messageID is required")
	}

	// Fetch the group chat document
	docRef := FirestoreClient.Collection("group_chats").Doc(groupChatID)
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch group chat: %v", err)
	}

	data := docSnap.Data()
	if data == nil {
		return fmt.Errorf("group chat not found")
	}

	// Retrieve participants to check user permissions
	participants := mappers.GetParticipantsGoArray(data, "participants")
	isAuthorized := false
	for _, participant := range participants {
		if participant.UserID == userID && (participant.Role == "owner" || participant.Role == "admin") {
			isAuthorized = true
			break
		}
	}

	if !isAuthorized {
		return fmt.Errorf("only an owner or admin can unpin messages")
	}

	// Retrieve pinned messages
	pinnedMessages := mappers.GetStringArray(data, "pinned_messages")
	if len(pinnedMessages) == 0 {
		return fmt.Errorf("no pinned messages to unpin")
	}

	// Remove the message ID from pinned messages
	var updatedPinnedMessages []string
	messageFound := false
	for _, pinnedID := range pinnedMessages {
		if pinnedID == messageID {
			messageFound = true
			continue
		}
		updatedPinnedMessages = append(updatedPinnedMessages, pinnedID)
	}

	if !messageFound {
		return fmt.Errorf("message with ID %s is not pinned", messageID)
	}

	// Update Firestore
	if _, err := docRef.Update(ctx, []firestore.Update{
		{Path: "pinned_messages", Value: updatedPinnedMessages},
		{Path: "updated_at", Value: time.Now()},
	}); err != nil {
		return fmt.Errorf("failed to unpin message: %v", err)
	}

	return nil
}

func ReactToMessageService(ctx context.Context, groupChatID, userID, messageID, reaction string) error {
	if groupChatID == "" || userID == "" || messageID == "" || reaction == "" {
		return fmt.Errorf("all parameters (groupChatID, userID, messageID, reaction) are required")
	}

	docRef := FirestoreClient.Collection("group_chats").Doc(groupChatID)
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch group chat: %v", err)
	}

	data := docSnap.Data()
	messages := mappers.GetBaseMessagesArrayFromFirestore(data, "messages")
	for i, msg := range messages {
		if msg.ID == messageID {
			if msg.Reactions == nil {
				msg.Reactions = map[string]int{}
			}
			msg.Reactions[reaction]++
			messages[i] = msg
			break
		}
	}

	firestorePayload := map[string]interface{}{
		"messages":   mappers.MapBaseMessagesArrayToFirestore(messages),
		"updated_at": time.Now(),
	}

	if _, err := docRef.Set(ctx, firestorePayload, firestore.MergeAll); err != nil {
		return fmt.Errorf("failed to update message reactions: %v", err)
	}

	return nil
}

func GetMessageReadReceiptsService(ctx context.Context, groupChatID, messageID string) (map[string]bool, error) {
	if groupChatID == "" || messageID == "" {
		return nil, fmt.Errorf("groupChatID and messageID are required")
	}

	docRef := FirestoreClient.Collection("group_chats").Doc(groupChatID)
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch group chat: %v", err)
	}

	data := docSnap.Data()
	messages := mappers.GetBaseMessagesArrayFromFirestore(data, "messages")
	for _, msg := range messages {
		if msg.ID == messageID {
			return msg.ReadStatus, nil
		}
	}

	return nil, fmt.Errorf("message not found")
}

func SetParticipantRoleService(ctx context.Context, groupChatID, userID, participantID, newRole string) error {
	docRef := FirestoreClient.Collection("group_chats").Doc(groupChatID)
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch group chat: %v", err)
	}

	data := docSnap.Data()
	participants := mappers.GetParticipantsGoArray(data, "participants")
	isAdmin := false
	for _, participant := range participants {
		if participant.UserID == userID && (participant.Role == "owner" || participant.Role == "admin") {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return fmt.Errorf("only admins or owners can set roles")
	}

	for i, participant := range participants {
		if participant.UserID == participantID {
			participants[i].Role = newRole
			break
		}
	}

	firestorePayload := map[string]interface{}{
		"participants": mappers.MapParticipantsArrayToFirestore(participants),
		"updated_at":   time.Now(),
	}

	if _, err := docRef.Update(ctx, []firestore.Update{
		{Path: "participants", Value: firestorePayload["participants"]},
		{Path: "updated_at", Value: firestorePayload["updated_at"]},
	}); err != nil {
		return fmt.Errorf("failed to update participant role: %v", err)
	}

	return nil
}

func MuteParticipantService(ctx context.Context, groupChatID, userID, participantID string, duration time.Duration) error {
	if groupChatID == "" || userID == "" || participantID == "" {
		return fmt.Errorf("groupChatID, userID, and participantID are required")
	}

	docRef := FirestoreClient.Collection("group_chats").Doc(groupChatID)
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch group chat: %v", err)
	}

	data := docSnap.Data()
	participants := mappers.GetParticipantsGoArray(data, "participants")

	isAdmin := false
	for _, participant := range participants {
		if participant.UserID == userID && (participant.Role == "owner" || participant.Role == "admin") {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return fmt.Errorf("only admins or owners can mute participants")
	}

	for i, participant := range participants {
		if participant.UserID == participantID {
			participants[i].MutedUntil = time.Now().Add(duration)
			break
		}
	}

	firestorePayload := map[string]interface{}{
		"participants": mappers.MapParticipantsArrayToFirestore(participants),
		"updated_at":   time.Now(),
	}

	if _, err := docRef.Update(ctx, []firestore.Update{
		{Path: "participants", Value: firestorePayload["participants"]},
		{Path: "updated_at", Value: firestorePayload["updated_at"]},
	}); err != nil {
		return fmt.Errorf("failed to mute participant: %v", err)
	}

	return nil
}

func UpdateLastSeenService(ctx context.Context, groupChatID, userID string) error {
	if groupChatID == "" || userID == "" {
		return fmt.Errorf("groupChatID and userID are required")
	}

	docRef := FirestoreClient.Collection("group_chats").Doc(groupChatID)
	_, err := docRef.Update(ctx, []firestore.Update{
		{Path: fmt.Sprintf("last_seen.%s", userID), Value: time.Now()},
		{Path: "updated_at", Value: time.Now()},
	})

	if err != nil {
		return fmt.Errorf("failed to update last seen status: %v", err)
	}

	return nil
}

func ArchiveGroupChatService(ctx context.Context, groupChatID, userID string) error {
	docRef := FirestoreClient.Collection("group_chats").Doc(groupChatID)
	_, err := docRef.Update(ctx, []firestore.Update{
		{Path: "is_archived", Value: true},
		{Path: "updated_at", Value: time.Now()},
	})

	if err != nil {
		return fmt.Errorf("failed to archive group chat: %v", err)
	}

	return nil
}

func LeaveGroupService(ctx context.Context, groupChatID, userID string) error {
	if groupChatID == "" || userID == "" {
		return fmt.Errorf("groupChatID and userID are required")
	}

	docRef := FirestoreClient.Collection("group_chats").Doc(groupChatID)
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch group chat: %v", err)
	}

	data := docSnap.Data()
	participants := mappers.GetParticipantsGoArray(data, "participants")

	updatedParticipants := []models.Participant{}
	userFound := false
	for _, participant := range participants {
		if participant.UserID == userID {
			userFound = true
			continue
		}
		updatedParticipants = append(updatedParticipants, participant)
	}

	if !userFound {
		return fmt.Errorf("user not found in the group chat")
	}

	if len(updatedParticipants) == 0 {
		// If no participants remain, delete the group chat
		_, err := docRef.Delete(ctx)
		return err
	}

	// Update Firestore with the updated participants
	firestorePayload := map[string]interface{}{
		"participants": mappers.MapParticipantsArrayToFirestore(updatedParticipants),
		"updated_at":   time.Now(),
	}

	if _, err := docRef.Update(ctx, []firestore.Update{
		{Path: "participants", Value: firestorePayload["participants"]},
		{Path: "updated_at", Value: firestorePayload["updated_at"]},
	}); err != nil {
		return fmt.Errorf("failed to update group chat participants: %v", err)
	}

	return nil
}

func CreatePollService(ctx context.Context, groupChatID, userID string, poll models.Poll) (*models.Poll, error) {
	if groupChatID == "" || userID == "" {
		return nil, fmt.Errorf("groupChatID and userID are required")
	}

	docRef := FirestoreClient.Collection("group_chats").Doc(groupChatID)
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch group chat: %v", err)
	}

	data := docSnap.Data()
	participants := mappers.GetParticipantsGoArray(data, "participants")
	userFound := false
	for _, participant := range participants {
		if participant.UserID == userID {
			userFound = true
			break
		}
	}

	if !userFound {
		return nil, fmt.Errorf("user is not a participant in the group chat")
	}

	poll.ID = uuid.New().String()
	poll.CreatedBy = userID
	poll.CreatedAt = time.Now()

	firestorePayload := map[string]interface{}{
		"polls":      append(mappers.GetPollsArray(data, "polls"), poll),
		"updated_at": time.Now(),
	}

	if _, err := docRef.Update(ctx, []firestore.Update{
		{Path: "polls", Value: firestorePayload["polls"]},
		{Path: "updated_at", Value: firestorePayload["updated_at"]},
	}); err != nil {
		return nil, fmt.Errorf("failed to create poll: %v", err)
	}

	return &poll, nil
}

func ReportMessageService(ctx context.Context, groupChatID, userID, messageID, reason string) error {
	if groupChatID == "" || userID == "" || messageID == "" || reason == "" {
		return fmt.Errorf("all fields are required")
	}

	docRef := FirestoreClient.Collection("group_chats").Doc(groupChatID)
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch group chat: %v", err)
	}

	data := docSnap.Data()
	reports := mappers.GetReportsArray(data, "reports")

	report := models.Report{
		ID:         uuid.New().String(),
		MessageID:  messageID,
		ReportedBy: userID,
		Reason:     reason,
		CreatedAt:  time.Now(),
	}

	reports = append(reports, report)

	firestorePayload := map[string]interface{}{
		"reports":    mappers.MapReportsArrayToFirestore(reports),
		"updated_at": time.Now(),
	}

	if _, err := docRef.Update(ctx, []firestore.Update{
		{Path: "reports", Value: firestorePayload["reports"]},
		{Path: "updated_at", Value: firestorePayload["updated_at"]},
	}); err != nil {
		return fmt.Errorf("failed to report message: %v", err)
	}

	return nil
}

// UpdateGroupSettingsService updates the settings of a group chat
func UpdateGroupSettingsService(ctx context.Context, groupChatID, userID string, updatedSettings models.GroupSettings) error {
	if groupChatID == "" || userID == "" {
		return fmt.Errorf("groupChatID and userID are required")
	}

	// Fetch the group chat document
	docRef := FirestoreClient.Collection("group_chats").Doc(groupChatID)
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch group chat: %v", err)
	}

	data := docSnap.Data()
	if data == nil {
		return fmt.Errorf("group chat not found")
	}

	// Retrieve participants to check user permissions
	participants := mappers.GetParticipantsGoArray(data, "participants")
	isAuthorized := false
	for _, participant := range participants {
		if participant.UserID == userID && (participant.Role == "owner" || participant.Role == "admin") {
			isAuthorized = true
			break
		}
	}

	if !isAuthorized {
		return fmt.Errorf("only an owner or admin can update group settings")
	}

	// Map updated settings to Firestore format
	settingsFirestore := mappers.MapGroupSettingsGoToFirestore(updatedSettings)

	// Update the Firestore document
	if _, err := docRef.Update(ctx, []firestore.Update{
		{Path: "group_settings", Value: settingsFirestore},
		{Path: "updated_at", Value: time.Now()},
	}); err != nil {
		return fmt.Errorf("failed to update group chat settings: %v", err)
	}

	return nil
}

// func BlockUnblockParticipantService(ctx context.Context, groupChatID, userID, participantID string, block bool) error {
// 	docRef := FirestoreClient.Collection("group_chats").Doc(groupChatID)
// 	docSnap, err := docRef.Get(ctx)
// 	if err != nil {
// 		return fmt.Errorf("failed to fetch group chat: %v", err)
// 	}

// 	data := docSnap.Data()
// 	participants := mappers.GetParticipantsArray(data, "participants")

// 	for i, participant := range participants {
// 		if participant.UserID == userID {
// 			if participant.Blocked == nil {
// 				participant.Blocked = map[string]bool{}
// 			}
// 			participant.Blocked[participantID] = block
// 			participants[i] = participant
// 			break
// 		}
// 	}

// 	firestorePayload := map[string]interface{}{
// 		"participants": mappers.MapParticipantsArrayToFirestore(participants),
// 		"updated_at":   time.Now(),
// 	}

// 	if _, err := docRef.Update(ctx, []firestore.Update{
// 		{Path: "participants", Value: firestorePayload["participants"]},
// 		{Path: "updated_at", Value: firestorePayload["updated_at"]},
// 	}); err != nil {
// 		return fmt.Errorf("failed to update block/unblock status: %v", err)
// 	}

// 	return nil
// }
