package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/services"
)

type UploadPDFHandlerToCloudinary struct {
	uploadService *services.UploadPDFService
}

func NewUploadPDFHandler(uploadService *services.UploadPDFService) *UploadPDFHandlerToCloudinary {
	return &UploadPDFHandlerToCloudinary{
		uploadService: uploadService,
	}
}

func (h *UploadPDFHandlerToCloudinary) HandleUploadPDF(c *fiber.Ctx) error {

	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	// Validate token and get user ID
	_, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}
	// Get the file from form
	file, err := c.FormFile("pdf")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error getting the file: " + err.Error(),
		})
	}

	// Call service to handle upload
	response, err := h.uploadService.UploadPDF(file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error uploading file: " + err.Error(),
		})
	}

	return c.JSON(response)
}
