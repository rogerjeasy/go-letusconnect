package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

func setupChatRoutes(api fiber.Router, sc *services.ServiceContainer) error {
	if api == nil {
		return fmt.Errorf("api router cannot be nil")
	}
	if sc == nil {
		return fmt.Errorf("service container cannot be nil")
	}
	if sc.ChatGPTService == nil {
		return fmt.Errorf("ChatGPTService must be initialized before setting up routes")
	}

	handler := handlers.NewChatGPTHandler(sc.ChatGPTService)
	if handler == nil {
		return fmt.Errorf("failed to create chat handler")
	}

	chat := api.Group("/chat")
	chat.Post("/", handler.HandleChat)
	chat.Get("/history", handler.GetChatHistory)
	chat.Delete("/history/:id", handler.DeleteChatHistory)

	return nil
}
