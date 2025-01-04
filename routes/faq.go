package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

func setupFAQRoutes(api fiber.Router, sc *services.ServiceContainer) error {
	if api == nil {
		return fmt.Errorf("api router cannot be nil")
	}
	if sc == nil {
		return fmt.Errorf("service container cannot be nil")
	}
	if sc.FAQService == nil {
		return fmt.Errorf("faq service cannot be nil")
	}

	handler := handlers.NewFAQHandler(sc.FAQService)
	if handler == nil {
		return fmt.Errorf("failed to create faq handler")
	}

	faqs := api.Group("/faqs")
	faqs.Get("/", handler.GetAllFAQs)
	faqs.Post("/", handler.CreateFAQ)
	faqs.Put("/:id", handler.UpdateFAQ)
	faqs.Delete("/:id", handler.DeleteFAQ)

	return nil
}
