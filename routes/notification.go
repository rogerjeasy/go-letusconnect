package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

// routes/notifications.go
func setupNotificationRoutes(api fiber.Router, notificationService *services.NotificationService) {
	notifications := api.Group("/notifications")
	handler := handlers.NewNotificationHandler(notificationService)

	notifications.Get("/targeted", handler.ListTargetedNotifications)
	notifications.Get("/unread-count", handler.GetUnreadNotificationCount)
	notifications.Get("/stats", handler.GetNotificationStats)
	notifications.Patch("/:id", handler.MarkNotificationAsRead)
	notifications.Post("/", handler.CreateNotification)
	notifications.Get("/", handler.ListNotifications)
	notifications.Get("/:id", handler.GetNotification)
	notifications.Put("/:id", handler.UpdateNotification)
	notifications.Delete("/:id", handler.DeleteNotification)
	notifications.Put("/:id/read", handler.MarkNotificationAsRead)
}
