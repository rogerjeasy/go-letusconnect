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

// CreateAlumniTestimonial handles the creation of a new alumni testimonial
func (h *TestimonialHandler) CreateAlumniTestimonial(c *fiber.Ctx) error {

	userID, err := ExtractAndValidateToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	var testimonial models.AlumniTestimonial
	if err := c.BodyParser(&testimonial); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	ctx := context.Background()
	createdTestimonial, err := h.testimonialService.CreateAlumniTestimonial(ctx, testimonial, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Alumni testimonial created successfully",
		"data":    mappers.MapAlumniTestimonialGoToFrontend(*createdTestimonial),
	})
}

// CreateStudentSpotlight handles the creation of a new student spotlight
func (h *TestimonialHandler) CreateStudentSpotlight(c *fiber.Ctx) error {

	userID, err := ExtractAndValidateToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	var spotlight models.StudentSpotlight
	if err := c.BodyParser(&spotlight); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	ctx := context.Background()
	createdSpotlight, err := h.testimonialService.CreateStudentSpotlight(ctx, spotlight, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Student spotlight created successfully",
		"data":    mappers.MapStudentSpotlightGoToFrontend(*createdSpotlight),
	})
}

// GetTestimonial handles retrieving a testimonial by ID
func (h *TestimonialHandler) GetTestimonial(c *fiber.Ctx) error {
	testimonialID := c.Params("id")
	if testimonialID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Testimonial ID is required",
		})
	}

	ctx := context.Background()
	testimonial, err := h.testimonialService.GetTestimonial(ctx, testimonialID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data":    mappers.MapTestimonialGoToFrontend(*testimonial),
		"message": "Testimonial retrieved successfully",
	})
}

// GetAlumniTestimonial handles retrieving an alumni testimonial by ID
func (h *TestimonialHandler) GetAlumniTestimonial(c *fiber.Ctx) error {
	testimonialID := c.Params("id")
	if testimonialID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Testimonial ID is required",
		})
	}

	ctx := context.Background()
	testimonial, err := h.testimonialService.GetAlumniTestimonial(ctx, testimonialID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data":    mappers.MapAlumniTestimonialGoToFrontend(*testimonial),
		"message": "Alumni testimonial retrieved successfully",
	})
}

// ListTestimonials handles retrieving all published testimonials
func (h *TestimonialHandler) ListTestimonials(c *fiber.Ctx) error {
	ctx := context.Background()
	testimonials, err := h.testimonialService.ListTestimonials(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	frontendTestimonials := make([]map[string]interface{}, 0, len(testimonials))
	for _, testimonial := range testimonials {
		frontendTestimonials = append(frontendTestimonials, mappers.MapTestimonialGoToFrontend(testimonial))
	}

	return c.JSON(fiber.Map{
		"data": frontendTestimonials,
	})
}

// PublishTestimonial handles publishing a testimonial
func (h *TestimonialHandler) PublishTestimonial(c *fiber.Ctx) error {

	userID, err := ExtractAndValidateToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	testimonialID := c.Params("id")
	if testimonialID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Testimonial ID is required",
		})
	}

	// Verify ownership before publishing
	testimonial, err := h.testimonialService.GetTestimonial(c.Context(), testimonialID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if testimonial.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Only the testimonial author can publish it",
		})
	}

	ctx := context.Background()
	err = h.testimonialService.PublishTestimonial(ctx, testimonialID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Testimonial published successfully",
	})
}

// AddLike handles adding a like to a testimonial
func (h *TestimonialHandler) AddLike(c *fiber.Ctx) error {

	_, err := ExtractAndValidateToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	testimonialID := c.Params("id")
	if testimonialID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Testimonial ID is required",
		})
	}

	ctx := context.Background()
	err = h.testimonialService.AddLike(ctx, testimonialID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Like added successfully",
	})
}

// DeleteTestimonial handles deleting a testimonial
func (h *TestimonialHandler) DeleteTestimonial(c *fiber.Ctx) error {

	userID, err := ExtractAndValidateToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	testimonialID := c.Params("id")
	if testimonialID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Testimonial ID is required",
		})
	}

	ctx := context.Background()
	err = h.testimonialService.DeleteTestimonial(ctx, testimonialID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Testimonial deleted successfully",
	})
}
