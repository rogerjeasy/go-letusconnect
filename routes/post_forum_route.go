package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

func setupForumRoutes(api fiber.Router, sc *services.ServiceContainer) error {
	if api == nil {
		return fmt.Errorf("api router cannot be nil")
	}
	if sc == nil {
		return fmt.Errorf("service container cannot be nil")
	}
	if sc.ForumService == nil {
		return fmt.Errorf("forum service cannot be nil")
	}
	if sc.UserService == nil {
		return fmt.Errorf("user service cannot be nil")
	}

	handler := handlers.NewForumHandler(sc.ForumService, sc.UserService)
	if handler == nil {
		return fmt.Errorf("failed to create forum handler")
	}

	forums := api.Group("/forums")

	// Basic CRUD operations
	forums.Post("/", handler.CreateForum)
	forums.Get("/:id", handler.GetForum)
	forums.Put("/:id", handler.UpdateForum)
	forums.Delete("/:id", handler.DeleteForum)

	// Group-related routes
	forums.Get("/group/:groupId", handler.ListForumsByGroup)

	// Posts
	forums.Post("/:id/posts", handler.CreatePost)
	forums.Get("/:id/posts/search", handler.SearchPosts)

	// Comments
	forums.Post("/:forumId/posts/:postId/comments", handler.CreateComment)

	// Reactions
	forums.Post("/:id/reactions", handler.AddReaction)

	// Moderator management
	forums.Post("/:id/moderators", handler.AddModerator)
	forums.Delete("/:id/moderators/:userId", handler.RemoveModerator)

	return nil
}
