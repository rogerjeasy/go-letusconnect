package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

// routes/users.go
func setupUserRoutes(api fiber.Router, userService *services.UserService) {
	users := api.Group("/users")
	handler := handlers.NewUserHandler(userService)

	users.Get("/completion", handler.GetProfileCompletion)
	users.Get("/:uid", handler.GetUser)
	users.Get("/", handler.GetAllUsers)
	users.Put("/:uid", handler.UpdateUser)
}
