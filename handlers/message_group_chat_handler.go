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

type GroupChatHandler struct {
	GroupChatService *services.GroupChatService
	UserService      *services.UserService
}

func NewGroupChatHandler(groupChatService *services.GroupChatService, userService *services.UserService) *GroupChatHandler {
	return &GroupChatHandler{
		GroupChatService: groupChatService,
		UserService:      userService,
	}
}

// CreateGroupChat handles the HTTP request for creating a new group chat
func (h *GroupChatHandler) CreateGroupChatF(c *fiber.Ctx) error {
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
	user, err := h.UserService.GetUserByUID(uid)
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
	groupChat, err := h.GroupChatService.CreateGroupChatService(context.Background(), input)
	if err != nil {
		log.Printf("Failed to create group chat: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create group chat",
		})
	}

	groupChatData, err := h.GroupChatService.GetGroupChatService(context.Background(), groupChat.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Group chat created successfully",
		"data":    groupChatData,
	})
}

// AddParticipantsToGroupChatHandler handles updating the participants list of a group chat
func (h *GroupChatHandler) AddParticipantsToGroupChatHandler(c *fiber.Ctx) error {
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
	if err := h.GroupChatService.AddParticipantsToGroupChat(ctx, groupChatID, projectID, uid, participants); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to add participants to group chat: %v", err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Participants added successfully to the group chat",
	})
}

func (h *GroupChatHandler) GetGroupChat(c *fiber.Ctx) error {
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
	groupChat, err := h.GroupChatService.GetGroupChatService(context.Background(), groupChatId)
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

func (h *GroupChatHandler) GetGroupChatsByProject(c *fiber.Ctx) error {
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
	groupChats, err := h.GroupChatService.GetGroupChatsByProjectService(context.Background(), projectId)
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

func (h *GroupChatHandler) GetMyGroupChats(c *fiber.Ctx) error {
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
	groupChats, err := h.GroupChatService.GetGroupChatsByUserService(context.Background(), uid)
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

func (h *GroupChatHandler) SendMessageHandler(c *fiber.Ctx) error {
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
	user, err := h.UserService.GetUserByUID(uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user details",
		})
	}

	// Call the service
	message, err := h.GroupChatService.SendMessageService(context.Background(), requestData.GroupChatID, uid, user["username"].(string), requestData.Content)
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
	participants, err := h.GroupChatService.GetGroupChatParticipants(context.Background(), requestData.GroupChatID)
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

		unreadCount, err := h.GroupChatService.CountUnreadMessagesService(context.Background(), requestData.GroupChatID, "", participant.UserID)
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

func (h *GroupChatHandler) MarkMessagesAsReadHandler(c *fiber.Ctx) error {
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
	err = h.GroupChatService.MarkMessagesAsReadService(ctx, groupChatID, userID)
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

func (h *GroupChatHandler) CountUnreadMessagesHandler(c *fiber.Ctx) error {
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
	unreadCount, err := h.GroupChatService.CountUnreadMessagesService(ctx, groupChatID, projectID, userID)
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

// CountUnreadGroupMessagesHandler handles the request to count all unread messages across all group chats for a user
func (h *GroupChatHandler) CountUnreadGroupMessagesFromAllChatHandler(c *fiber.Ctx) error {
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

	// Call the service to count all unread messages
	ctx := context.Background()
	unreadCount, err := h.GroupChatService.CountUnreadGroupMessagesFromAllChatService(ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to count unread messages: %v", err),
		})
	}

	// Trigger Pusher event for updating total unread counts
	notificationChannel := "group-total-unread-" + userID
	err = services.PusherClient.Trigger(
		notificationChannel,
		"update-total-unread",
		map[string]interface{}{
			"totalUnreadCount": unreadCount,
		},
	)
	if err != nil {
		log.Printf("Pusher trigger failed: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to trigger total unread count event",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"unreadCount": unreadCount,
	})
}

func (h *GroupChatHandler) RemoveParticipantsFromGroupChatHandler(c *fiber.Ctx) error {
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
	err = h.GroupChatService.RemoveParticipantsFromGroupChatService(ctx, groupChatID, ownerID, requestData.ParticipantIDs)
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

func (h *GroupChatHandler) ReplyToMessageHandler(c *fiber.Ctx) error {
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

	senderDetails, err := h.UserService.GetUserByUID(senderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user details",
		})
	}

	senderName := senderDetails["username"].(string)

	// Call the service to reply to the message
	ctx := context.Background()
	replyMessage, err := h.GroupChatService.ReplyToMessageService(ctx, requestData.GroupChatID, senderID, senderName, requestData.Content, requestData.MessageIDToReply)
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

func (h *GroupChatHandler) AttachFilesToMessageHandler(c *fiber.Ctx) error {
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
	senderDetails, err := h.UserService.GetUserByUID(senderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user details",
		})
	}

	senderName := senderDetails["username"].(string)
	message, err := h.GroupChatService.AttachFilesToMessageService(ctx, groupChatID, senderID, senderName, content, files)
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

func (h *GroupChatHandler) PinMessageHandler(c *fiber.Ctx) error {
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
	err = h.GroupChatService.PinMessageService(ctx, requestData.GroupChatID, userID, requestData.MessageID)
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

func (h *GroupChatHandler) GetPinnedMessagesHandler(c *fiber.Ctx) error {
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
	pinnedMessages, err := h.GroupChatService.GetPinnedMessagesService(ctx, requestData.GroupChatID)
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

func (h *GroupChatHandler) UnpinMessageHandler(c *fiber.Ctx) error {
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
	err = h.GroupChatService.UnpinMessageService(ctx, requestData.GroupChatID, userID, requestData.MessageID)
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

func (h *GroupChatHandler) ReactToMessageHandler(c *fiber.Ctx) error {
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

	err = h.GroupChatService.ReactToMessageService(context.Background(), requestData.GroupChatID, userID, requestData.MessageID, requestData.Reaction)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Reaction added successfully"})
}

func (h *GroupChatHandler) GetMessageReadReceiptsHandler(c *fiber.Ctx) error {
	groupChatID := c.Params("groupChatId")
	messageID := c.Params("messageId")

	receipts, err := h.GroupChatService.GetMessageReadReceiptsService(context.Background(), groupChatID, messageID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"readReceipts": receipts})
}

func (h *GroupChatHandler) SetParticipantRoleHandler(c *fiber.Ctx) error {
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

	err = h.GroupChatService.SetParticipantRoleService(context.Background(), requestData.GroupChatID, userID, requestData.ParticipantID, requestData.NewRole)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Role updated successfully"})
}

func (h *GroupChatHandler) MuteParticipantHandler(c *fiber.Ctx) error {
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

	err = h.GroupChatService.MuteParticipantService(context.Background(), requestData.GroupChatID, userID, requestData.ParticipantID, time.Duration(requestData.Duration)*time.Second)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Participant muted successfully"})
}

func (h *GroupChatHandler) UpdateLastSeenHandler(c *fiber.Ctx) error {
	groupChatID := c.Params("groupChatId")

	userID, err := validateToken(c.Get("Authorization"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	err = h.GroupChatService.UpdateLastSeenService(context.Background(), groupChatID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Last seen updated successfully"})
}

func (h *GroupChatHandler) ArchiveGroupChatHandler(c *fiber.Ctx) error {
	groupChatID := c.Params("groupChatId")

	userID, err := validateToken(c.Get("Authorization"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	err = h.GroupChatService.ArchiveGroupChatService(context.Background(), groupChatID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Group chat archived successfully"})
}

func (h *GroupChatHandler) LeaveGroupHandler(c *fiber.Ctx) error {
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

	err = h.GroupChatService.LeaveGroupService(context.Background(), groupChatID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Successfully left the group"})
}

func (h *GroupChatHandler) CreatePollHandler(c *fiber.Ctx) error {
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

	poll, err := h.GroupChatService.CreatePollService(context.Background(), requestData.GroupChatID, userID, requestData.Poll)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"poll": poll})
}

func (h *GroupChatHandler) ReportMessageHandler(c *fiber.Ctx) error {
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

	err = h.GroupChatService.ReportMessageService(context.Background(), requestData.GroupChatID, userID, requestData.MessageID, requestData.Reason)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Message reported successfully"})
}

// UpdateGroupSettingsHandler handles the HTTP request to update group settings
func (h *GroupChatHandler) UpdateGroupSettingsHandler(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	userID, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	groupChatId := c.Params("groupChatId")
	if groupChatId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "groupChatId is required",
		})
	}

	var requestData map[string]interface{}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	requestDataGo := mappers.MapGroupSettingsFrontendToGo(requestData)

	ctx := context.Background()
	err = h.GroupChatService.UpdateGroupSettingsService(ctx, groupChatId, userID, requestDataGo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to update group settings: %v", err),
		})
	}

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

func (h *GroupChatHandler) DeleteGroupChat(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	chatID := c.Params("id")
	if chatID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Group chat ID is required",
		})
	}

	err = h.GroupChatService.DeleteGroupChatService(context.Background(), chatID, uid)
	if err != nil {
		if strings.Contains(err.Error(), "unauthorized") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Group chat deleted successfully",
	})
}

func (h *GroupChatHandler) DeleteMultipleGroupChats(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	var requestData struct {
		ChatIDs []string `json:"chatIds"`
	}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if len(requestData.ChatIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No chat IDs provided",
		})
	}

	err = h.GroupChatService.DeleteMultipleGroupChatsService(context.Background(), requestData.ChatIDs, uid)
	if err != nil {
		if strings.Contains(err.Error(), "Unauthorized") || strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Group chats deleted successfully",
	})
}
