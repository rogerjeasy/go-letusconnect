package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

func setupLinkedInJobRoutes(api fiber.Router, sc *services.ServiceContainer) error {
	if api == nil {
		return fmt.Errorf("api router cannot be nil")
	}
	if sc == nil {
		return fmt.Errorf("service container cannot be nil")
	}
	if sc.LinkedInJobsService == nil {
		return fmt.Errorf("linkedin jobs service cannot be nil")
	}

	handler := handlers.NewLinkedInJobsHandler(sc.LinkedInJobsService)
	if handler == nil {
		return fmt.Errorf("failed to create linkedin jobs handler")
	}

	// LinkedIn Job Routes
	linkedinJobs := api.Group("/linkedin/jobs")
	linkedinJobs.Get("/applied", handler.GetAppliedJobs)
	linkedinJobs.Post("/applications", handler.StoreJobApplication)
	linkedinJobs.Get("/applications", handler.GetUserApplications)

	return nil
}
