package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

func setupAddressRoutes(api fiber.Router, addressService *services.AddressService) {
	addresses := api.Group("/addresses")
	handler := handlers.NewAddressHandler(addressService)

	addresses.Post("/", handler.CreateUserAddress)
	addresses.Put("/:id", handler.UpdateUserAddress)
	addresses.Get("/", handler.GetUserAddress)
	addresses.Delete("/:id", handler.DeleteUserAddress)
}
