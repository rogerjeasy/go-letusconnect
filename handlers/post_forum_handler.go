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
