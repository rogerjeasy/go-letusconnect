package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

func setupUserConnectionRoutes(api fiber.Router, connectionService *services.UserConnectionService) {
	connectionHandler := handlers.NewUserConnectionHandler(connectionService)

	// Main connections resource
	connections := api.Group("/connections")

	// Core connection endpoints
	connections.Get("/", connectionHandler.GetUserConnections)
	connections.Get("/:uid", connectionHandler.GetUserConnectionsByUID)
	connections.Delete("/:uid", connectionHandler.RemoveConnection)

	// Connection requests as a sub-resource
	requests := connections.Group("/requests")
	requests.Get("/:uid", connectionHandler.GetConnectionRequests)
	requests.Post("/", connectionHandler.SendConnectionRequest)

	// Request actions (keeping separate accept/reject endpoints)
	requests.Put("/:fromUid/accept", connectionHandler.AcceptConnectionRequest)
	requests.Put("/:fromUid/reject", connectionHandler.RejectConnectionRequest)
	requests.Delete("/:toUid", connectionHandler.CancelSentRequest)
}
