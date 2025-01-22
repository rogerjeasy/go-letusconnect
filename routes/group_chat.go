package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

func setupGroupChatRoutes(api fiber.Router, sc *services.ServiceContainer) error {
	if api == nil {
		return fmt.Errorf("api router cannot be nil")
	}
	if sc == nil {
		return fmt.Errorf("service container cannot be nil")
	}
	if sc.GroupChatService == nil {
		return fmt.Errorf("group chat service cannot be nil")
	}
	if sc.UserService == nil {
		return fmt.Errorf("user service cannot be nil")
	}

	handler := handlers.NewGroupChatHandler(sc.GroupChatService, sc.UserService)
	if handler == nil {
		return fmt.Errorf("failed to create group chat handler")
	}

	groupChats := api.Group("/group-chats")

	// Group Chat Routes
	groupChats.Post("/", handler.CreateGroupChatF)
	groupChats.Get("/:id", handler.GetGroupChat)
	groupChats.Get("/projects/:projectId/group-chats", handler.GetGroupChatsByProject)
	groupChats.Get("/my/group-chats", handler.GetMyGroupChats)
	groupChats.Post("/messages", handler.SendMessageHandler)
	groupChats.Patch("/:groupChatId/mark-messages-read", handler.MarkMessagesAsReadHandler)
	groupChats.Get("/unread-messages/count", handler.CountUnreadMessagesHandler)
	groupChats.Get("/unread/total", handler.CountUnreadGroupMessagesFromAllChatHandler)
	groupChats.Post("/:groupChatId/remove-participants", handler.RemoveParticipantsFromGroupChatHandler)
	groupChats.Post("/reply-message", handler.ReplyToMessageHandler)
	groupChats.Post("/attach-files", handler.AttachFilesToMessageHandler)
	groupChats.Post("/pin-message", handler.PinMessageHandler)
	groupChats.Get("/pinned-messages", handler.GetPinnedMessagesHandler)
	groupChats.Post("/unpin-message", handler.UnpinMessageHandler)
	groupChats.Post("/react-to-message", handler.ReactToMessageHandler)
	groupChats.Get("/message-read-receipts/:groupChatId/:messageId", handler.GetMessageReadReceiptsHandler)
	groupChats.Post("/set-role", handler.SetParticipantRoleHandler)
	groupChats.Post("/mute-participant", handler.MuteParticipantHandler)
	groupChats.Get("/online-status/:participantId", handler.UpdateLastSeenHandler)

	// delete group chat
	groupChats.Delete("/:id", handler.DeleteGroupChat)
	groupChats.Delete("/delete-multiple", handler.DeleteMultipleGroupChats)

	// Update Group Settings
	groupChats.Put("/update-settings/:groupChatId", handler.UpdateGroupSettingsHandler)

	// Archive Group Chat
	groupChats.Post("/archive", handler.ArchiveGroupChatHandler)

	groupChats.Delete("/:groupChatId/participants/me", handler.LeaveGroupHandler)
	groupChats.Put("/:groupChatId/participants", handler.AddParticipantsToGroupChatHandler)
	groupChats.Put("/projects/:projectId/participants", handler.AddParticipantsToGroupChatHandler)

	// Polls
	// groupChats.Post("/create-poll", handlers.CreatePollHandler)
	// groupChats.Get("/polls/:groupChatId", handlers.GetPollsHandler)
	// groupChats.Post("/vote", handlers.VoteOnPollHandler)

	// // Star/Favorite Messages
	// groupChats.Post("/star-message", handlers.StarMessageHandler)
	// groupChats.Post("/unstar-message", handlers.UnstarMessageHandler)
	// groupChats.Get("/starred-messages/:groupChatId", handlers.GetStarredMessagesHandler)

	// Report Messages
	groupChats.Post("/report-message", handler.ReportMessageHandler)
	// groupChats.Get("/reports/:groupChatId", handlers.GetReportsHandler)
	// groupChats.Post("/block-participant", handlers.BlockParticipantHandler)
	// groupChats.Post("/unblock-participant", handlers.UnblockParticipantHandler)

	return nil
}
