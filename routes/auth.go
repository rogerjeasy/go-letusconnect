// routes/auth.go
package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

func setupAuthRoutes(api fiber.Router, sc *services.ServiceContainer) error {
	// Validate inputs
	if api == nil {
		return fmt.Errorf("api router cannot be nil")
	}
	if sc == nil {
		return fmt.Errorf("service container cannot be nil")
	}
	if sc.AuthService == nil {
		return fmt.Errorf("auth service cannot be nil")
	}

	// Create handler
	handler := handlers.NewAuthHandler(sc.AuthService)
	if handler == nil {
		return fmt.Errorf("failed to create auth handler")
	}

	// Setup auth routes group
	auth := api.Group("/auth")

	// Register routes
	auth.Post("/login", handler.Login)
	auth.Post("/register", handler.Register)
	auth.Get("/session", handler.GetSession)
	auth.Patch("/logout", handler.Logout)

	return nil
}
