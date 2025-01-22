package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
	// "github.com/rogerjeasy/go-letusconnect/middleware"
)

func SetupRoutes(router fiber.Router, notificationService *services.NotificationService) {

	// User Work Experience Routes
	workExperiences := router.Group("/work-experiences")
	workExperiences.Post("/", handlers.CreateUserWorkExperience)
	workExperiences.Put("/:id", handlers.UpdateUserWorkExperience)
	workExperiences.Get("/", handlers.GetUserWorkExperience)
	workExperiences.Delete("/:id", handlers.DeleteUserWorkExperience)

	// Pusher Routes
	SetupPusherRoutes(router)

	// // Media File Routes
	// mediaFiles := api.Group("/media-files")
	// mediaFiles.Post("/upload-images", handlers.UploadImageHandler)
	// mediaFiles.Post("/upload-videos", handlers.UploadVideoHandler)
	// mediaFiles.Post("/upload-pdf", handlers.UploadPDFHandler)

}
