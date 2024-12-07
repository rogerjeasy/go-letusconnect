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

// creates a new UserSchoolExperience with an empty list of universities
func CreateSchoolExperience(c *fiber.Ctx) error {

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

	ctx := context.Background()

	// Use where() to check if a UserSchoolExperience already exists for this user
	existingQuery := services.FirestoreClient.Collection("user_school_experiences").Where("uid", "==", uid).Documents(ctx)
	defer existingQuery.Stop()

	_, err = existingQuery.Next()
	if err != iterator.Done {
		// If a document is found
		if err == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "School experience already exists for this user",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check existing school experience",
		})
	}

	// Initialize an empty UserSchoolExperience
	currentTime := time.Now()

	newExperience := models.UserSchoolExperience{
		UID:          uid,
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
		Universities: []models.University{},
	}

	backendData := mappers.MapUserSchoolExperienceFrontendToBackend(&newExperience)

	// Save to Firestore
	_, _, err = services.FirestoreClient.Collection("user_school_experiences").Add(ctx, backendData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create school experience",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "School experience created successfully",
		"data":    mappers.MapUserSchoolExperienceBackendToFrontend(backendData),
	})
}

// GetUserSchoolExperience retrieves a user's school experience
func GetSchoolExperience(c *fiber.Ctx) error {
	// Extract the UID from the route parameters
	uid := c.Params("uid")
	if uid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User UID is required",
		})
	}

	// Create Firestore context
	ctx := context.Background()

	// Query the Firestore collection to find documents where "uid" matches the provided UID
	query := services.FirestoreClient.Collection("user_school_experiences").Where("uid", "==", uid).Documents(ctx)

	// Get the first matching document
	doc, err := query.Next()
	if err == iterator.Done {
		// No document found
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "School experience not found",
		})
	}
	if err != nil {
		// Handle other errors
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch school experience",
		})
	}

	// Parse the Firestore document data into the desired frontend format
	var schoolExperience models.UserSchoolExperience
	if err := doc.DataTo(&schoolExperience); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse school experience data",
		})
	}

	frontendData := mappers.MapUserSchoolExperienceBackendToFrontend(doc.Data())

	// Return the school experience data to the frontend
	return c.Status(fiber.StatusOK).JSON(frontendData)
}

// UpdateUniversity updates a specific university in the user's education list
func UpdateUniversity(c *fiber.Ctx) error {
	uid := c.Params("uid")
	universityID := c.Params("universityID")

	if uid == "" || universityID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User UID and University ID are required",
		})
	}

	var requestData map[string]interface{}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	ctx := context.Background()

	// Locate the document by UID using where()
	query := services.FirestoreClient.Collection("user_school_experiences").Where("uid", "==", uid).Documents(ctx)
	defer query.Stop()

	doc, err := query.Next()
	if err == iterator.Done {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "School experience not found",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch school experience",
		})
	}

	var schoolExperience models.UserSchoolExperience
	if err := doc.DataTo(&schoolExperience); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse school experience data",
		})
	}

	// Find and update the university by ID
	updated := false
	for i, university := range schoolExperience.Universities {
		if university.ID == universityID {
			schoolExperience.Universities[i] = mappers.MapFrontendToUniversity(requestData)
			updated = true
			break
		}
	}

	if !updated {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "University not found",
		})
	}

	// Update Firestore
	schoolExperience.UpdatedAt = time.Now()
	backendData := mappers.MapUserSchoolExperienceFrontendToBackend(&schoolExperience)
	_, err = services.FirestoreClient.Collection("user_school_experiences").Doc(doc.Ref.ID).Set(ctx, backendData, firestore.MergeAll)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update university",
		})
	}

	return GetSchoolExperience(c)
}

// AddUniversity adds a new university to the user's school experience
func AddUniversity(c *fiber.Ctx) error {
	// Extract the UID from the route parameters
	uid := c.Params("uid")
	if uid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User UID is required",
		})
	}

	// Parse the request body to get the new university data
	var newUniversityData map[string]interface{}
	if err := c.BodyParser(&newUniversityData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Create Firestore context
	ctx := context.Background()

	// Query Firestore to locate the document by UID
	query := services.FirestoreClient.Collection("user_school_experiences").Where("uid", "==", uid).Documents(ctx)
	defer query.Stop()

	doc, err := query.Next()
	if err == iterator.Done {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "School experience not found",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch school experience",
		})
	}

	// Parse the school experience document
	var schoolExperience models.UserSchoolExperience
	if err := doc.DataTo(&schoolExperience); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse school experience data",
		})
	}

	// Map the incoming data to a University model
	newUniversity := mappers.MapFrontendToUniversity(newUniversityData)
	newUniversity.ID = services.GenerateID() // Generate a unique ID for the university

	// Add the new university to the user's education list
	schoolExperience.Universities = append(schoolExperience.Universities, newUniversity)
	schoolExperience.UpdatedAt = time.Now()

	// Update the document in Firestore
	backendData := mappers.MapUserSchoolExperienceFrontendToBackend(&schoolExperience)
	_, err = services.FirestoreClient.Collection("user_school_experiences").Doc(doc.Ref.ID).Set(ctx, backendData, firestore.MergeAll)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add university",
		})
	}

	// Reuse GetSchoolExperience to fetch and return the updated data
	return GetSchoolExperience(c)
}

func AddListOfUniversities(c *fiber.Ctx) error {
	uid := c.Params("uid")
	if uid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User UID is required",
		})
	}

	var universitiesData []map[string]interface{}
	if err := c.BodyParser(&universitiesData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload. Expected a list of universities.",
		})
	}

	ctx := context.Background()

	// Locate the document by UID using where()
	query := services.FirestoreClient.Collection("user_school_experiences").Where("uid", "==", uid).Documents(ctx)
	defer query.Stop()

	doc, err := query.Next()
	if err == iterator.Done {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "School experience not found",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch school experience",
		})
	}

	// Parse the school experience document
	var schoolExperience models.UserSchoolExperience
	if err := doc.DataTo(&schoolExperience); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse school experience data",
		})
	}

	// Map the incoming data to University models and append to the list
	for _, universityData := range universitiesData {
		newUniversity := mappers.MapFrontendToUniversity(universityData)
		newUniversity.ID = services.GenerateID() // Generate a unique ID for each university
		schoolExperience.Universities = append(schoolExperience.Universities, newUniversity)
	}
	schoolExperience.UpdatedAt = time.Now()

	// Update the document in Firestore
	backendData := mappers.MapUserSchoolExperienceFrontendToBackend(&schoolExperience)
	_, err = services.FirestoreClient.Collection("user_school_experiences").Doc(doc.Ref.ID).Set(ctx, backendData, firestore.MergeAll)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add universities",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Universities added successfully",
		"data":    mappers.MapUserSchoolExperienceBackendToFrontend(backendData),
	})
}

// DeleteUniversity removes a specific university from the user's education list
func DeleteUniversity(c *fiber.Ctx) error {

	// Extract the Authorization token from the headers
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	// Validate the token and retrieve the uid
	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}
	universityID := c.Params("universityID")

	if uid == "" || universityID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User UID and University ID are required",
		})
	}

	ctx := context.Background()

	// Locate the document by UID using where()
	query := services.FirestoreClient.Collection("user_school_experiences").Where("uid", "==", uid).Documents(ctx)
	defer query.Stop()

	doc, err := query.Next()
	if err == iterator.Done {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "School experience not found",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch school experience",
		})
	}

	var schoolExperience models.UserSchoolExperience
	if err := doc.DataTo(&schoolExperience); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse school experience data",
		})
	}

	// Remove the university by ID
	newUniversities := []models.University{}
	for _, university := range schoolExperience.Universities {
		if university.ID != universityID {
			newUniversities = append(newUniversities, university)
		}
	}

	if len(newUniversities) == len(schoolExperience.Universities) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "University not found",
		})
	}

	schoolExperience.Universities = newUniversities
	schoolExperience.UpdatedAt = time.Now()

	backendData := mappers.MapUserSchoolExperienceFrontendToBackend(&schoolExperience)
	_, err = services.FirestoreClient.Collection("user_school_experiences").Doc(doc.Ref.ID).Set(ctx, backendData, firestore.MergeAll)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete university",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "University deleted successfully",
	})
}
