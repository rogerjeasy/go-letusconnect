package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"github.com/rogerjeasy/go-letusconnect/services"
)

type JobHandler struct {
	JobService *services.JobService
}

func NewJobHandler(jobService *services.JobService) *JobHandler {
	return &JobHandler{JobService: jobService}
}

// CreateJobHandler handles job creation requests
func (h *JobHandler) CreateJobHandler(c *fiber.Ctx) error {
	uid, err := ExtractAndValidateToken(c)
	if err != nil {
		return err
	}

	var jobData map[string]interface{}
	if err := c.BodyParser(&jobData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	timeStr := jobData["applicationDate"].(string)
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid date format"})
	}

	job := mappers.MapJobFrontendToGo(jobData)
	job.UserID = uid
	job.ApplicationDate = t

	createdJob, err := h.JobService.CreateJob(c.Context(), &job)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create job: " + err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Job created successfully",
		"data":    mappers.MapJobGoToFrontend(*createdJob),
	})
}

// GetJobHandler fetches a single job by ID
func (h *JobHandler) GetJobHandler(c *fiber.Ctx) error {
	uid, err := ExtractAndValidateToken(c)
	if err != nil {
		return err
	}

	jobID := c.Params("id")
	if jobID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Job ID is required"})
	}

	job, err := h.JobService.GetJob(c.Context(), jobID, uid)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Job retrieved successfully",
		"data":    mappers.MapJobGoToFrontend(*job),
	})
}

// GetJobsByUserHandler fetches all jobs for a user
func (h *JobHandler) GetJobsByUserHandler(c *fiber.Ctx) error {
	uid, err := ExtractAndValidateToken(c)
	if err != nil {
		return err
	}

	jobs, err := h.JobService.GetJobsByUser(c.Context(), uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch jobs. Reason: " + err.Error()})
	}

	var jobList []map[string]interface{}
	for _, job := range jobs {
		jobList = append(jobList, mappers.MapJobGoToFrontend(job))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Jobs fetched successfully",
		"data":    jobList,
	})
}

// UpdateJobHandler updates an existing job
func (h *JobHandler) UpdateJobHandler(c *fiber.Ctx) error {
	uid, err := ExtractAndValidateToken(c)
	if err != nil {
		return err
	}

	jobID := c.Params("id")
	if jobID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Job ID is required"})
	}

	// Parse request payload
	var jobData map[string]interface{}
	if err := c.BodyParser(&jobData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	err = h.JobService.UpdateJob(c.Context(), jobID, uid, jobData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Job updated successfully"})
}

// DeleteJobHandler deletes a job
func (h *JobHandler) DeleteJobHandler(c *fiber.Ctx) error {
	uid, err := ExtractAndValidateToken(c)
	if err != nil {
		return err
	}

	jobID := c.Params("id")
	if jobID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Job ID is required"})
	}

	err = h.JobService.DeleteJob(c.Context(), jobID, uid)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Job deleted successfully"})
}

// AddInterviewRoundHandler adds an interview round to a job
func (h *JobHandler) AddInterviewRoundHandler(c *fiber.Ctx) error {
	uid, err := ExtractAndValidateToken(c)
	if err != nil {
		return err
	}

	jobID := c.Params("id")
	if jobID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Job ID is required"})
	}

	// Parse interview round data
	var interview models.InterviewRound
	if err := c.BodyParser(&interview); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Add interview round
	err = h.JobService.AddInterviewRound(c.Context(), jobID, uid, interview)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Interview round added successfully"})
}

// RemoveInterviewRoundHandler removes an interview round from a job
func (h *JobHandler) RemoveInterviewRoundHandler(c *fiber.Ctx) error {
	uid, err := ExtractAndValidateToken(c)
	if err != nil {
		return err
	}

	jobID := c.Params("id")
	roundNumber := c.Params("roundNumber")
	if jobID == "" || roundNumber == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Job ID and Round Number are required"})
	}

	err = h.JobService.RemoveInterviewRound(c.Context(), jobID, uid, stringToInt(roundNumber))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Interview round removed successfully"})
}

// GetJobsByStatusHandler fetches jobs by their status
func (h *JobHandler) GetJobsByStatusHandler(c *fiber.Ctx) error {
	uid, err := ExtractAndValidateToken(c)
	if err != nil {
		return err
	}

	status := c.Params("status")
	if status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Job status is required"})
	}

	jobs, err := h.JobService.GetJobsByStatus(c.Context(), uid, models.JobStatus(status))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch jobs"})
	}

	var jobList []map[string]interface{}
	for _, job := range jobs {
		jobList = append(jobList, mappers.MapJobGoToFrontend(job))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Jobs fetched successfully",
		"data":    jobList,
	})
}

// Helper function to convert string to int
func stringToInt(s string) int {
	var num int
	fmt.Sscanf(s, "%d", &num)
	return num
}
