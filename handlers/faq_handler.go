package handlers

import (
	"context"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/services"
	"google.golang.org/api/iterator"
)

type FAQHandler struct {
	FAQService  *services.FAQService
	UserService *services.UserService
}

// NewFAQHandler creates a new instance of FAQHandler with all required services
func NewFAQHandler(faqService *services.FAQService, userService *services.UserService) *FAQHandler {
	if faqService == nil {
		panic("faqService cannot be nil")
	}
	if userService == nil {
		panic("userService cannot be nil")
	}

	return &FAQHandler{
		FAQService:  faqService,
		UserService: userService,
	}
}

// CreateFAQ handles the creation of a new FAQ entry
func (f *FAQHandler) CreateFAQ(c *fiber.Ctx) error {
	// Extract and validate token
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

	// Get user details
	userDetails, err := f.UserService.GetUserByUID(uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch user details",
			"details": err.Error(),
		})
	}

	username, ok := userDetails["username"].(string)
	if !ok || username == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid or missing username",
		})
	}

	// Validate admin privileges
	roles, err := f.UserService.GetUserRole(uid)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   "Failed to fetch user roles",
			"details": err.Error(),
		})
	}

	isAdmin := false
	for _, role := range roles {
		if role == "admin" {
			isAdmin = true
			break
		}
	}
	if !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Admin privileges required",
		})
	}

	// Parse request body
	var requestData map[string]interface{}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
	}

	// Map frontend data to FAQ struct
	faq := mappers.MapFAQFrontendToGo(requestData)

	ctx := context.Background()

	// Create FAQ using service
	newFAQ, err := f.FAQService.CreateFAQ(ctx, faq, username, uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create FAQ",
			"details": err.Error(),
		})
	}

	// Convert FAQ to frontend format for response
	response := mappers.MapFAQGoToFrontend(*newFAQ)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "FAQ created successfully",
		"data":    response,
	})
}

// UpdateFAQ updates an existing FAQ (admin only)
func (f *FAQHandler) UpdateFAQ(c *fiber.Ctx) error {
	// Extract the Authorization token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	// Validate token and get UID
	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Fetch the user's username
	username, err := f.UserService.GetUsernameByUID(uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user details",
		})
	}

	// Fetch the user's roles
	roles, err := f.UserService.GetUserRole(uid)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Failed to fetch user roles",
		})
	}

	// Check if the user has the "admin" role
	isAdmin := false
	for _, role := range roles {
		if role == "admin" {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Admin privileges required",
		})
	}

	// Get the FAQ ID from the route parameters
	faqID := c.Params("id")
	if faqID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "FAQ ID is required",
		})
	}

	// Parse the request payload
	var requestData map[string]interface{}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Prepare updates for Firestore
	ctx := context.Background()
	updates := []firestore.Update{
		{Path: "question", Value: requestData["question"]},
		{Path: "response", Value: requestData["response"]},
		{Path: "category", Value: requestData["category"]},
		{Path: "status", Value: requestData["status"]},
		{Path: "updated_by", Value: username},
		{Path: "updated_at", Value: time.Now()},
	}

	// Update the FAQ in Firestore
	_, err = services.Firestore.Collection("faqs").Doc(faqID).Update(ctx, updates)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update FAQ",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "FAQ updated successfully",
	})
}

// DeleteFAQ deletes an existing FAQ (admin only)
func (f *FAQHandler) DeleteFAQ(c *fiber.Ctx) error {
	// Extract the Authorization token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	// Validate token and get UID
	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Fetch the user's roles
	roles, err := f.UserService.GetUserRole(uid)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Failed to fetch user roles",
		})
	}

	// Check if the user has the "admin" role
	isAdmin := false
	for _, role := range roles {
		if role == "admin" {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Admin privileges required",
		})
	}

	// Get the FAQ ID from the route parameters
	faqID := c.Params("id")
	if faqID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "FAQ ID is required",
		})
	}

	// Delete the FAQ from Firestore
	ctx := context.Background()
	_, err = services.Firestore.Collection("faqs").Doc(faqID).Delete(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete FAQ",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "FAQ deleted successfully",
	})
}

func (f *FAQHandler) GetAllFAQs(c *fiber.Ctx) error {
	ctx := context.Background()
	iter := services.Firestore.Collection("faqs").Documents(ctx)
	defer iter.Stop()

	faqs := make([]map[string]interface{}, 0)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch FAQs",
			})
		}

		data := doc.Data()
		data["id"] = doc.Ref.ID

		// Convert Firestore data to Go struct
		faq := mappers.MapFAQFirestoreToGo(data)

		// Convert Go struct to frontend format
		frontendFAQ := mappers.MapFAQGoToFrontend(faq)
		faqs = append(faqs, frontendFAQ)
	}

	return c.Status(fiber.StatusOK).JSON(faqs)
}

func (f *FAQHandler) GetFAQByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "FAQ ID is required",
		})
	}

	ctx := context.Background()

	faq, err := f.FAQService.GetFAQByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "FAQ not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "FAQ not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch FAQ",
			"details": err.Error(),
		})
	}

	// Convert FAQ to frontend format
	response := mappers.MapFAQGoToFrontend(*faq)

	return c.Status(fiber.StatusOK).JSON(response)
}
