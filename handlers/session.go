// handlers/session.go
package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/services"
	"google.golang.org/api/iterator"
)

func GetSession(c *fiber.Ctx) error {

	// Get token from Authorization header or cookie
	var token string
	authHeader := c.Get("Authorization")
	if authHeader != "" {
		token = strings.TrimPrefix(authHeader, "Bearer ")
	} else {
		token = c.Cookies("jwt")
		if token != "" {
			log.Printf("Token found in cookie: %s", token[:10])
		} else {
			log.Printf("No token found in either Authorization header or cookie")
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "No authentication token provided",
			})
		}
	}

	// Validate token
	uid, err := validateToken(token)
	if err != nil {
		log.Printf("Token validation failed: %v", err)
		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    "",
			Expires:  time.Now().Add(-time.Hour),
			HTTPOnly: true,
			Secure:   true,
			SameSite: "Lax",
			Path:     "/",
		})
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": fmt.Sprintf("Token validation failed: %v", err),
		})
	}

	ctx := context.Background()

	query := services.FirestoreClient.Collection("users").Where("uid", "==", uid).Documents(ctx)
	defer query.Stop()

	doc, err := query.Next()
	if err == iterator.Done {
		log.Printf("No user found for UID: %s", uid)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": fmt.Sprintf("No user found for UID: %s", uid),
		})
	}
	if err != nil {
		log.Printf("Firestore query error: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Database error: %v", err),
		})
	}

	var dbUser map[string]interface{}
	if err := doc.DataTo(&dbUser); err != nil {
		log.Printf("Error parsing user data: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to parse user data: %v", err),
		})
	}

	// Convert to user model
	backendUser := mappers.MapBackendToUser(dbUser)
	if backendUser.Email == "" {
		log.Printf("Mapping produced invalid user: %+v", backendUser)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to map user data correctly",
		})
	}

	// Generate new token
	newToken, err := GenerateJWT(&backendUser)
	if err != nil {
		log.Printf("Error generating new token: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to generate token: %v", err),
		})
	}

	// Set cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    newToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Path:     "/",
	})

	frontendUser := mappers.MapUserToFrontend(&backendUser)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"user":  frontendUser,
		"token": newToken,
	})
}
