package handlers

import (
	"context"
	"fmt"
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
	pendingRequestMap := frontend["pendingRequests"].(map[string]interface{})
	return c.JSON(fiber.Map{
		"connections":     connectionsMap,
		"pendingRequests": pendingRequestMap,
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

func (h *UserConnectionHandler) GetConnectionRequests(c *fiber.Ctx) error {
	uid, err := validateToken(strings.TrimPrefix(c.Get("Authorization"), "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	connections, err := h.connectionService.GetUserConnections(context.Background(), uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch connection requests",
		})
	}

	if connections == nil {
		return c.JSON(fiber.Map{
			"pendingRequests": make(map[string]interface{}),
			"sentRequests":    make(map[string]interface{}),
		})
	}

	// // Map the data to frontend format
	// frontend := mappers.MapConnectionsFirestoreToFrontend(
	// 	mappers.MapConnectionsGoToFirestore(*connections),
	// )

	frontend := mappers.MapConnectionsGoToFrontend(*connections)
	// Extract only the requests data
	pendingRequests := frontend["pendingRequests"]
	sentRequests := frontend["sentRequests"]

	return c.JSON(fiber.Map{
		"pendingRequests": pendingRequests,
		"sentRequests":    sentRequests,
		"message":         "Connection requests fetched successfully",
	})
}

func (h *UserConnectionHandler) UpdateRequestStatus(c *fiber.Ctx) error {
	toUID, err := validateToken(strings.TrimPrefix(c.Get("Authorization"), "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	requestID := c.Params("requestId")
	if requestID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Request ID is required",
		})
	}

	// Parse the status update request
	var updateRequest struct {
		Status string `json:"status"` // "accepted" or "rejected"
	}
	if err := c.BodyParser(&updateRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate status
	switch updateRequest.Status {
	case "accepted":
		err = h.connectionService.AcceptConnectionRequest(context.Background(), requestID, toUID)
	case "rejected":
		err = h.connectionService.RejectConnectionRequest(context.Background(), requestID, toUID)
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status. Must be 'accepted' or 'rejected'",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to %s connection request", updateRequest.Status),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("Connection request %s successfully", updateRequest.Status),
	})
}

func (h *UserConnectionHandler) CancelSentRequest(c *fiber.Ctx) error {
	uid, err := validateToken(strings.TrimPrefix(c.Get("Authorization"), "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	toUID := c.Params("toUid")
	if toUID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Target user ID is required",
		})
	}

	err = h.connectionService.CancelSentRequest(context.Background(), uid, toUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to cancel connection request",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Connection request cancelled successfully",
	})
}
