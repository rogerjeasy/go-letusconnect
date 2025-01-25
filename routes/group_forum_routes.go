package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

func setupGroupRoutes(api fiber.Router, sc *services.ServiceContainer) error {
	if api == nil {
		return fmt.Errorf("api router cannot be nil")
	}
	if sc == nil {
		return fmt.Errorf("service container cannot be nil")
	}
	if sc.GroupService == nil {
		return fmt.Errorf("group service cannot be nil")
	}
	if sc.UserService == nil {
		return fmt.Errorf("user service cannot be nil")
	}

	handler := handlers.NewGroupHandler(sc.GroupService, sc.UserService)
	if handler == nil {
		return fmt.Errorf("failed to create group handler")
	}

	groups := api.Group("/group-forums")

	groups.Post("/", handler.CreateGroup)
	groups.Get("/my-groups", handler.ListGroupsByUser)
	groups.Get("/search", handler.SearchGroups)
	groups.Get("/:id", handler.GetGroup)
	groups.Put("/:id", handler.UpdateGroup)
	groups.Delete("/:id", handler.DeleteGroup)
	groups.Get("/", handler.ListGroups)

	// Members
	groups.Post("/:id/members", handler.AddMember)
	groups.Delete("/:id/members/:userId", handler.RemoveMember)

	// Images
	groups.Post("/:id/image", handler.UploadGroupImage)

	// Events
	groups.Post("/:id/events", handler.AddEvent)
	groups.Delete("/:id/events/:eventId", handler.RemoveEvent)

	// Resources
	groups.Post("/:id/resources", handler.AddResource)
	groups.Delete("/:id/resources/:resourceId", handler.RemoveResource)

	// Settings
	groups.Put("/:id/settings", handler.UpdateGroupSettings)

	return nil
}
