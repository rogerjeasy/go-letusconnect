package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

const (
	errInvalidPayload     = "Invalid request payload"
	errExperienceNotFound = "School experience not found"
	errFetchExperience    = "Failed to fetch school experience"
	errParseExperience    = "Failed to parse school experience data"
	errAddUniversity      = "Failed to add university"
	errAddUniversities    = "Failed to add universities"
	msgAddSuccess         = "University added successfully"
	msgBulkAddSuccess     = "Universities added successfully"
)

type UserSchoolExperienceHandler struct {
	schoolExperienceService *services.UserSchoolExperienceService
}

func NewUserSchoolExperienceHandler(service *services.UserSchoolExperienceService) *UserSchoolExperienceHandler {
	return &UserSchoolExperienceHandler{
		schoolExperienceService: service,
	}
}

func (h *UserSchoolExperienceHandler) CreateSchoolExperience(c *fiber.Ctx) error {
	uid, err := ExtractAndValidateToken(c)
	if err != nil {
		return err
	}

	experience, err := h.schoolExperienceService.CreateSchoolExperience(c.Context(), uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create school experience",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "School experience created successfully",
		"data":    mappers.MapUserSchoolExperienceFromGoToFrontend(experience),
	})
}

func (h *UserSchoolExperienceHandler) GetSchoolExperience(c *fiber.Ctx) error {
	uid, err := ExtractAndValidateToken(c)
	if err != nil {
		return err
	}

	experience, err := h.schoolExperienceService.GetSchoolExperience(c.Context(), uid)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "School experience not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "School experience fetched successfully",
		"data":    mappers.MapUserSchoolExperienceFromGoToFrontend(experience),
	})
}

func (h *UserSchoolExperienceHandler) UpdateUniversity(c *fiber.Ctx) error {
	uid, err := ExtractAndValidateToken(c)
	if err != nil {
		return err
	}

	universityID := c.Params("universityID")
	if universityID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "University ID is required",
		})
	}

	var requestData map[string]interface{}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	experience, err := h.schoolExperienceService.UpdateUniversity(c.Context(), uid, universityID, requestData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update university",
		})
	}

	return c.Status(fiber.StatusOK).JSON(
		mappers.MapUserSchoolExperienceFromGoToFrontend(experience),
	)
}

func (h *UserSchoolExperienceHandler) AddUniversity(c *fiber.Ctx) error {
	uid, err := ExtractAndValidateToken(c)
	if err != nil {
		return err
	}

	var universityData map[string]interface{}
	if err := c.BodyParser(&universityData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	experience, err := h.schoolExperienceService.AddUniversity(c.Context(), uid, universityData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add university",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "University added successfully",
		"data":    mappers.MapUserSchoolExperienceFromGoToFrontend(experience),
	})
}

func (h *UserSchoolExperienceHandler) DeleteUniversity(c *fiber.Ctx) error {
	uid, err := ExtractAndValidateToken(c)
	if err != nil {
		return err
	}

	universityID := c.Params("id")
	if universityID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "University ID is required",
		})
	}

	err = h.schoolExperienceService.DeleteUniversity(c.Context(), uid, universityID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete university",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "University deleted successfully",
	})
}

func (h *UserSchoolExperienceHandler) AddListOfUniversities(c *fiber.Ctx) error {
	uid, err := ExtractAndValidateToken(c)
	if err != nil {
		return err
	}

	var requestBody struct {
		Universities []map[string]interface{} `json:"universities"`
	}

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errInvalidPayload,
		})
	}

	if len(requestBody.Universities) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No universities provided",
		})
	}

	experience, err := h.schoolExperienceService.AddListOfUniversities(c.Context(), uid, requestBody.Universities)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errAddUniversities,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": msgBulkAddSuccess,
		"data":    mappers.MapUserSchoolExperienceFromGoToFrontend(experience),
	})
}
