package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// User Routes
	users := api.Group("/users")
	users.Post("/register", handlers.Register)
	users.Post("/login", handlers.Login)
	users.Get("/:uid", handlers.GetUser)
	users.Put("/:uid", handlers.UpdateUser)

	// User Address Routes
	addresses := api.Group("/addresses")
	addresses.Post("/", handlers.CreateUserAddress)
	addresses.Put("/:id", handlers.UpdateUserAddress)
	addresses.Get("/", handlers.GetUserAddress)
	addresses.Delete("/:id", handlers.DeleteUserAddress)

	// User Work Experience Routes
	workExperiences := api.Group("/work-experiences")
	workExperiences.Post("/", handlers.CreateUserWorkExperience)
	workExperiences.Put("/:id", handlers.UpdateUserWorkExperience)
	workExperiences.Get("/", handlers.GetUserWorkExperience)
	workExperiences.Delete("/:id", handlers.DeleteUserWorkExperience)

	// User School Experience Routes
	schoolExperiences := api.Group("/school-experiences")
	schoolExperiences.Post("/", handlers.CreateSchoolExperience)
	schoolExperiences.Get("/:uid", handlers.GetSchoolExperience)
	schoolExperiences.Put("/:uid/universities/:universityID", handlers.UpdateUniversity)
	schoolExperiences.Delete("/:uid/universities/:universityID", handlers.DeleteUniversity)
	schoolExperiences.Post("/:uid/universities", handlers.AddUniversity)
	schoolExperiences.Post("/:uid/universities/bulk", handlers.AddListOfUniversities)
}
