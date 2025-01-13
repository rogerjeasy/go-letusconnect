package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

func setupDirectMessageRoutes(api fiber.Router, sc *services.ServiceContainer) error {
	if api == nil {
		return fmt.Errorf("api router cannot be nil")
	}
	if sc == nil {
		return fmt.Errorf("service container cannot be nil")
	}
	if sc.MessageService == nil {
		return fmt.Errorf("message service cannot be nil")
	}
	if sc.UserService == nil {
		return fmt.Errorf("user service cannot be nil")
	}

	handler := handlers.NewMessageHandler(sc.MessageService, sc.UserService)
	if handler == nil {
		return fmt.Errorf("failed to create message handler")
	}

	messages := api.Group("/messages")

	messages.Post("/send", handler.SendMessage)
	messages.Get("/", handler.GetMessages)
	messages.Post("/typing", handler.SendTyping)
	messages.Post("/direct", handler.SendDirectMessage)
	// messages.Post("/group", handlers.SendGroupMessage)
	messages.Get("/direct", handler.GetDirectMessages)
	messages.Get("/unread", handlers.GetUnreadMessagesCount)
	messages.Patch("/mark-as-read", handlers.MarkMessagesAsRead)

	return nil
}
