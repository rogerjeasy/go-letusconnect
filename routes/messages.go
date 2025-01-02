package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

// direct message routes
func setupDirectMessageRoutes(api fiber.Router, messageService *services.MessageService, userService *services.UserService) {
	messages := api.Group("/messages")
	handler := handlers.NewMessageHandler(messageService, userService)
	messages.Post("/send", handler.SendMessage)
	messages.Get("/", handler.GetMessages)
	messages.Post("/typing", handler.SendTyping)
	messages.Post("/direct", handler.SendDirectMessage)
	// messages.Post("/group", handlers.SendGroupMessage)
	messages.Get("/direct", handler.GetDirectMessages)
	messages.Get("/unread", handlers.GetUnreadMessagesCount)
	messages.Post("/mark-as-read", handlers.MarkMessagesAsRead)
}
