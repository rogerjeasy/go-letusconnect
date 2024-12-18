package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/services"
)

// SetupPusherRoutes sets up the routes for Pusher-related actions
func SetupPusherRoutes(app *fiber.App) {
	app.Post("/api/trigger", func(c *fiber.Ctx) error {
		var payload struct {
			Message string `json:"message"`
			Channel string `json:"channel"`
			Event   string `json:"event"`
		}

		if err := c.BodyParser(&payload); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request payload"})
		}

		err := services.PusherClient.Trigger(payload.Channel, payload.Event, map[string]string{
			"message": payload.Message,
		})

		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to trigger event"})
		}

		return c.JSON(fiber.Map{"success": "Event triggered successfully"})
	})
}
