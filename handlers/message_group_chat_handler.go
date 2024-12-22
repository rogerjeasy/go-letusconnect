package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
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
	var requestData map[string]interface{}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Validate mandatory fields
	if _, ok := requestData["name"]; !ok || requestData["name"] == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "name is required",
		})
	}

	// Prepare input for service
	input := services.GroupChatInput{
		ProjectID:      requestData["projectId"].(string),
		Name:           requestData["name"].(string),
		Description:    requestData["description"].(string),
		CreatedByUID:   uid,
		CreatedByName:  user["username"].(string),
		Email:          user["email"].(string),
		ProfilePicture: user["profile_picture"].(string),
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

// GetGroupChat handles fetching a single group chat by ID
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

// GetGroupChatsByProject handles fetching all group chats for a project
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

// GetMyGroupChats handles fetching all group chats for the authenticated user
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

// SendMessageHandler handles sending a message in a group chat
func SendMessageHandler(c *fiber.Ctx) error {
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

	// Parse the request payload
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

	// Call the service to mark messages as read
	ctx := context.Background()
	err = services.MarkMessagesAsReadService(ctx, requestData.GroupChatID, userID)
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

	// Call the service to count unread messages
	ctx := context.Background()
	unreadCount, err := services.CountUnreadMessagesService(ctx, requestData.GroupChatID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to count unread messages: %v", err),
		})
	}

	// Respond with the count
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"unreadCount": unreadCount,
	})
}

func RemoveParticipantFromGroupChatHandler(c *fiber.Ctx) error {
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

	// Parse the request body
	var requestData struct {
		GroupChatID   string `json:"groupChatId"`
		ParticipantID string `json:"participantId"`
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
	if requestData.ParticipantID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "participantId is required",
		})
	}

	// Call the service to remove the participant
	ctx := context.Background()
	err = services.RemoveParticipantFromGroupChatService(ctx, requestData.GroupChatID, ownerID, requestData.ParticipantID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to remove participant: %v", err),
		})
	}

	// Respond with success
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Participant removed successfully",
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
