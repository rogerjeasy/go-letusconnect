package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/middleware"
	"github.com/rogerjeasy/go-letusconnect/services"
)

func SetupAllRoutes(app *fiber.App, services *services.ServiceContainer) {
	api := app.Group("/api/v1")

	// Apply common middleware
	api.Use(middleware.ConfigureCORS())

	// Setup route groups
	setupUserRoutes(api, services.UserService)
	setupNotificationRoutes(api, services.NotificationService)
	setupAuthRoutes(api, services.AuthService)
	setupFAQRoutes(api, services.FAQService)
	setupProjectCoreRoutes(api, services.ProjectCoreService, services.UserService)
	setupProjectCollab(api, services.ProjectService)
	setupDirectMessageRoutes(api, services.MessageService, services.UserService)
	setupGroupChatRoutes(api, services.GroupChatService, services.UserService)
	setupUserConnectionRoutes(api, services.ConnectionService)
	setupAddressRoutes(api, services.AddressService)
	setupNewsletterRoutes(api, services.NewsletterService)
	setupContactUserRoutes(api, services.ContactUsService)
}
