package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

	user, err := services.GetUserByUID(uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user details",
		})
	}

	// Create the join request
	joinRequest := models.JoinRequest{
		UserID:         uid,
		Username:       user["username"].(string),
		ProfilePicture: user["profile_picture"].(string),
		Email:          user["email"].(string),
		Message:        "Request to join the project",
		RequestedAt:    time.Now(),
		Status:         "pending",
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

	if projectData["status"] == "completed" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "This project has been completed. You cannot join a completed project.",
		})
	}

	// Check if the user is the owner of the project
	if projectData["owner_id"] == uid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Owners cannot join their own project",
		})
	}

	// Check if the user has already applied
	joinRequests := mappers.GetJoinRequestsArray(projectData, "join_requests")
	for _, jr := range joinRequests {
		if jr.UserID == uid {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "You have already applied to join this project",
			})
		}
	}

	// Check if the user was previously rejected
	rejectedParticipants := projectData["rejected_participants"].([]interface{})
	for _, rejectedUID := range rejectedParticipants {
		if rejectedUID == uid {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Your request to join this project was previously rejected. Please contact the project owner for further assistance.",
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

	// Send confirmation email
	projectName := projectData["title"].(string)
	if err := SendJoinRequestSubmittedEmail(user["email"].(string), user["username"].(string), projectName); err != nil {
		log.Printf("Error sending join request submitted email: %v", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Join request submitted successfully",
	})
}

// RemoveParticipantCollab handles removing a participant from a project
func RemoveParticipantCollab(c *fiber.Ctx) error {
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

	// Extract project ID and participant user ID from URL parameters
	projectID := c.Params("id")
	participantID := c.Params("participantId")
	if projectID == "" || participantID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Project ID and Participant ID are required",
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

	// Check if the user is the project owner or an owner-level participant
	isOwner := project.OwnerID == uid
	for _, participant := range project.Participants {
		if participant.UserID == uid && participant.Role == "owner" {
			isOwner = true
			break
		}
	}

	if !isOwner {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You do not have permission to remove participants",
		})
	}

	// Check if the participant exists in the project
	participantExists := false
	var updatedParticipants []map[string]interface{}

	for _, participant := range projectData["participants"].([]interface{}) {
		participantMap := participant.(map[string]interface{})
		if participantMap["user_id"] == participantID {
			participantExists = true
			continue // Skip the participant to remove them
		}
		updatedParticipants = append(updatedParticipants, participantMap)
	}

	if !participantExists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Participant not found in the project",
		})
	}

	// Update the participants list in Firestore
	_, err = services.FirestoreClient.Collection("projects").Doc(projectID).Update(ctx, []firestore.Update{
		{Path: "participants", Value: updatedParticipants},
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove participant",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Participant removed successfully",
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

	// Parse the request payload for action ("accept" or "reject") and user details
	var requestData struct {
		Action         string `json:"action"`
		Role           string `json:"role"`
		Username       string `json:"username"`
		Email          string `json:"email"`
		ProfilePicture string `json:"profilePicture"`
		UserID         string `json:"user_id"`
	}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Ensure join_requests is a valid slice
	joinRequests, ok := projectData["join_requests"].([]interface{})
	if !ok {
		log.Printf("join_requests field is of unexpected type: %+v", projectData["join_requests"])
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Join requests data is invalid",
		})
	}

	// Handle the join request based on action
	var updatedJoinRequests []map[string]interface{}
	participantExists := false

	for _, jr := range joinRequests {
		jrMap, ok := jr.(map[string]interface{})
		if !ok {
			continue
		}

		if jrMap["user_id"] == userID {
			participantExists = true
			if requestData.Action == "accept" {
				// Add the user as a participant with the provided attributes
				newParticipant := models.Participant{
					UserID:         userID,
					Role:           requestData.Role,
					Username:       requestData.Username,
					Email:          requestData.Email,
					ProfilePicture: requestData.ProfilePicture,
					JoinedAt:       time.Now(),
				}
				_, err := services.FirestoreClient.Collection("projects").Doc(projectID).Update(ctx, []firestore.Update{
					{Path: "participants", Value: firestore.ArrayUnion(mappers.MapParticipantGoToFirestore(newParticipant))},
				})
				if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": "Failed to add participant",
					})
				}

				// Add the participant to the group chat using projectID
				if err := services.AddParticipantToGroupChat(ctx, "", projectID, newParticipant); err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": fmt.Sprintf("Failed to add participant to group chat: %v", err),
					})
				}
				// Send notification email for accepted request
				if err := SendJoinRequestAcceptedEmail(requestData.Email, requestData.Username, project.Title); err != nil {
					log.Printf("Failed to send acceptance email: %v", err)
				}
			} else if requestData.Action == "reject" {
				// Add the user ID to the RejectedParticipants list
				_, err := services.FirestoreClient.Collection("projects").Doc(projectID).Update(ctx, []firestore.Update{
					{Path: "rejected_participants", Value: firestore.ArrayUnion(userID)},
				})
				if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": "Failed to add user to rejected participants",
					})
				}
				// Send notification email for rejected request
				if err := SendJoinRequestRejectedEmail(requestData.Email, requestData.Username, project.Title); err != nil {
					log.Printf("Failed to send rejection email: %v", err)
				}
			}
			// Skip adding this join request to updatedJoinRequests to remove it
		} else {
			updatedJoinRequests = append(updatedJoinRequests, jrMap)
		}
	}

	if !participantExists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Join request not found",
		})
	}

	// Update the join requests in Firestore to remove the processed request
	_, err = services.FirestoreClient.Collection("projects").Doc(projectID).Update(ctx, []firestore.Update{
		{Path: "join_requests", Value: updatedJoinRequests},
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update join requests",
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
		EmailOrUsername string `json:"emailOrUsername"`
		Role            string `json:"role"`
	}
	if err := c.BodyParser(&requestData); err != nil || requestData.EmailOrUsername == "" || requestData.Role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload or missing fields",
		})
	}

	var user *models.User

	if strings.Contains(requestData.EmailOrUsername, "@") {
		// Handle invitation by email
		user, err = services.GetUserByEmail(requestData.EmailOrUsername)
	} else {
		// Handle invitation by username
		user, err = services.GetUserByUsername(requestData.EmailOrUsername)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Check if the user is already in the participants list
	for _, participant := range project.Participants {
		if participant.UserID == user.UID {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("%s is already a participant in this project", user.Username),
			})
		}
	}

	// Initialize invited_users if nil
	invitedUsersData, ok := projectData["invited_users"].([]interface{})
	if !ok || invitedUsersData == nil {
		fmt.Println("invited_users is nil or not a []interface{}. Initializing as empty slice.")
		invitedUsersData = []interface{}{}
	}

	// Check if the user is already in the invited_users list
	for _, existingUser := range invitedUsersData {
		existingUserMap, ok := existingUser.(map[string]interface{})
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid format for invited user",
			})
		}

		if existingUserMap["user_id"] == user.UID {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("%s has already been invited to this project", user.Username),
			})
		}
	}

	// Add the user to invited_users
	invite := map[string]interface{}{
		"user_id":         user.UID,
		"username":        user.Username,
		"email":           user.Email,
		"profile_picture": user.ProfilePicture,
		"role":            requestData.Role,
		"joined_at":       time.Now().Format(time.RFC3339),
	}

	_, err = services.FirestoreClient.Collection("projects").Doc(projectID).Update(ctx, []firestore.Update{
		{Path: "invited_users", Value: firestore.ArrayUnion(invite)},
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to invite user",
		})
	}

	// Send invitation email
	projectName := projectData["title"].(string)
	ownernerName, err := services.GetUsernameByUID(uid)
	if err != nil {
		log.Printf("Error fetching project owner's username: %v", err)
		ownernerName = "Project Owner"
	}
	if err := SendProjectInvitationEmail(user.Email, user.Username, projectName, ownernerName); err != nil {
		log.Printf("Error sending project invitation email: %v", err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "User invited successfully",
		"invitedUser": invite,
	})
}

// AddTask handles adding a task to a project
func AddTask(c *fiber.Ctx) error {
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

	// Check if the user is the project owner or a participant
	isAuthorized := project.OwnerID == uid
	for _, participant := range project.Participants {
		if participant.UserID == uid {
			isAuthorized = true
			break
		}
	}

	if !isAuthorized {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You do not have permission to add tasks to this project",
		})
	}

	// Parse the request payload for the task details
	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Validate that the title is not empty
	if task.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Task title is required",
		})
	}

	// Generate a unique task ID and set additional fields
	task.ID = uuid.New().String()
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	// Add the task to Firestore
	_, err = services.FirestoreClient.Collection("projects").Doc(projectID).Update(ctx, []firestore.Update{
		{Path: "tasks", Value: firestore.ArrayUnion(mappers.MapTaskGoToFirestore(task))},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add task to the project",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Task added successfully",
		"task":    mappers.MapTaskGoToFrontend(task),
	})
}

// UpdateTask handles updating task details
func UpdateTask(c *fiber.Ctx) error {
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
	taskID := c.Params("taskID")
	if projectID == "" || taskID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Project ID and Task ID are required",
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

	// Check if the user is authorized (project owner or participant)
	isAuthorized := project.OwnerID == uid
	for _, participant := range project.Participants {
		if participant.UserID == uid {
			isAuthorized = true
			break
		}
	}

	if !isAuthorized {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You do not have permission to update tasks in this project",
		})
	}

	// Parse the request payload for task updates
	var updatedTaskData map[string]interface{}
	if err := c.BodyParser(&updatedTaskData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Update the specific task
	taskUpdated := false
	for i, task := range project.Tasks {
		if task.ID == taskID {
			updatedTask := mappers.MapTaskFrontendToGo(updatedTaskData)
			updatedTask.ID = task.ID
			updatedTask.CreatedAt = task.CreatedAt
			updatedTask.UpdatedAt = time.Now()
			project.Tasks[i] = updatedTask
			taskUpdated = true
			break
		}
	}

	if !taskUpdated {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Task not found",
		})
	}

	// Save the updated project back to Firestore
	_, err = services.FirestoreClient.Collection("projects").Doc(projectID).Set(ctx, mappers.MapProjectGoToFirestore(project))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update task",
		})
	}

	// Map each task to frontend format
	var frontendTasks []map[string]interface{}
	for _, task := range project.Tasks {
		frontendTasks = append(frontendTasks, mappers.MapTaskGoToFrontend(task))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Task updated successfully",
		"tasks":   frontendTasks,
	})

}

// DeleteTask handles deleting a task from a project
func DeleteTask(c *fiber.Ctx) error {
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
	taskID := c.Params("taskID")
	if projectID == "" || taskID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Project ID and Task ID are required",
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

	// Check if the user is authorized
	isAuthorized := project.OwnerID == uid
	for _, participant := range project.Participants {
		if participant.UserID == uid {
			isAuthorized = true
			break
		}
	}

	if !isAuthorized {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You do not have permission to delete tasks from this project",
		})
	}

	// Filter out the task to be deleted
	var updatedTasks []map[string]interface{}
	for _, task := range project.Tasks {
		if task.ID != taskID {
			updatedTasks = append(updatedTasks, mappers.MapTaskGoToFirestore(task))
		}
	}

	// Debug: print the updated tasks before updating Firestore
	fmt.Printf("Updated Tasks: %+v\n", updatedTasks)

	// Update Firestore with the new tasks list
	_, err = services.FirestoreClient.Collection("projects").Doc(projectID).Set(ctx, map[string]interface{}{
		"tasks": updatedTasks,
	}, firestore.MergeAll)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete task",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Task deleted successfully",
	})
}
