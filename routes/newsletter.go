package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

func setupNewsletterRoutes(api fiber.Router, sc *services.ServiceContainer) error {
	if api == nil {
		return fmt.Errorf("api router cannot be nil")
	}
	if sc == nil {
		return fmt.Errorf("service container cannot be nil")
	}
	if sc.NewsletterService == nil {
		return fmt.Errorf("newsletter service cannot be nil")
	}

	handler := handlers.NewNewsletterHandler(sc.NewsletterService)
	if handler == nil {
		return fmt.Errorf("failed to create newsletter handler")
	}

	newsletters := api.Group("/newsletters")

	// newsletters
	newsletters.Post("/subscribe", handler.SubscribeNewsletter)
	newsletters.Post("/unsubscribe", handler.UnsubscribeNewsletter)
	newsletters.Get("/subscribers", handler.GetAllSubscribers)
	newsletters.Get("/subscribers/count", handler.GetTotalSubscribers)

	return nil
}
