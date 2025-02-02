package routes

// import (
// 	"fmt"

// 	"github.com/gofiber/fiber/v2"
// 	// "github.com/gofiber/websocket/v2"
// 	// "github.com/rogerjeasy/go-letusconnect/middleware"
// 	"github.com/rogerjeasy/go-letusconnect/middleware"
// 	"github.com/rogerjeasy/go-letusconnect/services"
// 	wsHandler "github.com/rogerjeasy/go-letusconnect/websocket"
// )

// func SetupWebSocketRoutes(api fiber.Router, sc *services.ServiceContainer) error {
// 	if api == nil {
// 		return fmt.Errorf("api router cannot be nil")
// 	}
// 	if sc == nil {
// 		return fmt.Errorf("service container cannot be nil")
// 	}
// 	if sc.WebSocketService == nil {
// 		return fmt.Errorf("websocket service cannot be nil")
// 	}

// 	// Create WebSocket group under the API router
// 	ws := api.Group("/ws")
// 	ws.Use(middleware.ConfigureWebSocketCORS())

// 	ws.Use(func(c *fiber.Ctx) error {
// 		if c.Get("Upgrade") == "websocket" {
// 			c.Set("Connection", "Upgrade")
// 			c.Set("Upgrade", "websocket")
// 			return c.Next()
// 		}
// 		return fiber.ErrUpgradeRequired
// 	})

// 	// Add WebSocket middleware to handle upgrade
// 	// ws.Use(func(c *fiber.Ctx) error {
// 	// 	// IsWebSocketUpgrade returns true if the client
// 	// 	// requested upgrade to the WebSocket protocol
// 	// 	if websocket.IsWebSocketUpgrade(c) {
// 	// 		c.Locals("allowed", true)
// 	// 		return c.Next()
// 	// 	}
// 	// 	return fiber.ErrUpgradeRequired
// 	// })

// 	// Add CORS middleware specific to WebSocket
// 	// ws.Use(middleware.ConfigureWebSocketCORS())

// 	// Create WebSocket handler
// 	handler := wsHandler.NewHandler(sc.WebSocketService.GetManager())

// 	// WebSocket connection endpoints
// 	ws.Get("/", handler.HandleWebSocket)              // Main WebSocket endpoint
// 	ws.Get("/chat", handler.HandleWebSocket)          // Chat-specific endpoint
// 	ws.Get("/notifications", handler.HandleWebSocket) // Notifications endpoint

// 	return nil
// }
