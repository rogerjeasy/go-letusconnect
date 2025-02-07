package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/models"
	"github.com/rogerjeasy/go-letusconnect/services"
)

type LinkedInJobsHandler struct {
	jobsService *services.LinkedInJobsService
}

func NewLinkedInJobsHandler(jobsService *services.LinkedInJobsService) *LinkedInJobsHandler {
	return &LinkedInJobsHandler{
		jobsService: jobsService,
	}
}

func (h *LinkedInJobsHandler) GetAppliedJobs(c *fiber.Ctx) error {
	accessToken := c.Cookies("linkedin_token")
	if accessToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "LinkedIn token not found",
		})
	}

	appliedJobs, err := h.jobsService.GetAppliedJobs(accessToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"jobs": appliedJobs,
	})
}

func (h *LinkedInJobsHandler) StoreJobApplication(c *fiber.Ctx) error {
	ctx := context.Background()
	var jobApp models.LinkedInJobApplication

	if err := c.BodyParser(&jobApp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	userID := c.Locals("userID").(string)

	if err := h.jobsService.StoreJobApplication(ctx, userID, jobApp); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":     "Job application stored successfully",
		"application": jobApp,
	})
}

func (h *LinkedInJobsHandler) GetUserApplications(c *fiber.Ctx) error {
	ctx := context.Background()
	userID := c.Locals("userID").(string)

	applications, err := h.jobsService.GetUserApplications(ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"applications": applications,
	})
}
