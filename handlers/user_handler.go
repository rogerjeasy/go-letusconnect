package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"github.com/rogerjeasy/go-letusconnect/services"
	"google.golang.org/api/iterator"
)

func UpdateUser(c *fiber.Ctx) error {
	uid := c.Params("uid")
	if strings.TrimSpace(uid) == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "UID is required",
		})
	}

	// Extract and validate the token
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	requesterUID, err := validateToken(tokenString)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Ensure the requester is authorized to update the user
	if requesterUID != uid {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": "You are not authorized to update this user",
		})
	}

	ctx := context.Background()
	userUpdates := make(map[string]interface{})

	// Parse the request body into a generic map
	if err := c.BodyParser(&userUpdates); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload. Unable to parse JSON",
		})
	}

	// Fetch the current user document from Firestore
	userQuery := services.FirestoreClient.Collection("users").Where("uid", "==", uid).Documents(ctx)
	defer userQuery.Stop()

	doc, err := userQuery.Next()
	if err == iterator.Done {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	if err != nil {
		log.Printf("Error querying user: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	// Map Firestore document data to a user model
	var currentUser models.User
	if err := doc.DataTo(&currentUser); err != nil {
		log.Printf("Error mapping user data: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse user data",
		})
	}

	// Validate and sanitize user updates
	if username, ok := userUpdates["username"].(string); ok && strings.TrimSpace(username) != "" {
		if len(username) < 3 {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Username must be at least 3 characters long",
			})
		}

		// Check for uniqueness of the new username (if changed)
		if username != currentUser.Username {
			usernameQuery := services.FirestoreClient.Collection("users").Where("username", "==", username).Documents(ctx)
			defer usernameQuery.Stop()
			if _, err := usernameQuery.Next(); err != iterator.Done {
				return c.Status(http.StatusConflict).JSON(fiber.Map{
					"error": "Username already in use",
				})
			}
		}
	}

	// Map updates to Firestore-compatible format
	backendUpdates := mappers.MapFrontendToBackend(&currentUser)
	for key, value := range userUpdates {
		if key == "account_creation_date" {
			// Skip updating the account creation date
			continue
		}
		backendUpdates[key] = value
	}

	// Apply updates to Firestore
	_, err = services.FirestoreClient.Collection("users").Doc(doc.Ref.ID).Set(ctx, backendUpdates, firestore.MergeAll)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	// Re-fetch the updated user document
	updatedDoc, err := services.FirestoreClient.Collection("users").Doc(doc.Ref.ID).Get(ctx)
	if err != nil {
		log.Printf("Error fetching updated user: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch updated user",
		})
	}

	var updatedUser models.User
	if err := updatedDoc.DataTo(&updatedUser); err != nil {
		log.Printf("Error mapping updated user data: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse updated user data",
		})
	}

	// Map updated user to frontend format
	frontendUser := mappers.MapBackendToFrontend(updatedUser)

	// Return the updated user
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "User updated successfully",
		"user":    frontendUser,
	})
}

func GetUser(c *fiber.Ctx) error {
	uid := c.Params("uid")
	if strings.TrimSpace(uid) == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "UID is required",
		})
	}

	ctx := context.Background()

	// Query Firestore for a user with the specified UID
	userQuery := services.FirestoreClient.Collection("users").Where("UID", "==", uid).Documents(ctx)
	defer userQuery.Stop()

	doc, err := userQuery.Next()
	if err == iterator.Done {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	if err != nil {
		log.Printf("Error querying user: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	// Map Firestore document data to a user model
	var user models.User
	if err := doc.DataTo(&user); err != nil {
		log.Printf("Error mapping user data: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse user data",
		})
	}

	// Convert the backend user to frontend format
	frontendUser := mappers.MapBackendToFrontend(user)

	// Return the user data in frontend format
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "User retrieved successfully",
		"user":    frontendUser,
	})
}
