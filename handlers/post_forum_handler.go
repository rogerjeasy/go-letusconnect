package handlers

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"github.com/rogerjeasy/go-letusconnect/services"
)

type ForumHandler struct {
	forumService *services.ForumService
	userService  *services.UserService
}

func NewForumHandler(forumService *services.ForumService, userService *services.UserService) *ForumHandler {
	return &ForumHandler{
		forumService: forumService,
		userService:  userService,
	}
}

func (h *ForumHandler) CreateForum(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	userID, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	var forum models.Forum
	if err := c.BodyParser(&forum); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	ctx := context.Background()
	createdForum, err := h.forumService.CreateForum(ctx, forum, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Forum created successfully",
		"data":    mappers.MapForumGoToFrontend(*createdForum),
	})
}

func (h *ForumHandler) GetForum(c *fiber.Ctx) error {
	forumID := c.Params("id")
	if forumID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Forum ID is required",
		})
	}

	ctx := context.Background()
	forum, err := h.forumService.GetForum(ctx, forumID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data":    mappers.MapForumGoToFrontend(*forum),
		"message": "Forum retrieved successfully",
	})
}

func (h *ForumHandler) CreatePost(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	userID, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	forumID := c.Params("id")
	if forumID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Forum ID is required",
		})
	}

	var post models.Post
	if err := c.BodyParser(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	ctx := context.Background()
	createdPost, err := h.forumService.CreatePost(ctx, forumID, post, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Post created successfully",
		"data":    createdPost,
	})
}

func (h *ForumHandler) CreateComment(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	userID, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	forumID := c.Params("forumId")
	postID := c.Params("postId")
	if forumID == "" || postID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Forum ID and Post ID are required",
		})
	}

	var comment models.Comment
	if err := c.BodyParser(&comment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	ctx := context.Background()
	createdComment, err := h.forumService.CreateComment(ctx, forumID, postID, comment, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Comment created successfully",
		"data":    createdComment,
	})
}

func (h *ForumHandler) AddReaction(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	userID, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	forumID := c.Params("id")
	if forumID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Forum ID is required",
		})
	}

	var reaction models.Reaction
	if err := c.BodyParser(&reaction); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	ctx := context.Background()
	err = h.forumService.AddReaction(ctx, forumID, reaction, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Reaction added successfully",
	})
}

func (h *ForumHandler) ListForumsByGroup(c *fiber.Ctx) error {
	groupID := c.Params("groupId")
	if groupID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Group ID is required",
		})
	}

	ctx := context.Background()
	forums, err := h.forumService.ListForumsByGroup(ctx, groupID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	frontendForums := make([]map[string]interface{}, 0, len(forums))
	for _, forum := range forums {
		frontendForums = append(frontendForums, mappers.MapForumGoToFrontend(forum))
	}

	return c.JSON(fiber.Map{
		"data": frontendForums,
	})
}

func (h *ForumHandler) SearchPosts(c *fiber.Ctx) error {
	forumID := c.Params("id")
	if forumID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Forum ID is required",
		})
	}

	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Search query is required",
		})
	}

	ctx := context.Background()
	posts, err := h.forumService.SearchPosts(ctx, forumID, query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": posts,
	})
}

func (h *ForumHandler) UpdateForum(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	userID, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	forumID := c.Params("id")
	if forumID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Forum ID is required",
		})
	}

	// Check if user is a moderator before allowing update
	ctx := context.Background()
	forum, err := h.forumService.GetForum(ctx, forumID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	isModerator := false
	for _, mod := range forum.Moderators {
		if mod.UID == userID {
			isModerator = true
			break
		}
	}

	if !isModerator {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Only forum moderators can update the forum",
		})
	}

	var updates models.Forum
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Preserve existing moderators
	updates.Moderators = forum.Moderators

	if err := h.forumService.UpdateForum(ctx, forumID, updates); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Forum updated successfully",
	})
}

func (h *ForumHandler) DeleteForum(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	userID, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	forumID := c.Params("id")
	if forumID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Forum ID is required",
		})
	}

	ctx := context.Background()
	if err := h.forumService.DeleteForum(ctx, forumID, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Forum deleted successfully",
	})
}

func (h *ForumHandler) AddModerator(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	currentUserID, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	forumID := c.Params("id")
	if forumID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Forum ID is required",
		})
	}

	// Check if current user is a moderator
	ctx := context.Background()
	forum, err := h.forumService.GetForum(ctx, forumID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	isModerator := false
	for _, mod := range forum.Moderators {
		if mod.UID == currentUserID {
			isModerator = true
			break
		}
	}

	if !isModerator {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Only existing moderators can add new moderators",
		})
	}

	// Get new moderator user ID from request body
	var reqBody struct {
		UserID string `json:"userId"`
	}
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if err := h.forumService.AddModerator(ctx, forumID, reqBody.UserID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Moderator added successfully",
	})
}

func (h *ForumHandler) RemoveModerator(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	currentUserID, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	forumID := c.Params("id")
	userIDToRemove := c.Params("userId")
	if forumID == "" || userIDToRemove == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Forum ID and user ID are required",
		})
	}

	// Check if current user is a moderator
	ctx := context.Background()
	forum, err := h.forumService.GetForum(ctx, forumID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	isModerator := false
	for _, mod := range forum.Moderators {
		if mod.UID == currentUserID {
			isModerator = true
			break
		}
	}

	if !isModerator {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Only existing moderators can remove moderators",
		})
	}

	// Prevent self-removal if last moderator
	if currentUserID == userIDToRemove && len(forum.Moderators) <= 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot remove the last moderator",
		})
	}

	if err := h.forumService.RemoveModerator(ctx, forumID, userIDToRemove); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Moderator removed successfully",
	})
}
