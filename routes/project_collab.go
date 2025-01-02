package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

// project collaboration routes
func setupProjectCollab(api fiber.Router, projectCollabService *services.ProjectService) {
	projects := api.Group("/projects")
	handler := handlers.NewProjectHandler(projectCollabService)

	// 2. Collaboration Endpoints
	projects.Post("/:id/join", handler.JoinProjectCollab)
	projects.Put("/:id/join-requests/:uid", handler.AcceptRejectJoinRequestCollab)
	projects.Post("/:id/invite", handler.InviteUserCollab)
	projects.Delete("/:id/participants/:uid", handler.RemoveParticipantCollab)

	// 3. Task Endpoints
	projects.Post("/:id/tasks", handler.AddTask)
	projects.Put("/:id/tasks/:taskID", handler.UpdateTask)
	projects.Delete("/:id/tasks/:taskID", handler.DeleteTask)

}
