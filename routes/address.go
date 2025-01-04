package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

func setupAddressRoutes(api fiber.Router, sc *services.ServiceContainer) error {
	// Validate inputs
	if api == nil {
		return fmt.Errorf("api router cannot be nil")
	}
	if sc == nil {
		return fmt.Errorf("service container cannot be nil")
	}
	if sc.AddressService == nil {
		return fmt.Errorf("address service cannot be nil")
	}

	// Create handler
	handler := handlers.NewAddressHandler(sc.AddressService)
	if handler == nil {
		return fmt.Errorf("failed to create address handler")
	}
	addresses := api.Group("/addresses")

	addresses.Post("/", handler.CreateUserAddress)
	addresses.Put("/:id", handler.UpdateUserAddress)
	addresses.Get("/", handler.GetUserAddress)
	addresses.Delete("/:id", handler.DeleteUserAddress)

	return nil
}
