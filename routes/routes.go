package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
	// "github.com/rogerjeasy/go-letusconnect/middleware"
)

func SetupRoutes(app *fiber.App, notificationService *services.NotificationService) {
	api := app.Group("/api")

	// User Routes
	users := api.Group("/users")

	users.Get("/session", handlers.GetSession)
	users.Patch("/logout", handlers.Logout)

	users.Post("/register", handlers.Register)
	users.Post("/login", handlers.Login)
	users.Get("/completion", handlers.GetProfileCompletion)
	users.Get("/", handlers.GetAllUsers)

	users.Get("/:uid", handlers.GetUser)
	users.Put("/:uid", handlers.UpdateUser)

	// Notification Routes
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	notifications := api.Group("/notifications")
	notifications.Get("/targeted", notificationHandler.ListTargetedNotifications)
	notifications.Get("/unread-count", notificationHandler.GetUnreadNotificationCount)
	notifications.Get("/stats", notificationHandler.GetNotificationStats)
	notifications.Patch("/:id", notificationHandler.MarkNotificationAsRead)
	notifications.Post("/", notificationHandler.CreateNotification)
	notifications.Get("/", notificationHandler.ListNotifications)
	notifications.Get("/:id", notificationHandler.GetNotification)
	notifications.Put("/:id", notificationHandler.UpdateNotification)
	notifications.Delete("/:id", notificationHandler.DeleteNotification)
	notifications.Put("/:id/read", notificationHandler.MarkNotificationAsRead)

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
	schoolExperiences.Get("/", handlers.GetSchoolExperience)
	schoolExperiences.Put("/universities/:universityID", handlers.UpdateUniversity)
	schoolExperiences.Delete("/:uid/universities/:universityID", handlers.DeleteUniversity)
	schoolExperiences.Post("/universities", handlers.AddUniversity)
	schoolExperiences.Post("/universities/bulk", handlers.AddListOfUniversities)

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
	// messages.Post("/group", handlers.SendGroupMessage)
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
	groupChats.Patch("/:groupChatId/mark-messages-read", handlers.MarkMessagesAsReadHandler)
	groupChats.Get("/unread-messages/count", handlers.CountUnreadMessagesHandler)
	groupChats.Post("/:groupChatId/remove-participants", handlers.RemoveParticipantsFromGroupChatHandler)
	groupChats.Post("/reply-message", handlers.ReplyToMessageHandler)
	groupChats.Post("/attach-files", handlers.AttachFilesToMessageHandler)
	groupChats.Post("/pin-message", handlers.PinMessageHandler)
	groupChats.Get("/pinned-messages", handlers.GetPinnedMessagesHandler)
	groupChats.Post("/unpin-message", handlers.UnpinMessageHandler)
	groupChats.Post("/react-to-message", handlers.ReactToMessageHandler)
	groupChats.Get("/message-read-receipts/:groupChatId/:messageId", handlers.GetMessageReadReceiptsHandler)
	groupChats.Post("/set-role", handlers.SetParticipantRoleHandler)
	groupChats.Post("/mute-participant", handlers.MuteParticipantHandler)
	groupChats.Get("/online-status/:participantId", handlers.UpdateLastSeenHandler)

	// Update Group Settings
	groupChats.Post("/update-settings", handlers.UpdateGroupSettingsHandler)

	// Archive Group Chat
	groupChats.Post("/archive", handlers.ArchiveGroupChatHandler)

	groupChats.Delete("/:groupChatId/participants/me", handlers.LeaveGroupHandler)
	groupChats.Put("/:groupChatId/participants", handlers.AddParticipantsToGroupChatHandler)
	groupChats.Put("/projects/:projectId/participants", handlers.AddParticipantsToGroupChatHandler) // Add participants to a group chat in a project

	// Polls
	// groupChats.Post("/create-poll", handlers.CreatePollHandler)
	// groupChats.Get("/polls/:groupChatId", handlers.GetPollsHandler)
	// groupChats.Post("/vote", handlers.VoteOnPollHandler)

	// // Star/Favorite Messages
	// groupChats.Post("/star-message", handlers.StarMessageHandler)
	// groupChats.Post("/unstar-message", handlers.UnstarMessageHandler)
	// groupChats.Get("/starred-messages/:groupChatId", handlers.GetStarredMessagesHandler)

	// Report Messages
	groupChats.Post("/report-message", handlers.ReportMessageHandler)
	// groupChats.Get("/reports/:groupChatId", handlers.GetReportsHandler)
	// groupChats.Post("/block-participant", handlers.BlockParticipantHandler)
	// groupChats.Post("/unblock-participant", handlers.UnblockParticipantHandler)

	// Media File Routes
	mediaFiles := api.Group("/media-files")
	mediaFiles.Post("/upload-images", handlers.UploadImageHandler)
	mediaFiles.Post("/upload-videos", handlers.UploadVideoHandler)
	mediaFiles.Post("/upload-pdf", handlers.UploadPDFHandler)

}
