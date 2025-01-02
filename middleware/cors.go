// middleware/cors.go
package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func ConfigureCORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "https://letusconnect.vercel.app, http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		AllowMethods:     "GET, HEAD, PUT, PATCH, POST, DELETE, OPTIONS",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length, Authorization",
		MaxAge:           86400,
		AllowOriginsFunc: func(origin string) bool {
			allowedOrigins := map[string]bool{
				"https://letusconnect.vercel.app": true,
				"http://localhost:3000":           true,
			}
			return allowedOrigins[origin]
		},
	})
}
