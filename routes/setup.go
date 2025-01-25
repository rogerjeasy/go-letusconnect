package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/services"
)

func SetupAllRoutes(app *fiber.App, sc *services.ServiceContainer) error {
	// Validate inputs
	if app == nil {
		return fmt.Errorf("fiber app cannot be nil")
	}
	if sc == nil {
		return fmt.Errorf("service container cannot be nil")
	}

	api := app.Group("/api/v1")

	// Apply common middleware
	// api.Use(middleware.ConfigureCORS())

	routeSetups := []struct {
		name string
		fn   func(fiber.Router, *services.ServiceContainer) error
	}{
		{"user", setupUserRoutes},
		{"notification", setupNotificationRoutes},
		{"auth", setupAuthRoutes},
		{"faq", setupFAQRoutes},
		{"projectCore", setupProjectCoreRoutes},
		{"projectCollab", setupProjectCollab},
		{"directMessage", setupDirectMessageRoutes},
		{"groupChat", setupGroupChatRoutes},
		{"userConnection", setupUserConnectionRoutes},
		{"address", setupAddressRoutes},
		{"newsletter", setupNewsletterRoutes},
		{"contactUser", setupContactUserRoutes},
		{"chat", setupChatRoutes},
		{"webSocket", SetupWebSocketRoutes},
		{"schoolExperience", setupUserSchoolExperienceRoutes},
		{"group", setupGroupRoutes},
	}

	for _, setup := range routeSetups {
		if err := setup.fn(api, sc); err != nil {
			return fmt.Errorf("failed to setup %s routes: %w", setup.name, err)
		}
	}

	return nil
}
