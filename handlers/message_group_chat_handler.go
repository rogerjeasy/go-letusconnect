package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"github.com/rogerjeasy/go-letusconnect/services"
)

// CreateGroupChat handles the HTTP request for creating a new group chat
func CreateGroupChatF(c *fiber.Ctx) error {
	// Extract the Authorization token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	// Validate token and get UID
	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Fetch the user's details
	user, err := services.GetUserByUID(uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user details",
		})
	}

	// Parse the request payload
	var requestData struct {
		ProjectID    string                   `json:"projectId"`
		Name         string                   `json:"name"`
		Description  string                   `json:"description"`
		Participants []map[string]interface{} `json:"participants"`
	}

	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Validate mandatory fields
	if requestData.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "name is required",
		})
	}

	// Map participants from frontend format to backend models.Participant
	var participants []models.Participant
	for _, participant := range requestData.Participants {
		participants = append(participants, models.Participant{
			UserID:         participant["userId"].(string),
			Username:       participant["username"].(string),
			Role:           participant["role"].(string),
			JoinedAt:       time.Now(), // Set joined_at to the current time
			Email:          participant["email"].(string),
			ProfilePicture: participant["profilePicture"].(string),
		})
	}

	// Prepare input for service
	input := services.GroupChatInput{
		ProjectID:      requestData.ProjectID,
		Name:           requestData.Name,
		Description:    requestData.Description,
		CreatedByUID:   uid,
		CreatedByName:  user["username"].(string),
		Email:          user["email"].(string),
		ProfilePicture: user["profile_picture"].(string),
		Participants:   participants,
	}

	// Call service function
	groupChat, err := services.CreateGroupChatService(context.Background(), input)
	if err != nil {
		log.Printf("Failed to create group chat: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create group chat",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "Group chat created successfully",
		"groupChatId": groupChat.ID,
	})
}

// AddParticipantsToGroupChatHandler handles updating the participants list of a group chat
func AddParticipantsToGroupChatHandler(c *fiber.Ctx) error {
	// Extract the Authorization token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	// Validate token and get UID
	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Extract groupChatID or projectID from the URL
	groupChatID := c.Params("groupChatId")
	projectID := c.Params("projectId")

	if groupChatID == "" && projectID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Either groupChatId or projectId must be provided",
		})
	}

	// Parse request body for participants
	var requestData struct {
		Participants []map[string]interface{} `json:"participants"`
	}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if len(requestData.Participants) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Participants list cannot be empty",
		})
	}

	// Map participants from frontend format to backend models.Participant
	var participants []models.Participant
	for _, participant := range requestData.Participants {
		participants = append(participants, models.Participant{
			UserID:         participant["userId"].(string),
			Username:       participant["username"].(string),
			Role:           participant["role"].(string),
			JoinedAt:       time.Now(), // Set joined_at to the current time
			Email:          participant["email"].(string),
			ProfilePicture: participant["profilePicture"].(string),
		})
	}

	ctx := context.Background()

	// Call the service function to add participants
	if err := services.AddParticipantsToGroupChat(ctx, groupChatID, projectID, uid, participants); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to add participants to group chat: %v", err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Participants added successfully to the group chat",
	})
}

func GetGroupChat(c *fiber.Ctx) error {
	// Extract the Authorization token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	// Validate token and get UID
	_, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Get group chat ID from params
	groupChatId := c.Params("id")
	if groupChatId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Group chat ID is required",
		})
	}

	// Fetch group chat using service
	groupChat, err := services.GetGroupChatService(context.Background(), groupChatId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Group chat fetched successfully",
		"data":    groupChat,
	})
}

func GetGroupChatsByProject(c *fiber.Ctx) error {
	// Extract the Authorization token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	// Validate token and get UID
	_, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Get project ID from params
	projectId := c.Params("projectId")
	if projectId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Project ID is required",
		})
	}

	// Fetch group chats using service
	groupChats, err := services.GetGroupChatsByProjectService(context.Background(), projectId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Group chats fetched successfully",
		"data":    groupChats,
	})
}

func GetMyGroupChats(c *fiber.Ctx) error {
	// Extract the Authorization token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	// Validate token and get UID
	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Fetch group chats using service
	groupChats, err := services.GetGroupChatsByUserService(context.Background(), uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Group chats fetched successfully",
		"data":    groupChats,
	})
}

func SendMessageHandler(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	// Validate token and get UID
	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	var requestData struct {
		GroupChatID string `json:"groupChatId"`
		Content     string `json:"content"`
	}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if requestData.GroupChatID == "" || requestData.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "GroupChatID and content are required",
		})
	}

	// Fetch the user's details
	user, err := services.GetUserByUID(uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user details",
		})
	}

	// Call the service
	message, err := services.SendMessageService(context.Background(), requestData.GroupChatID, uid, user["username"].(string), requestData.Content)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Trigger Pusher event for the group chat
	channelName := "group-messages-" + requestData.GroupChatID
	err = services.PusherClient.Trigger(
		channelName,
		"new-group-message",
		mappers.MapBaseMessageGoToFrontend(*message),
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to trigger group chat event",
		})
	}

	// Notify participants via their notification channels
	participants, err := services.GetGroupChatParticipants(context.Background(), requestData.GroupChatID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch group participants",
		})
	}

	for _, participant := range participants {
		if participant.UserID == uid {
			continue // Skip notifying the sender
		}

		// Notify participant about the new message
		notificationChannelNewMessage := "user-notifications-new-msg-" + participant.UserID
		err = services.PusherClient.Trigger(
			notificationChannelNewMessage,
			"new-unread-message",
			map[string]string{
				"groupChatId": requestData.GroupChatID,
				"senderName":  user["username"].(string),
				"content":     requestData.Content,
				"messageId":   message.ID,
			},
		)
		if err != nil {
			fmt.Printf("Failed to notify participant %s: %v", participant.Username, err)
		}

		notificationChannel := "user-notifications-" + participant.UserID
		err = services.PusherClient.Trigger(
			notificationChannel,
			"update-unread-count",
			map[string]string{
				"groupChatId": requestData.GroupChatID,
				"senderName":  user["username"].(string),
				"content":     requestData.Content,
				"messageId":   message.ID,
			},
		)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to trigger notification event for participant " + participant.Username,
			})
		}

		unreadCount, err := services.CountUnreadMessagesService(context.Background(), requestData.GroupChatID, "", participant.UserID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to count unread messages",
			})
		}

		unreadCountNotificationChannel := "group-unread-counts-" + participant.UserID
		err = services.PusherClient.Trigger(
			unreadCountNotificationChannel,
			"update-unread-count",
			map[string]interface{}{
				"groupChatId": requestData.GroupChatID,
				"unreadCount": unreadCount,
			},
		)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to trigger unread count event for participant " + participant.Username,
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Message sent successfully",
		"data":    message,
	})
}

func MarkMessagesAsReadHandler(c *fiber.Ctx) error {
	// Extract the Authorization token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	// Validate token and get user ID
	userID, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Get groupChatId from path parameters
	groupChatID := c.Params("groupChatId")
	if groupChatID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "groupChatId is required",
		})
	}

	// Call the service to mark messages as read
	ctx := context.Background()
	err = services.MarkMessagesAsReadService(ctx, groupChatID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to mark messages as read: %v", err),
		})
	}

	// Respond with success
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Messages marked as read successfully",
	})
}

func CountUnreadMessagesHandler(c *fiber.Ctx) error {
	// Extract the Authorization token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	// Validate token and get user ID
	userID, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Extract query parameters
	groupChatID := c.Query("groupChatId")
	projectID := c.Query("projectId")

	if groupChatID == "" && projectID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Either groupChatId or projectId is required",
		})
	}

	// Call the service to count unread messages
	ctx := context.Background()
	unreadCount, err := services.CountUnreadMessagesService(ctx, groupChatID, projectID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to count unread messages: %v", err),
		})
	}

	// Trigger Pusher event for updating unread counts
	notificationChannel := "group-unread-counts-" + userID
	err = services.PusherClient.Trigger(
		notificationChannel,
		"update-unread-count",
		map[string]interface{}{
			"groupChatId": groupChatID,
			"unreadCount": unreadCount,
		},
	)
	if err != nil {
		log.Printf("Pusher trigger failed: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to trigger unread count event",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"unreadCount": unreadCount,
	})
}

func RemoveParticipantsFromGroupChatHandler(c *fiber.Ctx) error {
	// Extract the Authorization token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	// Validate token and get user ID
	ownerID, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Extract groupChatId from URL parameter
	groupChatID := c.Params("groupChatId")
	if groupChatID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "groupChatId is required",
		})
	}

	// Parse the request body
	var requestData struct {
		ParticipantIDs []string `json:"participantIds"`
	}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Validate participant IDs
	if len(requestData.ParticipantIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "At least one participant ID is required",
		})
	}

	// Call the service to remove the participants
	ctx := context.Background()
	err = services.RemoveParticipantsFromGroupChatService(ctx, groupChatID, ownerID, requestData.ParticipantIDs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to remove participants: %v", err),
		})
	}

	// Respond with success
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Participants removed successfully",
	})
}

func ReplyToMessageHandler(c *fiber.Ctx) error {
	// Extract the Authorization token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	// Validate token and get user ID
	senderID, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Parse the request body
	var requestData struct {
		GroupChatID      string `json:"groupChatId"`
		Content          string `json:"content"`
		MessageIDToReply string `json:"messageIdToReply"`
	}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Validate required fields
	if requestData.GroupChatID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "groupChatId is required",
		})
	}
	if requestData.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "content is required",
		})
	}
	if requestData.MessageIDToReply == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "messageIdToReply is required",
		})
	}

	senderDetails, err := services.GetUserByUID(senderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user details",
		})
	}

	senderName := senderDetails["username"].(string)

	// Call the service to reply to the message
	ctx := context.Background()
	replyMessage, err := services.ReplyToMessageService(ctx, requestData.GroupChatID, senderID, senderName, requestData.Content, requestData.MessageIDToReply)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to reply to the message: %v", err),
		})
	}

	// Respond with the new reply message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": replyMessage,
	})
}

func AttachFilesToMessageHandler(c *fiber.Ctx) error {
	// Extract the Authorization token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	// Validate token and get user ID
	senderID, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Parse form fields
	groupChatID := c.FormValue("groupChatId")
	content := c.FormValue("content")
	if groupChatID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "groupChatId is required",
		})
	}

	// Get uploaded files
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse multipart form",
		})
	}
	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "At least one file is required",
		})
	}

	// Call the service to attach files to the message
	ctx := context.Background()
	senderDetails, err := services.GetUserByUID(senderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user details",
		})
	}

	senderName := senderDetails["username"].(string)
	message, err := services.AttachFilesToMessageService(ctx, groupChatID, senderID, senderName, content, files)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to attach files to message: %v", err),
		})
	}

	// Respond with the new message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": message,
	})
}

func PinMessageHandler(c *fiber.Ctx) error {
	// Extract the Authorization token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	// Validate token and get user ID
	userID, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Parse the request body
	var requestData struct {
		GroupChatID string `json:"groupChatId"`
		MessageID   string `json:"messageId"`
	}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Validate required fields
	if requestData.GroupChatID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "groupChatId is required",
		})
	}
	if requestData.MessageID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "messageId is required",
		})
	}

	// Call the service to pin the message
	ctx := context.Background()
	err = services.PinMessageService(ctx, requestData.GroupChatID, userID, requestData.MessageID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to pin message: %v", err),
		})
	}

	// Respond with success
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Message pinned successfully",
	})
}

func GetPinnedMessagesHandler(c *fiber.Ctx) error {
	// Extract the Authorization token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	// Validate token and get user ID
	_, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Parse the request body
	var requestData struct {
		GroupChatID string `json:"groupChatId"`
	}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Validate required fields
	if requestData.GroupChatID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "groupChatId is required",
		})
	}

	// Call the service to get pinned messages
	ctx := context.Background()
	pinnedMessages, err := services.GetPinnedMessagesService(ctx, requestData.GroupChatID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to fetch pinned messages: %v", err),
		})
	}

	// Respond with pinned messages
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"pinnedMessages": pinnedMessages,
	})
}

func UnpinMessageHandler(c *fiber.Ctx) error {
	// Extract the Authorization token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	// Validate token and get user ID
	userID, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Parse the request body
	var requestData struct {
		GroupChatID string `json:"groupChatId"`
		MessageID   string `json:"messageId"`
	}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Validate required fields
	if requestData.GroupChatID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "groupChatId is required",
		})
	}
	if requestData.MessageID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "messageId is required",
		})
	}

	// Call the service to unpin the message
	ctx := context.Background()
	err = services.UnpinMessageService(ctx, requestData.GroupChatID, userID, requestData.MessageID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to unpin message: %v", err),
		})
	}

	// Respond with success
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Message unpinned successfully",
	})
}

func ReactToMessageHandler(c *fiber.Ctx) error {
	var requestData struct {
		GroupChatID string `json:"groupChatId"`
		MessageID   string `json:"messageId"`
		Reaction    string `json:"reaction"`
	}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	userID, err := validateToken(c.Get("Authorization"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	err = services.ReactToMessageService(context.Background(), requestData.GroupChatID, userID, requestData.MessageID, requestData.Reaction)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Reaction added successfully"})
}

func GetMessageReadReceiptsHandler(c *fiber.Ctx) error {
	groupChatID := c.Params("groupChatId")
	messageID := c.Params("messageId")

	receipts, err := services.GetMessageReadReceiptsService(context.Background(), groupChatID, messageID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"readReceipts": receipts})
}

func SetParticipantRoleHandler(c *fiber.Ctx) error {
	var requestData struct {
		GroupChatID   string `json:"groupChatId"`
		ParticipantID string `json:"participantId"`
		NewRole       string `json:"newRole"`
	}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	userID, err := validateToken(c.Get("Authorization"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	err = services.SetParticipantRoleService(context.Background(), requestData.GroupChatID, userID, requestData.ParticipantID, requestData.NewRole)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Role updated successfully"})
}

func MuteParticipantHandler(c *fiber.Ctx) error {
	var requestData struct {
		GroupChatID   string `json:"groupChatId"`
		ParticipantID string `json:"participantId"`
		Duration      int64  `json:"duration"` // Duration in seconds
	}

	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	userID, err := validateToken(c.Get("Authorization"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	err = services.MuteParticipantService(context.Background(), requestData.GroupChatID, userID, requestData.ParticipantID, time.Duration(requestData.Duration)*time.Second)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Participant muted successfully"})
}

func UpdateLastSeenHandler(c *fiber.Ctx) error {
	groupChatID := c.Params("groupChatId")

	userID, err := validateToken(c.Get("Authorization"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	err = services.UpdateLastSeenService(context.Background(), groupChatID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Last seen updated successfully"})
}

func ArchiveGroupChatHandler(c *fiber.Ctx) error {
	groupChatID := c.Params("groupChatId")

	userID, err := validateToken(c.Get("Authorization"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	err = services.ArchiveGroupChatService(context.Background(), groupChatID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Group chat archived successfully"})
}

func LeaveGroupHandler(c *fiber.Ctx) error {
	groupChatID := c.Params("groupChatId")

	// Extract the Authorization token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	// Validate token and get UID
	userID, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	err = services.LeaveGroupService(context.Background(), groupChatID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Successfully left the group"})
}

func CreatePollHandler(c *fiber.Ctx) error {
	var requestData struct {
		GroupChatID string      `json:"groupChatId"`
		Poll        models.Poll `json:"poll"`
	}

	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	userID, err := validateToken(c.Get("Authorization"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	poll, err := services.CreatePollService(context.Background(), requestData.GroupChatID, userID, requestData.Poll)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"poll": poll})
}

func ReportMessageHandler(c *fiber.Ctx) error {
	var requestData struct {
		GroupChatID string `json:"groupChatId"`
		MessageID   string `json:"messageId"`
		Reason      string `json:"reason"`
	}

	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	userID, err := validateToken(c.Get("Authorization"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	err = services.ReportMessageService(context.Background(), requestData.GroupChatID, userID, requestData.MessageID, requestData.Reason)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Message reported successfully"})
}

// UpdateGroupSettingsHandler handles the HTTP request to update group settings
func UpdateGroupSettingsHandler(c *fiber.Ctx) error {
	// Extract the Authorization token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	// Validate token and get user ID
	userID, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Parse the request payload
	var requestData struct {
		GroupChatID   string               `json:"groupChatId"`
		GroupSettings models.GroupSettings `json:"groupSettings"`
	}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Validate required fields
	if requestData.GroupChatID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "groupChatId is required",
		})
	}

	// Call the service to update the group settings
	ctx := context.Background()
	err = services.UpdateGroupSettingsService(ctx, requestData.GroupChatID, userID, requestData.GroupSettings)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to update group settings: %v", err),
		})
	}

	// Respond with success
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Group settings updated successfully",
	})
}

// func BlockUnblockParticipantHandler(c *fiber.Ctx) error {
// 	var requestData struct {
// 		GroupChatID   string `json:"groupChatId"`
// 		ParticipantID string `json:"participantId"`
// 		Block         bool   `json:"block"`
// 	}

// 	if err := c.BodyParser(&requestData); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
// 	}

// 	userID, err := validateToken(c.Get("Authorization"))
// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
// 	}

// 	err = services.BlockUnblockParticipantService(context.Background(), requestData.GroupChatID, userID, requestData.ParticipantID, requestData.Block)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	return c.JSON(fiber.Map{"message": "Participant block/unblock status updated successfully"})
// }
