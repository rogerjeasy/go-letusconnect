package services

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"google.golang.org/api/iterator"
)

type AddressService struct {
	FirestoreClient *firestore.Client
}

func NewAddressService(client *firestore.Client) *AddressService {
	return &AddressService{
		FirestoreClient: client,
	}
}

func (a *AddressService) CreateUserAddress(uid string) (models.UserAddress, error) {
	if a.FirestoreClient == nil {
		return models.UserAddress{}, errors.New("firestore client is not initialized")
	}

	// Initialize UserAddress with only UID
	address := models.UserAddress{
		UID: uid,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	backendAddress := mappers.MapUserAddressToBackend(address)

	docRef, _, err := a.FirestoreClient.Collection("user_addresses").Add(ctx, backendAddress)
	if err != nil {
		return models.UserAddress{}, errors.New("failed to create address: " + err.Error())
	}

	backendAddress["id"] = docRef.ID

	_, err = docRef.Set(ctx, backendAddress, firestore.MergeAll)
	if err != nil {
		return models.UserAddress{}, errors.New("failed to update address with ID: " + err.Error())
	}

	createdAddress := mappers.MapBackendToUserAddress(backendAddress)
	return createdAddress, nil
}

func (a *AddressService) GetUserAddresses(uid string) ([]models.UserAddress, error) {
	if a.FirestoreClient == nil {
		return nil, errors.New("firestore client is not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Query Firestore for user addresses
	iter := a.FirestoreClient.Collection("user_addresses").Where("uid", "==", uid).Documents(ctx)
	defer iter.Stop()

	var addresses []models.UserAddress
	hasAddress := false

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.New("error iterating through addresses: " + err.Error())
		}
		// Get the raw address data and ensure ID is included
		backendAddress := doc.Data()
		backendAddress["id"] = doc.Ref.ID
		// Convert to UserAddress struct
		address := mappers.MapBackendToUserAddress(backendAddress)
		addresses = append(addresses, address)
		hasAddress = true
	}

	// If no addresses found, create a new one
	if !hasAddress {
		newAddress, err := a.CreateUserAddress(uid)
		if err != nil {
			return nil, errors.New("failed to create initial address: " + err.Error())
		}
		addresses = append(addresses, newAddress)
	}

	return addresses, nil
}

func (a *AddressService) UpdateUserAddress(addressID string, uid string, updatedAddress models.UserAddress) (models.UserAddress, error) {
	if a.FirestoreClient == nil {
		return models.UserAddress{}, errors.New("firestore client is not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get document reference
	docRef := a.FirestoreClient.Collection("user_addresses").Doc(addressID)

	// Get current address
	docSnapshot, err := docRef.Get(ctx)
	if err != nil {
		return models.UserAddress{}, errors.New("failed to fetch address: " + err.Error())
	}

	data := docSnapshot.Data()
	if data["uid"] != uid {
		return models.UserAddress{}, errors.New("unauthorized to update this address")
	}

	updatedAddress.UID = uid
	updatedAddress.ID = addressID

	// Convert to backend format
	backendUpdates := mappers.MapUserAddressToBackend(updatedAddress)

	_, err = docRef.Set(ctx, backendUpdates, firestore.MergeAll)
	if err != nil {
		return models.UserAddress{}, errors.New("failed to update address: " + err.Error())
	}

	// Return the updated address
	return updatedAddress, nil
}

func (a *AddressService) GetUserSchoolExperience(uid string) (*models.UserSchoolExperience, error) {
	ctx := context.Background()

	// Query the Firestore collection to find school experience document
	query := a.FirestoreClient.Collection("school_experiences").Where("uid", "==", uid).Limit(1).Documents(ctx)
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
func (a *AddressService) GetUserWorkExperience(uid string) (*models.UserWorkExperience, error) {
	ctx := context.Background()

	// Query the Firestore collection to find work experience document
	query := a.FirestoreClient.Collection("work_experiences").Where("uid", "==", uid).Limit(1).Documents(ctx)
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
