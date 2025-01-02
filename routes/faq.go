package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

func setupFAQRoutes(api fiber.Router, faqService *services.FAQService) {
	faqs := api.Group("/faqs")
	handler := handlers.NewFAQHandler(faqService)

	// FAQ Routes
	faqs.Get("/", handler.GetAllFAQs)
	faqs.Post("/", handler.CreateFAQ)
	faqs.Put("/:id", handler.UpdateFAQ)
	faqs.Delete("/:id", handler.DeleteFAQ)
}
