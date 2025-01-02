package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

// routes/auth.go
func setupAuthRoutes(router fiber.Router, authService *services.AuthService) {
	auth := router.Group("/auth")
	handler := handlers.NewAuthHandler(authService)

	auth.Post("/login", handler.Login)
	auth.Post("/register", handler.Register)
	auth.Get("/session", handler.GetSession)
	auth.Patch("/logout", handler.Logout)
}
