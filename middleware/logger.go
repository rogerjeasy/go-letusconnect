package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// middleware/logger.go
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		log.Printf(
			"%s %s - %v - %s",
			c.Method(),
			c.Path(),
			time.Since(start),
			c.Get("X-Real-IP"),
		)
		return err
	}
}
