package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

// CreateProject handles the creation of a new project
func CreateProject(c *fiber.Ctx) error {
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

	// Fetch the user's details (username)
	username, err := services.GetUsernameByUID(uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user details",
		})
	}

	// Parse the request payload
	var requestData map[string]interface{}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Validate mandatory fields
	mandatoryFields := []string{"title", "description", "collaborationType", "academicFields"}
	for _, field := range mandatoryFields {
		if _, ok := requestData[field]; !ok || requestData[field] == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("%s is required", field),
			})
		}
	}

	// Set additional fields for the project
	requestData["ownerId"] = uid
	requestData["ownerUsername"] = username
	requestData["participants"] = []map[string]interface{}{
		{
			"userId":   uid,
			"role":     "owner",
			"joinedAt": time.Now().Format(time.RFC3339),
		},
	}

	// Convert participants to the correct []interface{} format
	if participants, ok := requestData["participants"].([]map[string]interface{}); ok {
		var participantInterfaces []interface{}
		for _, p := range participants {
			participantInterfaces = append(participantInterfaces, p)
		}
		requestData["participants"] = participantInterfaces
	}

	requestData["status"] = "open"
	requestData["createdAt"] = time.Now().Format(time.RFC3339)
	requestData["updatedAt"] = time.Now().Format(time.RFC3339)

	// Map request data to Project model
	newProject := mappers.MapProjectFrontendToGo(requestData)

	ctx := context.Background()

	// Save to Firestore
	docRef, _, err := services.FirestoreClient.Collection("projects").Add(ctx, mappers.MapProjectGoToFirestore(newProject))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create project",
		})
	}

	c.Locals("id", docRef.ID)
	c.Locals("message", "Project created successfully")

	return GetProject(c)
}

// UpdateProject handles updating project details
func UpdateProject(c *fiber.Ctx) error {
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

	projectID := c.Params("id")
	if projectID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Project ID is required",
		})
	}

	ctx := context.Background()

	// Fetch the project from Firestore
	doc, err := services.FirestoreClient.Collection("projects").Doc(projectID).Get(ctx)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Project not found",
		})
	}

	// Map Firestore data to a Project struct
	projectData := doc.Data()
	project := mappers.MapProjectFirestoreToGo(projectData)

	// Check if the user is the project owner
	isOwner := project.OwnerID == uid

	// Check if the user is an invited user with the role "owner"
	for _, participant := range project.Participants {
		if participant.UserID == uid && participant.Role == "owner" {
			isOwner = true
			break
		}
	}

	if !isOwner {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You do not have permission to update this project",
		})
	}

	// Parse the request payload
	var requestData map[string]interface{}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Validate mandatory fields
	mandatoryFields := []string{"title", "description", "collaborationType", "academicFields"}
	for _, field := range mandatoryFields {
		if _, ok := requestData[field]; !ok || requestData[field] == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("%s is required", field),
			})
		}
	}

	// Update the `updatedAt` field
	requestData["updatedAt"] = time.Now().Format(time.RFC3339)

	// Map request data to Project model
	updatedProject := mappers.MapProjectFrontendToGo(requestData)

	// If OwnerID is empty, preserve the original OwnerID
	if updatedProject.OwnerID == "" {
		updatedProject.OwnerID = project.OwnerID
	}

	// If OwnerUsername is empty, preserve the original OwnerUsername
	if updatedProject.OwnerUsername == "" {
		updatedProject.OwnerUsername = project.OwnerUsername
	}

	// Update project in Firestore
	_, err = services.FirestoreClient.Collection("projects").Doc(projectID).Set(ctx, mappers.MapProjectGoToFirestore(updatedProject), firestore.MergeAll)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update project",
		})
	}

	c.Locals("id", projectID)
	c.Locals("message", "Project updated successfully")
	c.Locals("token", token)

	return GetProject(c)
}

// GetProject handles fetching a project by its ID
func GetProject(c *fiber.Ctx) error {
	// Extract the Authorization token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	// Validate token and get UID
	_, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Get the project ID from the route parameter
	projectID := c.Params("id")
	if projectID == "" {
		if id, ok := c.Locals("id").(string); ok {
			projectID = id
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Project ID is required",
			})
		}
	}

	ctx := context.Background()

	// Fetch the project document from Firestore
	doc, err := services.FirestoreClient.Collection("projects").Doc(projectID).Get(ctx)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Project not found",
		})
	}

	projectData := doc.Data()

	// map project to client format
	projectFrontend := mappers.MapProjectFirestoreToFrontend(projectData)
	projectFrontend["id"] = projectID

	// Get custom message if set, otherwise use default message
	message := "Project fetched successfully"
	if customMessage, ok := c.Locals("message").(string); ok {
		message = customMessage
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": message,
		"data":    projectFrontend,
	})
}

// DeleteProject handles deleting a project by its ID
func DeleteProject(c *fiber.Ctx) error {
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

	// Get the project ID from the route parameter
	projectID := c.Params("id")
	if projectID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Project ID is required",
		})
	}

	ctx := context.Background()

	// Fetch the project document from Firestore
	doc, err := services.FirestoreClient.Collection("projects").Doc(projectID).Get(ctx)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Project not found",
		})
	}

	// Map Firestore data to a Project struct
	projectData := doc.Data()
	project := mappers.MapProjectFirestoreToGo(projectData)

	// Check if the user is the project owner
	isOwner := project.OwnerID == uid

	// Check if the user is a participant with the role "owner"
	for _, participant := range project.Participants {
		if participant.UserID == uid && participant.Role == "owner" {
			isOwner = true
			break
		}
	}

	if !isOwner {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You do not have permission to delete this project",
		})
	}

	// Delete the project from Firestore
	_, err = services.FirestoreClient.Collection("projects").Doc(projectID).Delete(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete project",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Project deleted successfully",
	})
}
