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

const (
	errInvalidPayload     = "Invalid request payload"
	errExperienceNotFound = "School experience not found"
	errFetchExperience    = "Failed to fetch school experience"
	errParseExperience    = "Failed to parse school experience data"
	errAddUniversity      = "Failed to add university"
	errAddUniversities    = "Failed to add universities"
	msgAddSuccess         = "University added successfully"
	msgBulkAddSuccess     = "Universities added successfully"
)

// schoolExperienceDoc represents the document and its data
type schoolExperienceDoc struct {
	ref        *firestore.DocumentRef
	experience *models.UserSchoolExperience
}

func CreateSchoolExperience(c *fiber.Ctx) error {

	// Extract and validate token
	uid, err := extractAndValidateToken(c)
	if err != nil {
		return err
	}

	// context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = checkExistingExperience(ctx, uid)
	if err != nil {
		return handleFirestoreError(c, err)
	}

	experience := createNewExperience(uid)

	// Save to Firestore with retry mechanism
	err = saveExperience(ctx, experience)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errCreateExperience,
		})
	}

	// Return success response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": msgCreateSuccess,
		"data": mappers.MapUserSchoolExperienceBackendToFrontend(
			mappers.MapUserSchoolExperienceFrontendToBackend(experience),
		),
	})
}

// GetUserSchoolExperience retrieves a user's school experience
func GetSchoolExperience(c *fiber.Ctx) error {
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
	// Extract and validate token
	uid, err := extractAndValidateToken(c)
	if err != nil {
		return err
	}

	universityID := c.Params("universityID")

	if universityID == "" {
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

func AddUniversity(c *fiber.Ctx) error {
	// Extract and validate token
	uid, err := extractAndValidateToken(c)
	if err != nil {
		return err
	}

	// Parse and validate request body
	var universityData map[string]interface{}
	if err := validateRequestBody(c, &universityData); err != nil {
		return err
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get school experience document
	doc, err := getSchoolExperience(ctx, uid)
	if err != nil {
		return handleFirestoreError(c, err)
	}

	// Add university with transaction
	if err := addUniversityTransaction(ctx, doc, universityData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errAddUniversity,
		})
	}

	// Return updated data
	return GetSchoolExperience(c)
}

func AddListOfUniversities(c *fiber.Ctx) error {
	// Extract and validate token
	uid, err := extractAndValidateToken(c)
	if err != nil {
		return err
	}

	// Parse and validate request body
	var universitiesData []map[string]interface{}
	if err := validateRequestBody(c, &universitiesData); err != nil {
		return err
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Get school experience document
	doc, err := getSchoolExperience(ctx, uid)
	if err != nil {
		return handleFirestoreError(c, err)
	}

	// Add universities with transaction
	if err := addUniversitiesTransaction(ctx, doc, universitiesData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errAddUniversities,
		})
	}

	// Return success response with updated data
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": msgBulkAddSuccess,
		"data": mappers.MapUserSchoolExperienceBackendToFrontend(
			mappers.MapUserSchoolExperienceFrontendToBackend(doc.experience),
		),
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
