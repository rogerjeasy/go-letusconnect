package routes

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/middleware"
	"github.com/rogerjeasy/go-letusconnect/services"
)

func setupWebSocketRoutes(api fiber.Router, sc *services.ServiceContainer) error {
	if sc.WebSocketService == nil {
		return fmt.Errorf("websocket service cannot be nil")
	}

	ws := api.Group("/ws")

	handler := handlers.NewWebSocketHandler(sc.WebSocketService)

	ws.Use(middleware.WebSocketMiddleware())

	ws.Get("/:id", websocket.New(func(c *websocket.Conn) {
		if err := handler.HandleWebSocket(c); err != nil {
			log.Printf("WebSocket handler error: %v", err)
			c.WriteJSON(fiber.Map{
				"error": "WebSocket error occurred",
			})
			return
		}
	}))

	return nil
}
