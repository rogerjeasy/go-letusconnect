package handlers

import (
	"context"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"github.com/rogerjeasy/go-letusconnect/services"
	"google.golang.org/api/iterator"
)

// CreateFAQ handles the creation of a new FAQ entry
func CreateFAQ(c *fiber.Ctx) error {
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
	username, err := services.GetUsernameByUID(uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user details",
		})
	}

	// Fetch the user's role list
	roles, err := services.GetUserRole(uid)
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

	// Parse the request payload
	var requestData map[string]interface{}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	ctx := context.Background()

	// Map request data to FAQ model
	newFAQ := mappers.FrontendToFAQ(requestData, username, uid)

	// Save to Firestore
	docRef, _, err := services.FirestoreClient.Collection("faqs").Add(ctx, mappers.FAQToFirestore(newFAQ))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create FAQ",
		})
	}

	// Set the ID of the new FAQ
	newFAQ.ID = docRef.ID

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "FAQ created successfully",
		"data":    newFAQ,
	})
}

// GetAllFAQs retrieves all FAQs
func GetAllFAQs(c *fiber.Ctx) error {
	ctx := context.Background()
	iter := services.FirestoreClient.Collection("faqs").Documents(ctx)
	defer iter.Stop()

	var faqs []models.FAQ

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch FAQs"})
		}

		data := doc.Data()
		data["id"] = doc.Ref.ID
		faq := mappers.FirestoreToFAQ(data)
		faqs = append(faqs, *faq)
	}

	return c.Status(fiber.StatusOK).JSON(faqs)
}

// UpdateFAQ updates an existing FAQ (admin only)
func UpdateFAQ(c *fiber.Ctx) error {
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
	username, err := services.GetUsernameByUID(uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user details",
		})
	}

	// Fetch the user's roles
	roles, err := services.GetUserRole(uid)
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
		{Path: "updated_by", Value: username},
		{Path: "updated_at", Value: time.Now()},
	}

	// Update the FAQ in Firestore
	_, err = services.FirestoreClient.Collection("faqs").Doc(faqID).Update(ctx, updates)
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
func DeleteFAQ(c *fiber.Ctx) error {
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
	roles, err := services.GetUserRole(uid)
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
	_, err = services.FirestoreClient.Collection("faqs").Doc(faqID).Delete(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete FAQ",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "FAQ deleted successfully",
	})
}
