package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"github.com/rogerjeasy/go-letusconnect/services"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error responses
const (
	errTokenRequired    = "Authorization token is required"
	errInvalidToken     = "Invalid token"
	errExperienceExists = "School experience already exists for this user"
	errCheckExisting    = "Failed to check existing school experience"
	errCreateExperience = "Failed to create school experience"
	msgCreateSuccess    = "School experience created successfully"
)

func ValidateToken(tokenString string) (string, error) {
	if strings.TrimSpace(tokenString) == "" {
		log.Printf("Empty token received")
		return "", errors.New("token cannot be empty")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Invalid signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		log.Printf("Token parsing error: %v", err)
		return "", fmt.Errorf("failed to parse token: %v", err)
	}

	if !token.Valid {
		log.Printf("Token is invalid")
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Printf("Failed to extract token claims")
		return "", errors.New("invalid token claims")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			log.Printf("Token expired at %v", time.Unix(int64(exp), 0))
			return "", errors.New("token has expired")
		}
	} else {
		log.Printf("Token missing expiration claim")
		return "", errors.New("invalid token: missing expiration")
	}

	uid, ok := claims["uid"].(string)
	if !ok {
		log.Printf("Missing or invalid UID in token claims")
		return "", errors.New("missing or invalid UID in token")
	}

	if strings.TrimSpace(uid) == "" {
		log.Printf("Empty UID in token claims")
		return "", errors.New("empty UID in token")
	}

	return uid, nil
}

// extractAndValidateToken extracts and validates the authorization token
// Usage example: uid, err := extractAndValidateToken(c)
func ExtractAndValidateToken(c *fiber.Ctx) (string, error) {
	token := c.Get("Authorization")
	if token == "" {
		return "", c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": errTokenRequired,
		})
	}

	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return "", c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": errInvalidToken,
		})
	}

	return uid, nil
}

// handleFirestoreError handles Firestore-specific errors
func handleFirestoreError(c *fiber.Ctx, err error) error {
	if status.Code(err) == codes.AlreadyExists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": errCheckExisting,
	})
}

// isRetryableError determines if an error is retryable
func isRetryableError(err error) bool {
	code := status.Code(err)
	return code == codes.DeadlineExceeded ||
		code == codes.Unavailable ||
		code == codes.Internal
}

// checkExistingExperience checks if a school experience already exists for the user
func checkExistingExperience(ctx context.Context, uid string) error {
	query := services.FirestoreClient.Collection("user_school_experiences").
		Where("uid", "==", uid).
		Limit(1).
		Documents(ctx)
	defer query.Stop()

	_, err := query.Next()
	if err == nil {
		return status.Error(codes.AlreadyExists, errExperienceExists)
	}
	if err != iterator.Done {
		return err
	}
	return nil
}

// createNewExperience creates a new UserSchoolExperience instance
func createNewExperience(uid string) *models.UserSchoolExperience {
	currentTime := time.Now().UTC()
	return &models.UserSchoolExperience{
		UID:          uid,
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
		Universities: make([]models.University, 0),
	}
}

// saveExperience saves the experience to Firestore with retry mechanism
func saveExperience(ctx context.Context, experience *models.UserSchoolExperience) error {
	backendData := mappers.MapUserSchoolExperienceFrontendToBackend(experience)

	// Implement retry mechanism
	maxRetries := 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		_, _, err := services.FirestoreClient.Collection("user_school_experiences").Add(ctx, backendData)
		if err == nil {
			return nil
		}

		// If this is not the last attempt and the error is retryable
		if attempt < maxRetries && isRetryableError(err) {
			time.Sleep(time.Duration(attempt*100) * time.Millisecond)
			continue
		}
		return err
	}
	return nil
}

// validateRequestBody validates the request body
func validateRequestBody(c *fiber.Ctx, data interface{}) error {
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errInvalidPayload,
		})
	}
	return nil
}

func getSchoolExperience(ctx context.Context, uid string) (*schoolExperienceDoc, error) {
	query := services.FirestoreClient.Collection("user_school_experiences").
		Where("uid", "==", uid).
		Limit(1).
		Documents(ctx)
	defer query.Stop()

	doc, err := query.Next()
	if err == iterator.Done {
		return nil, status.Error(codes.NotFound, errExperienceNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch document: %w", err)
	}

	var experience models.UserSchoolExperience
	if err := doc.DataTo(&experience); err != nil {
		return nil, fmt.Errorf("failed to parse document: %w", err)
	}

	return &schoolExperienceDoc{
		ref:        doc.Ref,
		experience: &experience,
	}, nil
}

// addUniversityTransaction adds a single university using a transaction
func addUniversityTransaction(ctx context.Context, doc *schoolExperienceDoc, universityData map[string]interface{}) error {
	return services.FirestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		newUniversity := mappers.MapFrontendToUniversity(universityData)
		newUniversity.ID = services.GenerateID()

		doc.experience.Universities = append(doc.experience.Universities, newUniversity)
		doc.experience.UpdatedAt = time.Now()

		backendData := mappers.MapUserSchoolExperienceFrontendToBackend(doc.experience)
		return tx.Set(doc.ref, backendData, firestore.MergeAll)
	})
}

func addUniversitiesTransaction(ctx context.Context, doc *schoolExperienceDoc, universitiesData []map[string]interface{}) error {
	return services.FirestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		for _, universityData := range universitiesData {
			newUniversity := mappers.MapFrontendToUniversity(universityData)
			newUniversity.ID = services.GenerateID()
			doc.experience.Universities = append(doc.experience.Universities, newUniversity)
		}

		doc.experience.UpdatedAt = time.Now()
		backendData := mappers.MapUserSchoolExperienceFrontendToBackend(doc.experience)
		return tx.Set(doc.ref, backendData, firestore.MergeAll)
	})
}
