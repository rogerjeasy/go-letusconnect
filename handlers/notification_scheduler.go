package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/models"
	"github.com/rogerjeasy/go-letusconnect/services"
)

type NotificationSchedulerHandler struct {
	notificationService *services.SchedulerNotificationService
}

func NewNotificationSchedulerHandler(notificationService *services.SchedulerNotificationService) *NotificationSchedulerHandler {
	return &NotificationSchedulerHandler{
		notificationService: notificationService,
	}
}

func (h *NotificationSchedulerHandler) ScheduleNotification(c *fiber.Ctx) error {
	var req models.NotificationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.notificationService.ScheduleNotification(c.Context(), &req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Notification scheduled successfully",
	})
}

func (h *NotificationSchedulerHandler) CancelNotification(c *fiber.Ctx) error {
	notificationID := c.Params("id")
	if notificationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Notification ID is required",
		})
	}

	if err := h.notificationService.CancelNotification(c.Context(), notificationID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Notification cancelled successfully",
	})
}

func (h *NotificationSchedulerHandler) GetNotifications(c *fiber.Ctx) error {
	userID := c.Query("userId")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	notifications, err := h.notificationService.GetNotifications(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(notifications)
}
