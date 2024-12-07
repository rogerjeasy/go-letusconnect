package handlers

import (
	"context"
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
	// Extract Authorization token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	// Validate token
	requesterUID, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Get User UID from parameters
	uid := c.Params("uid")
	if uid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User UID is required",
		})
	}

	// Ensure the requester is authorized to update the user
	if requesterUID != uid {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You are not authorized to update this user",
		})
	}

	// Parse request body into a map
	var requestData map[string]interface{}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Query the user document from Firestore using where()
	ctx := context.Background()
	userQuery := services.FirestoreClient.Collection("users").Where("uid", "==", uid).Documents(ctx)
	defer userQuery.Stop()

	doc, err := userQuery.Next()
	if err == iterator.Done {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	// Validate and sanitize user updates
	if username, ok := requestData["username"].(string); ok && strings.TrimSpace(username) != "" {
		if len(username) < 3 {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Username must be at least 3 characters long",
			})
		}

		// Check for uniqueness of the new username (if changed)
		if username != doc.Data()["username"] {
			usernameQuery := services.FirestoreClient.Collection("users").Where("username", "==", username).Documents(ctx)
			defer usernameQuery.Stop()
			if _, err := usernameQuery.Next(); err != iterator.Done {
				return c.Status(http.StatusConflict).JSON(fiber.Map{
					"error": "Username already in use",
				})
			}
		}
	}

	// Map the frontend data to a User struct
	updatedUser := mappers.MapFrontendToUser(requestData)

	updatedUser.UID = uid

	// Convert the updated User struct to Firestore-compatible format
	backendUpdates := mappers.MapUserFrontendToBackend(&updatedUser)

	// Update Firestore document
	docRef := services.FirestoreClient.Collection("users").Doc(doc.Ref.ID)
	_, err = docRef.Set(ctx, backendUpdates, firestore.MergeAll)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	// Map the backend data to frontend format for response
	frontendResponse := mappers.MapUserBackendToFrontend(backendUpdates)

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User updated successfully",
		"user":    frontendResponse,
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

	doc, err := services.FirestoreClient.Collection("users").Doc(uid).Get(ctx)
	if err != nil {
		if err == iterator.Done {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	var user models.User
	if err := doc.DataTo(&user); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse user data",
		})
	}

	frontendUser := mappers.MapUserBackendToFrontend(mappers.MapUserFrontendToBackend(&user))
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "User retrieved successfully",
		"user":    frontendUser,
	})
}
