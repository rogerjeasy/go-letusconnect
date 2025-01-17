package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

// routes/users.go
func setupUserRoutes(api fiber.Router, sc *services.ServiceContainer) error {
	// Validate inputs
	if api == nil {
		return fmt.Errorf("api router cannot be nil")
	}
	if sc == nil {
		return fmt.Errorf("service container cannot be nil")
	}
	if sc.UserService == nil {
		return fmt.Errorf("user service cannot be nil")
	}

	users := api.Group("/users")
	handler := handlers.NewUserHandler(sc.UserService)

	users.Get("/completion", handler.GetProfileCompletion)
	users.Get("/:uid", handler.GetUser)
	users.Get("/", handler.GetAllUsers)
	users.Put("/", handler.UpdateUser)

	return nil
}
