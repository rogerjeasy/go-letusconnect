package handlers

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"github.com/rogerjeasy/go-letusconnect/services"
)

type GroupHandler struct {
	groupService *services.GroupService
	userService  *services.UserService
}

func NewGroupHandler(groupService *services.GroupService, userService *services.UserService) *GroupHandler {
	return &GroupHandler{
		groupService: groupService,
		userService:  userService,
	}
}

func (h *GroupHandler) CreateGroup(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	userId, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	var group models.Group
	if err := c.BodyParser(&group); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	ctx := context.Background()
	createdGroup, err := h.groupService.CreateGroup(ctx, group, userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Group created successfully",
		"data":    mappers.MapGroupGoToFrontend(*createdGroup),
	})
}

func (h *GroupHandler) GetGroup(c *fiber.Ctx) error {
	groupID := c.Params("id")
	if groupID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Group ID is required",
		})
	}

	ctx := context.Background()
	group, err := h.groupService.GetGroup(ctx, groupID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data":    mappers.MapGroupGoToFrontend(*group),
		"message": "Group retrieved successfully",
	})
}

func (h *GroupHandler) UpdateGroup(c *fiber.Ctx) error {
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

	groupID := c.Params("id")

	ctx := context.Background()
	existingGroup, err := h.groupService.GetGroup(ctx, groupID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	isAdmin := false
	for _, admin := range existingGroup.Admins {
		if admin.UID == userID {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Only group admins can update the group",
		})
	}

	var updates models.Group
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	updates.Admins = existingGroup.Admins

	ctx = context.Background()
	if err := h.groupService.UpdateGroup(ctx, groupID, updates); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Group updated successfully",
	})
}

func (h *GroupHandler) DeleteGroup(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	groupID := c.Params("id")
	ctx := context.Background()

	if err := h.groupService.DeleteGroup(ctx, groupID, uid); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Group deleted successfully",
	})
}

func (h *GroupHandler) ListGroups(c *fiber.Ctx) error {
	filters := make(map[string]interface{})

	if category := c.Query("category"); category != "" {
		filters["category.name"] = category
	}

	if privacy := c.Query("privacy"); privacy != "" {
		filters["privacy"] = privacy
	}

	if featured := c.Query("featured"); featured != "" {
		filters["featured"] = featured == "true"
	}

	ctx := context.Background()
	groups, err := h.groupService.ListGroups(ctx, filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	frontendGroups := make([]map[string]interface{}, 0, len(groups))
	for _, group := range groups {
		frontendGroups = append(frontendGroups, mappers.MapGroupGoToFrontend(group))
	}

	return c.JSON(fiber.Map{
		"data": frontendGroups,
	})
}

func (h *GroupHandler) AddMember(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	groupID := c.Params("id")
	var member models.Member
	if err := c.BodyParser(&member); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	ctx := context.Background()
	if err := h.groupService.AddMember(ctx, groupID, member); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Member added successfully",
	})
}

func (h *GroupHandler) RemoveMember(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	groupID := c.Params("id")
	userID := c.Params("userId")

	ctx := context.Background()
	if err := h.groupService.RemoveMember(ctx, groupID, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Member removed successfully",
	})
}

func (h *GroupHandler) UploadGroupImage(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	groupID := c.Params("id")
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Image file is required",
		})
	}

	ctx := context.Background()
	if err := h.groupService.UploadGroupImage(ctx, groupID, file); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Group image uploaded successfully",
	})
}

func (h *GroupHandler) AddEvent(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	groupID := c.Params("id")
	var event models.Event
	if err := c.BodyParser(&event); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	ctx := context.Background()
	if err := h.groupService.AddEvent(ctx, groupID, event); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Event added successfully",
	})
}

func (h *GroupHandler) RemoveEvent(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	groupID := c.Params("id")
	eventID := c.Params("eventId")

	ctx := context.Background()
	if err := h.groupService.RemoveEvent(ctx, groupID, eventID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Event removed successfully",
	})
}

func (h *GroupHandler) AddResource(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	groupID := c.Params("id")
	var resource models.Resource
	if err := c.BodyParser(&resource); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	ctx := context.Background()
	if err := h.groupService.AddResource(ctx, groupID, resource); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Resource added successfully",
	})
}

func (h *GroupHandler) RemoveResource(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	groupID := c.Params("id")
	resourceID := c.Params("resourceId")

	ctx := context.Background()
	if err := h.groupService.RemoveResource(ctx, groupID, resourceID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Resource removed successfully",
	})
}

func (h *GroupHandler) UpdateGroupSettings(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	groupID := c.Params("id")
	var settings struct {
		Privacy  string `json:"privacy"`
		Featured bool   `json:"featured"`
	}
	if err := c.BodyParser(&settings); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	ctx := context.Background()
	if err := h.groupService.UpdateGroupSettings(ctx, groupID, settings.Privacy, settings.Featured); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Group settings updated successfully",
	})
}

func (h *GroupHandler) SearchGroups(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Search query is required",
		})
	}

	ctx := context.Background()
	groups, err := h.groupService.SearchGroups(ctx, query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": groups,
	})
}
