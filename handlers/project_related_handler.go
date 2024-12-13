package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"github.com/rogerjeasy/go-letusconnect/services"
)

// JoinProject handles applying to join a project
func JoinProjectCollab(c *fiber.Ctx) error {
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

	// Parse the join request message
	var requestData struct {
		Message string `json:"message"`
	}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	username, err := services.GetUsernameByUID(uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user details",
		})
	}

	// Create the join request
	joinRequest := models.JoinRequest{
		UserID:      uid,
		UserName:    username,
		Message:     "Request to join the project",
		RequestedAt: time.Now(),
		Status:      "pending",
	}

	ctx := context.Background()

	// Fetch the project from Firestore
	doc, err := services.FirestoreClient.Collection("projects").Doc(projectID).Get(ctx)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Project not found",
		})
	}

	projectData := doc.Data()

	// Check if the user has already applied
	joinRequests := mappers.GetJoinRequestsArray(projectData, "join_requests")
	for _, jr := range joinRequests {
		if jr.UserID == uid {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "You have already applied to join this project",
			})
		}
	}

	// convert joinRequest to map[string]interface{} for Firestore
	joinRequestMap := mappers.MapJoinRequestGoToFirestore(joinRequest)

	_, err = services.FirestoreClient.Collection("projects").Doc(projectID).Update(ctx, []firestore.Update{
		{Path: "join_requests", Value: firestore.ArrayUnion(joinRequestMap)},
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to apply to join project",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Join request submitted successfully",
	})
}

// HandleJoinRequest handles accepting or rejecting join requests
func AcceptRejectJoinRequestCollab(c *fiber.Ctx) error {
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
	userID := c.Params("uid")
	if projectID == "" || userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Project ID and User ID are required",
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

	projectData := doc.Data()
	project := mappers.MapProjectFirestoreToGo(projectData)

	// Check if the user is the project owner or an owner participant
	isOwner := project.OwnerID == uid
	for _, participant := range project.Participants {
		if participant.UserID == uid && participant.Role == "owner" {
			isOwner = true
			break
		}
	}

	if !isOwner {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You do not have permission to handle join requests",
		})
	}

	// Parse the request payload for action ("accept" or "reject")
	var requestData struct {
		Action string `json:"action"`
	}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Handle the join request based on action
	var updatedJoinRequests []map[string]interface{}
	for _, jr := range projectData["join_requests"].([]interface{}) {
		jrMap := jr.(map[string]interface{})
		if jrMap["user_id"] == userID {
			if requestData.Action == "accept" {
				// Add the user as a participant
				newParticipant := map[string]interface{}{
					"user_id":   userID,
					"role":      "member",
					"joined_at": time.Now().Format(time.RFC3339),
				}
				services.FirestoreClient.Collection("projects").Doc(projectID).Update(ctx, []firestore.Update{
					{Path: "participants", Value: firestore.ArrayUnion(newParticipant)},
				})
			}
			// Skip adding the join request to updatedJoinRequests to remove it
		} else {
			updatedJoinRequests = append(updatedJoinRequests, jrMap)
		}
	}

	// Update the join requests in Firestore
	_, err = services.FirestoreClient.Collection("projects").Doc(projectID).Update(ctx, []firestore.Update{
		{Path: "join_requests", Value: updatedJoinRequests},
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update join request",
		})
	}

	action := requestData.Action

	var actionPastTense string
	switch action {
	case "accept":
		actionPastTense = "accepted"
	case "reject":
		actionPastTense = "rejected"
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid action",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("Join request %s successfully", actionPastTense),
	})

}

// InviteUserCollab handles inviting a user to a project
func InviteUserCollab(c *fiber.Ctx) error {
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

	projectData := doc.Data()
	project := mappers.MapProjectFirestoreToGo(projectData)

	// Check if the user is the project owner or an owner participant
	isOwner := project.OwnerID == uid
	for _, participant := range project.Participants {
		if participant.UserID == uid && participant.Role == "owner" {
			isOwner = true
			break
		}
	}

	if !isOwner {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You do not have permission to invite users to this project",
		})
	}

	// Parse the request payload for the user to invite
	var requestData struct {
		UserID string `json:"user_id"`
	}
	if err := c.BodyParser(&requestData); err != nil || requestData.UserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload or missing user_id",
		})
	}

	// Check if the user is already in the invited_users list
	username, err := services.GetUsernameByUID(requestData.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch the user's details",
		})
	}

	// Convert invited_users from Firestore format to Go struct format
	for _, invitedUser := range projectData["invited_users"].([]interface{}) {
		invitedUserMap, ok := invitedUser.(map[string]interface{})
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to parse invited users",
			})
		}

		invitedUserStruct := mappers.MapInvitedUserFirestoreToGo(invitedUserMap)

		if invitedUserStruct.UserID == requestData.UserID {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("%s has already been invited to this project", username),
			})
		}
	}

	// Add the user to invited_users
	invite := map[string]interface{}{
		"user_id":   requestData.UserID,
		"role":      "invited",
		"joined_at": time.Now().Format(time.RFC3339),
	}

	_, err = services.FirestoreClient.Collection("projects").Doc(projectID).Update(ctx, []firestore.Update{
		{Path: "invited_users", Value: firestore.ArrayUnion(invite)},
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to invite user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User invited successfully",
	})
}
