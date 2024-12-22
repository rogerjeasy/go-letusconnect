package handlers

import (
	"context"
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
