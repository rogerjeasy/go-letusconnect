package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"github.com/rogerjeasy/go-letusconnect/services"
)

type TestimonialHandler struct {
	testimonialService *services.TestimonialService
	userService        *services.UserService
}

func NewTestimonialHandler(testimonialService *services.TestimonialService, userService *services.UserService) *TestimonialHandler {
	return &TestimonialHandler{
		testimonialService: testimonialService,
		userService:        userService,
	}
}

// CreateTestimonial handles the creation of a new testimonial
func (h *TestimonialHandler) CreateTestimonial(c *fiber.Ctx) error {

	userID, err := ExtractAndValidateToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	var testimonial models.Testimonial
	if err := c.BodyParser(&testimonial); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	ctx := context.Background()
	createdTestimonial, err := h.testimonialService.CreateTestimonial(ctx, testimonial, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Testimonial created successfully",
		"data":    mappers.MapTestimonialGoToFrontend(*createdTestimonial),
	})
}
