package handlers

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

type UserConnectionHandler struct {
	connectionService *services.UserConnectionService
}

func NewUserConnectionHandler(connectionService *services.UserConnectionService) *UserConnectionHandler {
	return &UserConnectionHandler{
		connectionService: connectionService,
	}
}

func (h *UserConnectionHandler) GetUserConnections(c *fiber.Ctx) error {
	uid, err := validateToken(strings.TrimPrefix(c.Get("Authorization"), "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	connections, err := h.connectionService.GetUserConnections(context.Background(), uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch connections",
		})
	}

	if connections == nil {
		return c.JSON(fiber.Map{
			"message": "No connections found",
			"data":    nil,
		})
	}

	return c.JSON(mappers.MapConnectionsFirestoreToFrontend(
		mappers.MapConnectionsGoToFirestore(*connections),
	))
}

func (h *UserConnectionHandler) GetUserConnectionsByUID(c *fiber.Ctx) error {
	targetUID := c.Params("uid")
	if targetUID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	connections, err := h.connectionService.GetUserConnections(context.Background(), targetUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch connections",
		})
	}

	if connections == nil {
		return c.JSON(fiber.Map{
			"message": "No connections found",
			"data":    nil,
		})
	}
	frontend := mappers.MapConnectionsFirestoreToFrontend(
		mappers.MapConnectionsGoToFirestore(*connections),
	)
	connectionsMap := frontend["connections"].(map[string]interface{})
	return c.JSON(fiber.Map{
		"connections": connectionsMap,
	})
}

func (h *UserConnectionHandler) SendConnectionRequest(c *fiber.Ctx) error {
	fromUID, err := validateToken(strings.TrimPrefix(c.Get("Authorization"), "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Parse frontend format
	var frontendRequest map[string]interface{}
	if err := c.BodyParser(&frontendRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Convert to Go struct format
	request := mappers.MapConnectionRequestFrontendToGo(frontendRequest)

	err = h.connectionService.SendConnectionRequest(context.Background(), fromUID, request.ToUID, request.Message)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send connection request",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Connection request sent successfully",
	})
}

func (h *UserConnectionHandler) AcceptConnectionRequest(c *fiber.Ctx) error {
	toUID, err := validateToken(strings.TrimPrefix(c.Get("Authorization"), "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	fromUID := c.Params("fromUid")
	if fromUID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "FromUID is required",
		})
	}

	err = h.connectionService.AcceptConnectionRequest(context.Background(), fromUID, toUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to accept connection request",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Connection request accepted",
	})
}

func (h *UserConnectionHandler) RejectConnectionRequest(c *fiber.Ctx) error {
	toUID, err := validateToken(strings.TrimPrefix(c.Get("Authorization"), "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	fromUID := c.Params("fromUid")
	if fromUID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "FromUID is required",
		})
	}

	err = h.connectionService.RejectConnectionRequest(context.Background(), fromUID, toUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to reject connection request",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Connection request rejected",
	})
}

func (h *UserConnectionHandler) RemoveConnection(c *fiber.Ctx) error {
	uid1, err := validateToken(strings.TrimPrefix(c.Get("Authorization"), "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	uid2 := c.Params("uid")
	if uid2 == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Target UID is required",
		})
	}

	err = h.connectionService.RemoveConnection(context.Background(), uid1, uid2)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove connection",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Connection removed successfully",
	})
}
