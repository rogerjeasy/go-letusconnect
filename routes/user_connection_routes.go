package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

func setupUserConnectionRoutes(api fiber.Router, sc *services.ServiceContainer) error {
	if api == nil {
		return fmt.Errorf("api router cannot be nil")
	}
	if sc == nil {
		return fmt.Errorf("service container cannot be nil")
	}
	if sc.ConnectionService == nil {
		return fmt.Errorf("connection service cannot be nil")
	}

	connectionHandler := handlers.NewUserConnectionHandler(sc.ConnectionService)
	if connectionHandler == nil {
		return fmt.Errorf("failed to create connection handler")
	}

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

	return nil
}
