package routes

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
	"github.com/rogerjeasy/go-letusconnect/utils"
)

// SetupPusherRoutes sets up the routes for Pusher-related actions
func SetupPusherRoutes(app *fiber.App) {
	app.Post("/api/pusher/auth", handlers.PusherAuth)
	app.Post("/api/trigger", func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Authorization token is required"})
		}

		// Validate the token
		_, err := utils.ValidateToken(strings.TrimPrefix(token, "Bearer "))
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
		}

		var payload struct {
			Message string `json:"message"`
			Channel string `json:"channel"`
			Event   string `json:"event"`
		}

		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request payload"})
		}

		if payload.Channel == "" || payload.Event == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Channel and Event are required"})
		}

		err = services.PusherClient.Trigger(payload.Channel, payload.Event, map[string]string{
			"message": payload.Message,
		})

		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to trigger event"})
		}

		return c.JSON(fiber.Map{"success": "Event triggered successfully"})
	})

}
