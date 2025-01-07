package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
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
	var request struct {
		Message string `json:"message"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	response, err := h.service.GenerateResponse(c.Context(), request.Message, uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(response)
}

func (h *ChatGPTHandler) GetChatHistory(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)

	history, err := h.service.GetChatHistory(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(history)
}

func (h *ChatGPTHandler) DeleteChatHistory(c *fiber.Ctx) error {
	historyID := c.Params("id")
	userID := c.Locals("userId").(string)

	if err := h.service.DeleteChatHistory(c.Context(), historyID, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
