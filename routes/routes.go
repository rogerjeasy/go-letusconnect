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
	users.Get("/", handlers.GetAllUsers)

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

	// Contact Us Routes
	contacts := api.Group("/contact-us")
	contacts.Post("/", handlers.CreateContact)
	contacts.Get("/", handlers.GetAllContacts)
	contacts.Get("/:id", handlers.GetContactByID)
	contacts.Put("/:id", handlers.UpdateContactStatus)

	// FAQ Routes
	faqs := api.Group("/faqs")
	faqs.Get("/", handlers.GetAllFAQs)
	faqs.Post("/", handlers.CreateFAQ)
	faqs.Put("/:id", handlers.UpdateFAQ)
	faqs.Delete("/:id", handlers.DeleteFAQ)

	// Project Management Routes
	projects := api.Group("/projects")
	projects.Post("/", handlers.CreateProject)
	// projects.Get("/", handlers.GetAllProjects)
	projects.Get("/:id", handlers.GetProject)
	projects.Put("/:id", handlers.UpdateProject)
	projects.Delete("/:id", handlers.DeleteProject)
	// 2. Collaboration Endpoints
	projects.Post("/:id/join", handlers.JoinProjectCollab)
	projects.Put("/:id/join-requests/:uid", handlers.AcceptRejectJoinRequestCollab)
	projects.Post("/:id/invite", handlers.InviteUserCollab)

}
