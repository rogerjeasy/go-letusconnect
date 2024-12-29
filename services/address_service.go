package services

import (
	"context"
	"errors"
	"time"

	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"google.golang.org/api/iterator"
)

// GetUserAddresses retrieves all addresses for a given user UID
func GetUserAddresses(uid string) ([]models.UserAddress, error) {
	ctx := context.Background()

	// Query Firestore for user addresses
	docRef := FirestoreClient.Collection("user_addresses").Where("uid", "==", uid).Documents(ctx)
	defer docRef.Stop()

	var addresses []models.UserAddress
	for {
		doc, err := docRef.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.New("failed to fetch addresses")
		}

		// Get the raw address data
		backendAddress := doc.Data()
		// Convert to frontend format
		frontendAddress := mappers.MapBackendToUserAddress(backendAddress)
		addresses = append(addresses, frontendAddress)
	}

	return addresses, nil
}

// GetUserSchoolExperience fetches the school experience for a given user UID
func GetUserSchoolExperience(uid string) (*models.UserSchoolExperience, error) {
	ctx := context.Background()

	// Query the Firestore collection to find school experience document
	query := FirestoreClient.Collection("school_experiences").Where("uid", "==", uid).Limit(1).Documents(ctx)
	defer query.Stop()

	doc, err := query.Next()
	if err == iterator.Done {
		return nil, nil // Return nil without error if no school experience found
	}
	if err != nil {
		return nil, errors.New("failed to fetch school experience data")
	}

	// Extract the school experience data
	data := doc.Data()

	// Parse universities array
	universitiesData, ok := data["universities"].([]interface{})
	if !ok {
		return nil, errors.New("invalid universities data format")
	}

	var universities []models.University
	for _, uniInterface := range universitiesData {
		uniData, ok := uniInterface.(map[string]interface{})
		if !ok {
			continue
		}

		// Convert awards and extracurriculars to string slices
		awards := convertInterfaceToStringSlice(uniData["awards"])
		extracurriculars := convertInterfaceToStringSlice(uniData["extracurriculars"])

		university := models.University{
			ID:               uniData["id"].(string),
			Name:             uniData["name"].(string),
			Program:          uniData["program"].(string),
			Country:          uniData["country"].(string),
			City:             uniData["city"].(string),
			StartYear:        int(uniData["start_year"].(int64)),
			EndYear:          int(uniData["end_year"].(int64)),
			Degree:           uniData["degree"].(string),
			Experience:       uniData["experience"].(string),
			Awards:           awards,
			Extracurriculars: extracurriculars,
		}
		universities = append(universities, university)
	}

	// Create the school experience object
	schoolExp := &models.UserSchoolExperience{
		UID:          uid,
		CreatedAt:    data["created_at"].(time.Time),
		UpdatedAt:    data["updated_at"].(time.Time),
		Universities: universities,
	}

	return schoolExp, nil
}

// GetUserWorkExperience fetches the work experience for a given user UID
func GetUserWorkExperience(uid string) (*models.UserWorkExperience, error) {
	ctx := context.Background()

	// Query the Firestore collection to find work experience document
	query := FirestoreClient.Collection("work_experiences").Where("uid", "==", uid).Limit(1).Documents(ctx)
	defer query.Stop()

	doc, err := query.Next()
	if err == iterator.Done {
		return nil, nil // Return nil without error if no work experience found
	}
	if err != nil {
		return nil, errors.New("failed to fetch work experience data")
	}

	// Extract the work experience data
	data := doc.Data()

	// Parse work experiences array
	workExpsData, ok := data["work_experiences"].([]interface{})
	if !ok {
		return nil, errors.New("invalid work experiences data format")
	}

	var workExperiences []models.WorkExperience
	for _, workInterface := range workExpsData {
		workData, ok := workInterface.(map[string]interface{})
		if !ok {
			continue
		}

		// Convert responsibilities and achievements to string slices
		responsibilities := convertInterfaceToStringSlice(workData["responsibilities"])
		achievements := convertInterfaceToStringSlice(workData["achievements"])

		// Parse dates
		startDate := workData["start_date"].(time.Time)
		endDate := workData["end_date"].(time.Time)

		workExp := models.WorkExperience{
			ID:               workData["id"].(string),
			Company:          workData["company"].(string),
			Position:         workData["position"].(string),
			City:             workData["city"].(string),
			Country:          workData["country"].(string),
			StartDate:        startDate,
			EndDate:          endDate,
			Responsibilities: responsibilities,
			Achievements:     achievements,
		}
		workExperiences = append(workExperiences, workExp)
	}

	// Create the work experience object
	workExp := &models.UserWorkExperience{
		ID:              data["id"].(string),
		UID:             uid,
		CreatedAt:       data["created_at"].(time.Time),
		UpdatedAt:       data["updated_at"].(time.Time),
		WorkExperiences: workExperiences,
	}

	return workExp, nil
}

// Helper function to convert interface{} slice to string slice
func convertInterfaceToStringSlice(input interface{}) []string {
	var result []string
	if input == nil {
		return result
	}

	interfaceSlice, ok := input.([]interface{})
	if !ok {
		return result
	}

	for _, item := range interfaceSlice {
		if str, ok := item.(string); ok {
			result = append(result, str)
		}
	}
	return result
}
