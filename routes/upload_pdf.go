package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

func SetupPDFRoutes(api fiber.Router, container *services.ServiceContainer) error {
	if container == nil {
		return fiber.NewError(fiber.StatusInternalServerError, "service container is nil")
	}

	if container.UploadPDFService == nil {
		return fiber.NewError(fiber.StatusInternalServerError, "upload PDF service is not initialized")
	}

	// Initialize handler with the service from container
	handler := handlers.NewUploadPDFHandler(container.UploadPDFService)

	// Setup route
	uploadPDF := api.Group("/uploads")
	uploadPDF.Post("/pdf", handler.HandleUploadPDF)

	return nil
}
