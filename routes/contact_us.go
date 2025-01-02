package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

func setupContactUserRoutes(api fiber.Router, contactUserService *services.ContactUsService) {
	contacts := api.Group("/contact_users")
	handler := handlers.NewContactUsHandler(contactUserService)

	// Contact Us Routes
	contacts.Post("/", handler.CreateContact)
	contacts.Get("/", handler.GetAllContacts)
	contacts.Get("/:id", handler.GetContactByID)
	contacts.Put("/:id", handler.UpdateContactStatus)
}
