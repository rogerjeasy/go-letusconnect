package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

func setupTestimonialRoutes(api fiber.Router, sc *services.ServiceContainer) error {
	if api == nil {
		return fmt.Errorf("api router cannot be nil")
	}
	if sc == nil {
		return fmt.Errorf("service container cannot be nil")
	}
	if sc.TestimonialService == nil {
		return fmt.Errorf("testimonial service cannot be nil")
	}
	if sc.UserService == nil {
		return fmt.Errorf("user service cannot be nil")
	}

	handler := handlers.NewTestimonialHandler(sc.TestimonialService, sc.UserService)
	if handler == nil {
		return fmt.Errorf("failed to create testimonial handler")
	}

	testimonials := api.Group("/testimonials")

	// Basic testimonial routes
	testimonials.Post("/", handler.CreateTestimonial)
	testimonials.Get("/", handler.ListTestimonials)
	testimonials.Get("/:id", handler.GetTestimonial)
	testimonials.Delete("/:id", handler.DeleteTestimonial)

	// Alumni testimonial routes
	testimonials.Post("/alumni", handler.CreateAlumniTestimonial)
	testimonials.Get("/alumni/:id", handler.GetAlumniTestimonial)

	// Student spotlight routes
	testimonials.Post("/spotlight", handler.CreateStudentSpotlight)
	// testimonials.Get("/spotlight/:id", handler.GetStudentSpotlight)

	// Publishing and interaction routes
	testimonials.Post("/:id/publish", handler.PublishTestimonial)
	testimonials.Post("/:id/like", handler.AddLike)

	return nil
}
