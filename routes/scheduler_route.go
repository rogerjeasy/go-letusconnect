package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

// setupNotificationRoutes initializes notification-related routes
func setupNotificationSchedulerRoutes(api fiber.Router, sc *services.ServiceContainer) error {
	// Validate dependencies
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
	handler := handlers.NewNotificationSchedulerHandler(sc.SchedulerNotificationService)
	if handler == nil {
		return fmt.Errorf("failed to create notification handler")
	}

	// Notification Routes
	notifications := api.Group("/notifications-scheduler")

	// Schedule and manage notifications
	notifications.Post("/", handler.ScheduleNotification)
	notifications.Delete("/:id", handler.CancelNotification)
	notifications.Get("/", handler.GetNotifications)

	return nil
}
