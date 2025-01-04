package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

func setupNotificationRoutes(api fiber.Router, sc *services.ServiceContainer) error {
	// Validate inputs
	if api == nil {
		return fmt.Errorf("api router cannot be nil")
	}
	if sc == nil {
		return fmt.Errorf("service container cannot be nil")
	}
	if sc.NotificationService == nil {
		return fmt.Errorf("notification service cannot be nil")
	}

	// Create handler
	handler := handlers.NewNotificationHandler(sc.NotificationService)
	if handler == nil {
		return fmt.Errorf("failed to create notification handler")
	}

	notifications := api.Group("/notifications")

	notifications.Get("/targeted", handler.ListTargetedNotifications)
	notifications.Get("/unread-count", handler.GetUnreadNotificationCount)
	notifications.Get("/stats", handler.GetNotificationStats)
	notifications.Patch("/:id", handler.MarkNotificationAsRead)
	notifications.Post("/", handler.CreateNotification)
	notifications.Get("/", handler.ListNotifications)
	notifications.Get("/:id", handler.GetNotification)
	notifications.Put("/:id", handler.UpdateNotification)
	notifications.Delete("/:id", handler.DeleteNotification)
	// notifications.Patch("/:id/read", handler.MarkNotificationAsRead)

	return nil
}
