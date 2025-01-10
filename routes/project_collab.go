package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

func setupProjectCollab(api fiber.Router, sc *services.ServiceContainer) error {
	if api == nil {
		return fmt.Errorf("api router cannot be nil")
	}
	if sc == nil {
		return fmt.Errorf("service container cannot be nil")
	}
	if sc.ProjectService == nil {
		return fmt.Errorf("project service cannot be nil")
	}
	if sc.UserService == nil {
		return fmt.Errorf("user service cannot be nil")
	}

	handler := handlers.NewProjectHandler(sc.ProjectService, sc.UserService)
	if handler == nil {
		return fmt.Errorf("failed to create project handler")
	}

	projects := api.Group("/projects")

	// 2. Collaboration Endpoints
	projects.Post("/:id/join", handler.JoinProjectCollab)
	projects.Put("/:id/join-requests/:uid", handler.AcceptRejectJoinRequestCollab)
	projects.Post("/:id/invite", handler.InviteUserCollab)
	projects.Delete("/:id/participants/:uid", handler.RemoveParticipantCollab)

	// 3. Task Endpoints
	projects.Post("/:id/tasks", handler.AddTask)
	projects.Put("/:id/tasks/:taskID", handler.UpdateTask)
	projects.Delete("/:id/tasks/:taskID", handler.DeleteTask)

	return nil
}
