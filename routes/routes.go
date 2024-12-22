package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	// "github.com/rogerjeasy/go-letusconnect/middleware"
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
	projects.Get("/owner", handlers.GetOwnerProjects)
	projects.Get("/participation", handlers.GetParticipationProjects)
	projects.Get("/public", handlers.GetAllPublicProjects)
	// projects.Use(middleware.AuthMiddleware)
	projects.Post("/", handlers.CreateProject)
	// projects.Get("/", handlers.GetAllProjects)
	projects.Get("/:id", handlers.GetProject)
	projects.Put("/:id", handlers.UpdateProject)
	projects.Delete("/:id", handlers.DeleteProject)
	// 2. Collaboration Endpoints
	projects.Post("/:id/join", handlers.JoinProjectCollab)
	projects.Put("/:id/join-requests/:uid", handlers.AcceptRejectJoinRequestCollab)
	projects.Post("/:id/invite", handlers.InviteUserCollab)
	projects.Delete("/:id/participants/:uid", handlers.RemoveParticipantCollab)

	// 3. Task Endpoints
	projects.Post("/:id/tasks", handlers.AddTask)
	projects.Put("/:id/tasks/:taskID", handlers.UpdateTask)
	projects.Delete("/:id/tasks/:taskID", handlers.DeleteTask)

	// newsletters
	newsletters := api.Group("/newsletters")
	newsletters.Post("/subscribe", handlers.SubscribeNewsletter)
	newsletters.Post("/unsubscribe", handlers.UnsubscribeNewsletter)
	newsletters.Get("/subscribers", handlers.GetAllSubscribers)
	newsletters.Get("/subscribers/count", handlers.GetTotalSubscribers)

	// Pusher Routes
	SetupPusherRoutes(app)

	// Message Routes
	messages := api.Group("/messages")
	messages.Post("/send", handlers.SendMessage)
	messages.Get("/", handlers.GetMessages)
	messages.Post("/typing", handlers.SendTyping)
	messages.Post("/direct", handlers.SendDirectMessage)
	messages.Post("/group", handlers.SendGroupMessage)
	messages.Get("/direct", handlers.GetDirectMessages)
	messages.Get("/unread", handlers.GetUnreadMessagesCount)
	messages.Post("/mark-as-read", handlers.MarkMessagesAsRead)

	// Group Chat Routes
	groupChats := api.Group("/group-chats")
	groupChats.Post("/", handlers.CreateGroupChatF)
	groupChats.Get("/:id", handlers.GetGroupChat)
	groupChats.Get("/projects/:projectId/group-chats", handlers.GetGroupChatsByProject)
	groupChats.Get("/my/group-chats", handlers.GetMyGroupChats)
	groupChats.Post("/messages", handlers.SendMessageHandler)
	groupChats.Post("/mark-messages-read", handlers.MarkMessagesAsReadHandler)
	groupChats.Get("/unread-messages/count", handlers.CountUnreadMessagesHandler)
	groupChats.Post("/remove-participant", handlers.RemoveParticipantFromGroupChatHandler)

	// Media File Routes
	mediaFiles := api.Group("/media-files")
	mediaFiles.Post("/upload-images", handlers.UploadImageHandler)
	mediaFiles.Post("/upload-videos", handlers.UploadVideoHandler)
	mediaFiles.Post("/upload-pdf", handlers.UploadPDFHandler)

}
