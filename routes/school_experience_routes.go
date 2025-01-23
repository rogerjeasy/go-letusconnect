package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

func setupUserSchoolExperienceRoutes(api fiber.Router, sc *services.ServiceContainer) error {
	if api == nil {
		return fmt.Errorf("api router cannot be nil")
	}
	if sc == nil {
		return fmt.Errorf("service container cannot be nil")
	}
	if sc.UserSchoolExperienceService == nil {
		return fmt.Errorf("school experience service cannot be nil")
	}

	schoolExperienceHandler := handlers.NewUserSchoolExperienceHandler(sc.UserSchoolExperienceService)
	if schoolExperienceHandler == nil {
		return fmt.Errorf("failed to create school experience handler")
	}

	schoolExperience := api.Group("/school-experiences")

	schoolExperience.Post("/", schoolExperienceHandler.CreateSchoolExperience)
	schoolExperience.Get("/", schoolExperienceHandler.GetSchoolExperience)

	universities := schoolExperience.Group("/universities")
	universities.Post("/", schoolExperienceHandler.AddUniversity)
	universities.Post("/bulk", schoolExperienceHandler.AddListOfUniversities)
	universities.Put("/:id", schoolExperienceHandler.UpdateUniversity)
	universities.Delete("/:id", schoolExperienceHandler.DeleteUniversity)

	return nil
}
