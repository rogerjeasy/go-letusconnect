package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

// project management routes
func setupProjectCoreRoutes(api fiber.Router, projectCoreService *services.ProjectCoreService, userService *services.UserService) {
	projects := api.Group("/projects")
	handler := handlers.NewProjectHandlerSetup(projectCoreService, userService)

	// Project Management Routes
	projects.Get("/owner", handler.GetOwnerProjects)
	projects.Get("/participation", handler.GetParticipationProjects)
	projects.Get("/public", handler.GetAllPublicProjects)
	// projects.Use(middleware.AuthMiddleware)
	projects.Post("/", handler.CreateProject)
	// projects.Get("/", handlers.GetAllProjects)
	projects.Get("/:id", handler.GetProject)
	projects.Put("/:id", handler.UpdateProject)
	projects.Delete("/:id", handler.DeleteProject)

}
