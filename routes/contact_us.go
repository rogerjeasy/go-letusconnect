package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

func setupContactUserRoutes(api fiber.Router, sc *services.ServiceContainer) error {
	if api == nil {
		return fmt.Errorf("api router cannot be nil")
	}
	if sc == nil {
		return fmt.Errorf("service container cannot be nil")
	}
	if sc.ContactUsService == nil {
		return fmt.Errorf("contact user service cannot be nil")
	}

	handler := handlers.NewContactUsHandler(sc.ContactUsService)
	if handler == nil {
		return fmt.Errorf("failed to create contact user handler")
	}

	contacts := api.Group("/contact_users")

	// Contact Us Routes
	contacts.Post("/", handler.CreateContact)
	contacts.Get("/", handler.GetAllContacts)
	contacts.Get("/:id", handler.GetContactByID)
	contacts.Put("/:id", handler.UpdateContactStatus)

	return nil
}
