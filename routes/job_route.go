package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

// setupJobRoutes initializes job-related routes
func setupJobRoutes(api fiber.Router, sc *services.ServiceContainer) error {
	if api == nil {
		return fmt.Errorf("api router cannot be nil")
	}
	if sc == nil {
		return fmt.Errorf("service container cannot be nil")
	}
	if sc.JobService == nil {
		return fmt.Errorf("job service cannot be nil")
	}

	handler := handlers.NewJobHandler(sc.JobService)
	if handler == nil {
		return fmt.Errorf("failed to create job handler")
	}

	// Job Routes
	jobs := api.Group("/jobs")

	jobs.Post("/", handler.CreateJobHandler)
	jobs.Get("/:id", handler.GetJobHandler)
	jobs.Get("/", handler.GetJobsByUserHandler)
	jobs.Put("/:id", handler.UpdateJobHandler)
	jobs.Delete("/:id", handler.DeleteJobHandler)
	jobs.Get("/status/:status", handler.GetJobsByStatusHandler)

	// Interview Rounds
	jobs.Post("/:id/interviews", handler.AddInterviewRoundHandler)
	jobs.Delete("/:id/interviews/:roundNumber", handler.RemoveInterviewRoundHandler)

	return nil
}
