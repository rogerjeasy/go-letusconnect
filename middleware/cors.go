// middleware/cors.go
package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// Common allowed origins
var allowedOrigins = map[string]bool{
	"https://letusconnect.vercel.app": true,
	"http://localhost:3000":           true,
	"ws://localhost:3000":             true,
	"wss://letusconnect.vercel.app":   true,
}

// ConfigureCORS returns CORS middleware configuration for HTTP requests
func ConfigureCORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     strings.Join(getOriginsSlice(), ", "),
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		AllowMethods:     "GET, HEAD, PUT, PATCH, POST, DELETE, OPTIONS",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length, Authorization",
		MaxAge:           86400,
		AllowOriginsFunc: validateOrigin,
	})
}

// ConfigureWebSocketCORS returns CORS middleware configuration specifically for WebSocket connections
func ConfigureWebSocketCORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     strings.Join(getOriginsSlice(), ", "),
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, Sec-WebSocket-Protocol, Sec-WebSocket-Version, Sec-WebSocket-Key, Sec-WebSocket-Extensions",
		AllowMethods:     "GET, OPTIONS",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length, Authorization, Sec-WebSocket-Accept",
		MaxAge:           300,
		AllowOriginsFunc: validateOrigin,
	})
}

// validateOrigin checks if the origin is allowed
func validateOrigin(origin string) bool {
	// Convert HTTP to WS and HTTPS to WSS for WebSocket origins
	wsOrigin := strings.Replace(strings.Replace(origin, "http:", "ws:", 1), "https:", "wss:", 1)
	return allowedOrigins[origin] || allowedOrigins[wsOrigin]
}

// getOriginsSlice returns a slice of all allowed origins
func getOriginsSlice() []string {
	origins := make([]string, 0, len(allowedOrigins))
	for origin := range allowedOrigins {
		origins = append(origins, origin)
	}
	return origins
}
