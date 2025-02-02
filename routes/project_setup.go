package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

func setupProjectCoreRoutes(api fiber.Router, sc *services.ServiceContainer) error {
	if api == nil {
		return fmt.Errorf("api router cannot be nil")
	}
	if sc == nil {
		return fmt.Errorf("service container cannot be nil")
	}
	if sc.ProjectCoreService == nil {
		return fmt.Errorf("project core service cannot be nil")
	}
	if sc.UserService == nil {
		return fmt.Errorf("user service cannot be nil")
	}

	handler := handlers.NewProjectHandlerSetup(sc.ProjectCoreService, sc.UserService, sc.GroupChatService)
	if handler == nil {
		return fmt.Errorf("failed to create project handler")
	}

	projects := api.Group("/projects")

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

	return nil

}
