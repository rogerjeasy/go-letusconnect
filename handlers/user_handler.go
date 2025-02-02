package handlers

import (
	"context"
	"math"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/services"
)

type UserHandler struct {
	userService      *services.UserService
	containerService *services.ServiceContainer
}

func NewUserHandler(userService *services.UserService, containerService *services.ServiceContainer) *UserHandler {
	return &UserHandler{
		userService:      userService,
		containerService: containerService,
	}
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	// Extract Authorization token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	// Validate token
	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
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
	fetchUser, err := h.userService.GetUserByUID(uid)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
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
		if username != fetchUser["username"] {
			res, err := h.userService.CheckUsernameUniqueness(username)
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to check username uniqueness",
				})
			}
			if !res {
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{
					"error": "Username already exists",
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
	err = h.containerService.AuthService.UpdateUser(ctx, uid, backendUpdates)
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

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	uid := c.Params("uid")
	if strings.TrimSpace(uid) == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "UID is required",
		})
	}

	ctx := context.Background()
	dbUser, err := h.userService.GetUserByUID(uid)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	userGoStruct := mappers.MapBackendToUser(dbUser)
	if userGoStruct.IsPrivate {
		token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		if token == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error":     "This account is private",
				"isPrivate": true,
			})
		}

		// Validate token and check if requester is connected to user
		requesterUID, err := validateToken(token)
		if err != nil || !isConnected(ctx, requesterUID, uid) {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error":     "This account is private",
				"isPrivate": true,
			})
		}
	}

	frontendUser := mappers.MapUserBackendToFrontend(dbUser)
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "User retrieved successfully",
		"user":    frontendUser,
	})
}

func isConnected(ctx context.Context, requesterUID, targetUID string) bool {
	connectionsRef := services.Firestore.Collection("connections")
	query := connectionsRef.Where("status", "==", "accepted").
		Where("users", "array-contains", requesterUID)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return false
	}

	for _, doc := range docs {
		users := doc.Data()["users"].([]interface{})
		for _, user := range users {
			if user.(string) == targetUID {
				return true
			}
		}
	}
	return false
}

// GetAllUsers retrieves all users from the Firestore "users" collection
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	ctx := context.Background()

	// Query all documents in the "users" collection
	users, err := h.userService.GetAllUsers(ctx)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve users",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Users retrieved successfully",
		"users":   users,
	})
}

func (h *UserHandler) GetProfileCompletion(c *fiber.Ctx) error {
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

	userDetails, err := h.userService.GetUserByUID(uid)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Calculate filled fields for each component
	filledFields := 0
	totalFields := 19 // Total number of fields excluding system fields
	missingFields := make([]string, 0)
	user := mappers.MapBackendToUser(userDetails)

	// User basic info with missing fields tracking
	if user.Username != "" {
		filledFields++
	} else {
		missingFields = append(missingFields, "username")
	}

	if user.FirstName != "" {
		filledFields++
	} else {
		missingFields = append(missingFields, "First name")
	}

	if user.LastName != "" {
		filledFields++
	} else {
		missingFields = append(missingFields, "Last name")
	}

	if user.Email != "" {
		filledFields++
	} else {
		missingFields = append(missingFields, "Email")
	}

	if user.PhoneNumber != "" {
		filledFields++
	} else {
		missingFields = append(missingFields, "Phone number")
	}

	if user.ProfilePicture != "" {
		filledFields++
	} else {
		missingFields = append(missingFields, "Profile picture")
	}

	if user.Bio != "" {
		filledFields++
	} else {
		missingFields = append(missingFields, "Bio")
	}

	if len(user.Role) > 0 {
		filledFields++
	} else {
		missingFields = append(missingFields, "Role")
	}

	if user.GraduationYear > 0 {
		filledFields++
	} else {
		missingFields = append(missingFields, "Graduation year")
	}

	if user.CurrentJobTitle != "" {
		filledFields++
	} else {
		missingFields = append(missingFields, "Current job title")
	}

	if len(user.AreasOfExpertise) > 0 {
		filledFields++
	} else {
		missingFields = append(missingFields, "Areas of expertise")
	}

	if len(user.Interests) > 0 {
		filledFields++
	} else {
		missingFields = append(missingFields, "Interests")
	}

	if user.Program != "" {
		filledFields++
	} else {
		missingFields = append(missingFields, "Program")
	}

	// if user.DateOfBirth != "" {
	// 	filledFields++
	// } else {
	// 	missingFields = append(missingFields, "dateOfBirth")
	// }

	// if user.PhoneCode != "" {
	// 	filledFields++
	// } else {
	// 	missingFields = append(missingFields, "phoneCode")
	// }

	if len(user.Languages) > 0 {
		filledFields++
	} else {
		missingFields = append(missingFields, "Languages")
	}

	if len(user.Skills) > 0 {
		filledFields++
	} else {
		missingFields = append(missingFields, "Skills")
	}

	if len(user.Certifications) > 0 {
		filledFields++
	} else {
		missingFields = append(missingFields, "Certifications")
	}

	if len(user.Projects) > 0 {
		filledFields++
	} else {
		missingFields = append(missingFields, "Projects")
	}

	// Calculate completion percentage
	var completionPercentage float64
	if totalFields > 0 {
		completionPercentage = float64(filledFields) / float64(totalFields) * 100
	}

	return c.JSON(fiber.Map{
		"completionPercentage": math.Round(completionPercentage),
		"filledFields":         filledFields,
		"totalFields":          totalFields,
		"missingFields":        missingFields,
		"profileStatus":        getProfileStatus(completionPercentage),
	})
}

// Helper function to determine profile status
func getProfileStatus(percentage float64) string {
	switch {
	case percentage >= 100:
		return "Complete"
	case percentage >= 80:
		return "Almost Complete"
	case percentage >= 50:
		return "Partially Complete"
	default:
		return "Incomplete"
	}
}
