package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"github.com/rogerjeasy/go-letusconnect/services"
)

// NotificationHandler handles HTTP requests related to notifications
type NotificationHandler struct {
	notificationService *services.NotificationService
}

// NewNotificationHandler creates a new NotificationHandler
func NewNotificationHandler(notificationService *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// CreateNotification handles the HTTP request for creating a new notification
func (h *NotificationHandler) CreateNotification(c *fiber.Ctx) error {
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
	var notification models.Notification
	if err := c.BodyParser(&notification); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Set the UserID from the token
	notification.UserID = uid

	// Call service function
	createdNotification, err := h.notificationService.CreateNotification(context.Background(), notification)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create notification",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(mappers.MapNotificationGoToFrontend(*createdNotification))
}

// UpdateNotification handles the HTTP request for updating an existing notification
func (h *NotificationHandler) UpdateNotification(c *fiber.Ctx) error {
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

	notificationID := c.Params("id")
	if notificationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Notification ID is required",
		})
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	updatedNotification, err := h.notificationService.UpdateNotification(context.Background(), notificationID, updates)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update notification",
		})
	}

	return c.Status(fiber.StatusOK).JSON(mappers.MapNotificationGoToFrontend(*updatedNotification))
}

// DeleteNotification handles the HTTP request for deleting a notification
func (h *NotificationHandler) DeleteNotification(c *fiber.Ctx) error {
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

	notificationID := c.Params("id")
	if notificationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Notification ID is required",
		})
	}

	err = h.notificationService.DeleteNotification(context.Background(), notificationID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete notification",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Notification deleted successfully",
	})
}

// GetNotification handles the HTTP request for fetching a single notification
func (h *NotificationHandler) GetNotification(c *fiber.Ctx) error {
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

	notificationID := c.Params("id")
	if notificationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Notification ID is required",
		})
	}

	notification, err := h.notificationService.GetNotification(context.Background(), notificationID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch notification",
		})
	}

	if notification == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Notification not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(mappers.MapNotificationGoToFrontend(*notification))
}

// ListNotifications handles the HTTP request for fetching notifications for a user
func (h *NotificationHandler) ListNotifications(c *fiber.Ctx) error {
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

	limit := c.QueryInt("limit", 20)
	lastNotificationID := c.Query("lastNotificationId")

	notifications, err := h.notificationService.ListNotifications(context.Background(), uid, limit, lastNotificationID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch notifications",
		})
	}

	notificationsResponse := make([]map[string]interface{}, len(notifications))
	for i, notification := range notifications {
		notificationsResponse[i] = mappers.MapNotificationGoToFrontend(notification)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"notifications": notificationsResponse,
	})
}

// MarkNotificationAsRead handles the HTTP request for marking a notification as read
func (h *NotificationHandler) MarkNotificationAsRead(c *fiber.Ctx) error {
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

	notificationID := c.Params("id")
	if notificationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Notification ID is required",
		})
	}

	err = h.notificationService.MarkNotificationAsRead(context.Background(), notificationID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to mark notification as read",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Notification marked as read",
	})
}

// ListTargetedNotifications handles the HTTP request for fetching notifications targeted at a user
func (h *NotificationHandler) ListTargetedNotifications(c *fiber.Ctx) error {
	// Get authorization token
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

	// Get query parameters
	limit := c.QueryInt("limit", 20)
	lastNotificationID := c.Query("lastNotificationId")

	// Fetch notifications
	notifications, err := h.notificationService.ListTargetedNotifications(context.Background(), uid, limit, lastNotificationID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to fetch notifications: %v", err),
		})
	}

	// Map notifications to frontend format
	notificationsResponse := make([]map[string]interface{}, len(notifications))
	for i, notification := range notifications {
		notification.TargetedUsers = nil
		notification.IsRead = notification.ReadStatus[uid]
		notification.ReadStatus = nil
		notificationsResponse[i] = mappers.MapNotificationGoToFrontend(notification)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "Notifications retrieved successfully",
		"notifications": notificationsResponse,
	})
}
