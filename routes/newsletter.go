package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

func setupNewsletterRoutes(api fiber.Router, newsletterService *services.NewsletterService) {
	newsletters := api.Group("/newsletters")
	handler := handlers.NewNewsletterHandler(newsletterService)

	// newsletters
	newsletters.Post("/subscribe", handler.SubscribeNewsletter)
	newsletters.Post("/unsubscribe", handler.UnsubscribeNewsletter)
	newsletters.Get("/subscribers", handler.GetAllSubscribers)
	newsletters.Get("/subscribers/count", handler.GetTotalSubscribers)
}
