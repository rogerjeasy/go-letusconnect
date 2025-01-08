package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

type ChatGPTHandler struct {
	service *services.ChatGPTService
}

func NewChatGPTHandler(service *services.ChatGPTService) *ChatGPTHandler {
	return &ChatGPTHandler{service: service}
}

func (h *ChatGPTHandler) HandleChat(c *fiber.Ctx) error {
	token := c.Get("Authorization")

	var uid string
	var err error

	if token == "" {
		uid = uuid.New().String()
	} else {
		// Validate token and get UID
		uid, err = validateToken(strings.TrimPrefix(token, "Bearer "))
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}
	}

	var request struct {
		Message string `json:"message"`
		ID      string `json:"id"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if request.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Message is required",
		})
	}

	conversation, err := h.service.GenerateResponse(c.Context(), request.Message, uid, request.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := mappers.MapConversationGoToFrontend(*conversation)
	return c.JSON(response)
}

func (h *ChatGPTHandler) GetUserConversations(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token required",
		})
	}

	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	conversations, err := h.service.GetUserConversations(c.Context(), uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Convert conversations to frontend format
	response := make([]map[string]interface{}, len(conversations))
	for i, conv := range conversations {
		response[i] = mappers.MapConversationGoToFrontend(conv)
	}

	return c.JSON(response)
}

func (h *ChatGPTHandler) DeleteConversation(c *fiber.Ctx) error {
	// Get and validate token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token required",
		})
	}

	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	conversationID := c.Params("id")
	if conversationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Conversation ID is required",
		})
	}

	err = h.service.DeleteConversation(c.Context(), conversationID, uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Conversation deleted successfully",
	})
}

func (h *ChatGPTHandler) GetConversation(c *fiber.Ctx) error {
	// Get and validate token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token required",
		})
	}

	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	conversationID := c.Params("id")
	if conversationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Conversation ID is required",
		})
	}

	conversation, err := h.service.GetConversation(c.Context(), conversationID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Verify ownership
	if conversation.UserID != uid {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Unauthorized access to conversation",
		})
	}

	response := mappers.MapConversationGoToFrontend(*conversation)
	return c.JSON(response)
}
