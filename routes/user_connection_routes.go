package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

func SetupUserConnectionRoutes(app *fiber.App, connectionService *services.UserConnectionService) {
	api := app.Group("/api")
	connectionHandler := handlers.NewUserConnectionHandler(connectionService)

	connections := api.Group("/connections")

	// Get user's connections and pending requests
	connections.Get("/", connectionHandler.GetUserConnections)
	connections.Get("/:uid", connectionHandler.GetUserConnectionsByUID)

	// Connection requests
	connections.Post("/requests", connectionHandler.SendConnectionRequest)
	connections.Put("/requests/:fromUid/accept", connectionHandler.AcceptConnectionRequest)
	connections.Put("/requests/:fromUid/reject", connectionHandler.RejectConnectionRequest)

	// Remove existing connection
	connections.Delete("/:uid", connectionHandler.RemoveConnection)
}
