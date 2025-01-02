package handlers

import (
	"context"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/mappers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"github.com/rogerjeasy/go-letusconnect/services"
	"google.golang.org/api/iterator"
)

type AddressHandler struct {
	AddressService *services.AddressService
}

func NewAddressHandler(addressService *services.AddressService) *AddressHandler {
	return &AddressHandler{
		AddressService: addressService,
	}
}

// CRUD operations for user addresses
func (a *AddressHandler) CreateUserAddress(c *fiber.Ctx) error {
	// Extract Authorization token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	// Validate the token and extract the UID
	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Parse the request body into a map
	var requestData map[string]interface{}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Map the frontend data to the UserAddress struct
	address := mappers.MapFrontendToUserAddress(requestData)
	address.UID = uid

	// Convert the address to firestore format
	backendAddress := mappers.MapUserAddressToBackend(address)

	// Add the address to the Firestore collection
	docRef, _, err := services.FirestoreClient.Collection("user_addresses").Add(context.Background(), backendAddress)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save address",
		})
	}

	// Set the ID field to the document ID
	backendAddress["id"] = docRef.ID

	// Update the document with the ID field
	_, err = docRef.Set(context.Background(), backendAddress, firestore.MergeAll)
	if err != nil {
		log.Printf("Error updating Firestore: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update address with document ID",
		})
	}

	// Map the backend address back to frontend format
	frontendAddress := mappers.MapBackendToFrontend(backendAddress)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Address created successfully",
		"address": frontendAddress,
	})
}

func (a *AddressHandler) UpdateUserAddress(c *fiber.Ctx) error {
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

	// Get Address ID
	addressID := c.Params("id")
	if addressID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Address ID is required",
		})
	}

	// Parse request body into a map
	var requestData map[string]interface{}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Fetch current address from Firestore
	docRef := services.FirestoreClient.Collection("user_addresses").Doc(addressID)
	docSnapshot, err := docRef.Get(context.Background())
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Address not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch address",
		})
	}

	// Verify if the address belongs to the current user
	data := docSnapshot.Data()
	if data["uid"] != uid {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You are not authorized to update this address",
		})
	}

	// Map frontend data to UserAddress struct
	updatedAddress := mappers.MapFrontendToUserAddress(requestData)

	// Preserve UID
	updatedAddress.UID = uid

	// Map to backend format
	backendUpdates := mappers.MapUserAddressToBackend(updatedAddress)

	// Update Firestore document
	_, err = docRef.Set(context.Background(), backendUpdates, firestore.MergeAll)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update address",
		})
	}

	// Map the backend data to frontend format
	frontendAddress := mappers.MapBackendToFrontend(backendUpdates)

	// Return the updated address
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Address updated successfully",
		"address": frontendAddress,
	})
}

func (a *AddressHandler) GetUserAddress(c *fiber.Ctx) error {
	// Extract Authorization token
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	// Validate the token and extract the UID
	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Get addresses using the service function
	addresses, err := a.AddressService.GetUserAddresses(uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch addresses",
		})
	}

	var frontendAddresses []map[string]interface{}
	for _, address := range addresses {
		frontendAddress := mappers.MapUserAddressToFrontend(address)
		frontendAddresses = append(frontendAddresses, frontendAddress)
	}

	// Return a proper message if no addresses are found
	if len(frontendAddresses) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No addresses found for the user",
		})
	}

	// Return the addresses in frontend format
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "Addresses retrieved successfully",
		"addresses": frontendAddresses,
	})
}

func (a *AddressHandler) DeleteUserAddress(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	addressID := c.Params("id")
	if addressID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Address ID is required",
		})
	}

	// Verify the address belongs to the user
	docRef := services.FirestoreClient.Collection("user_addresses").Doc(addressID)
	docSnapshot, err := docRef.Get(context.Background())
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Address not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch address",
		})
	}

	// Check if UID matches the authenticated user
	data := docSnapshot.Data()
	if data["UID"] != uid {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You are not authorized to delete this address",
		})
	}

	_, err = docRef.Delete(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete address",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Address deleted successfully",
	})
}

// CRUD operations for user work experiences

func CreateUserWorkExperience(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	workExperience := new(models.UserWorkExperience)
	if err := c.BodyParser(workExperience); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	workExperience.UID = uid
	workExperience.CreatedAt = time.Now()

	_, _, err = services.FirestoreClient.Collection("user_work_experiences").Add(context.Background(), workExperience)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save work experience",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Work experience created successfully",
		"data":    workExperience,
	})
}

func UpdateUserWorkExperience(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	workExperienceID := c.Params("id")
	if workExperienceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Work experience ID is required",
		})
	}

	updatedData := new(models.UserWorkExperience)
	if err := c.BodyParser(updatedData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	docRef := services.FirestoreClient.Collection("user_work_experiences").Doc(workExperienceID)
	docSnapshot, err := docRef.Get(context.Background())
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Work experience not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch work experience",
		})
	}

	data := docSnapshot.Data()
	if data["UID"] != uid {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You are not authorized to update this work experience",
		})
	}

	updatedData.UID = uid
	updatedData.UpdatedAt = time.Now()

	_, err = docRef.Set(context.Background(), updatedData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update work experience",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Work experience updated successfully",
		"data":    updatedData,
	})
}

func GetUserWorkExperience(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	workExperienceID := c.Params("id")
	if workExperienceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Work experience ID is required",
		})
	}

	docRef := services.FirestoreClient.Collection("user_work_experiences").Doc(workExperienceID)
	docSnapshot, err := docRef.Get(context.Background())
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Work experience not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch work experience",
		})
	}

	data := docSnapshot.Data()
	if data["UID"] != uid {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You are not authorized to access this work experience",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Work experience retrieved successfully",
		"data":    data,
	})
}

func DeleteUserWorkExperience(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	workExperienceID := c.Params("id")
	if workExperienceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Work experience ID is required",
		})
	}

	docRef := services.FirestoreClient.Collection("user_work_experiences").Doc(workExperienceID)
	docSnapshot, err := docRef.Get(context.Background())
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Work experience not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch work experience",
		})
	}

	data := docSnapshot.Data()
	if data["UID"] != uid {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You are not authorized to delete this work experience",
		})
	}

	_, err = docRef.Delete(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete work experience",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Work experience deleted successfully",
	})
}

// CRUD operations for user school experiences
func CreateUserSchoolExperience(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	schoolExperience := new(models.UserSchoolExperience)
	if err := c.BodyParser(schoolExperience); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	schoolExperience.UID = uid
	schoolExperience.CreatedAt = time.Now()

	_, _, err = services.FirestoreClient.Collection("user_school_experiences").Add(context.Background(), schoolExperience)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save school experience",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "School experience created successfully",
		"data":    schoolExperience,
	})
}

func UpdateUserSchoolExperience(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	schoolExperienceID := c.Params("id")
	if schoolExperienceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "School experience ID is required",
		})
	}

	updatedData := new(models.UserSchoolExperience)
	if err := c.BodyParser(updatedData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	docRef := services.FirestoreClient.Collection("user_school_experiences").Doc(schoolExperienceID)
	docSnapshot, err := docRef.Get(context.Background())
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "School experience not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch school experience",
		})
	}

	data := docSnapshot.Data()
	if data["UID"] != uid {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You are not authorized to update this school experience",
		})
	}

	updatedData.UID = uid
	updatedData.UpdatedAt = time.Now()

	_, err = docRef.Set(context.Background(), updatedData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update school experience",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "School experience updated successfully",
		"data":    updatedData,
	})
}

func GetUserSchoolExperience(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	query := services.FirestoreClient.Collection("user_school_experiences").Where("UID", "==", uid).Documents(context.Background())
	defer query.Stop()

	var experiences []models.UserSchoolExperience
	for {
		doc, err := query.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve school experiences",
			})
		}

		var experience models.UserSchoolExperience
		if err := doc.DataTo(&experience); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to parse school experience data",
			})
		}
		experiences = append(experiences, experience)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "School experiences retrieved successfully",
		"data":    experiences,
	})
}

func DeleteUserSchoolExperience(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization token is required",
		})
	}

	uid, err := validateToken(strings.TrimPrefix(token, "Bearer "))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	schoolExperienceID := c.Params("id")
	if schoolExperienceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "School experience ID is required",
		})
	}

	docRef := services.FirestoreClient.Collection("user_school_experiences").Doc(schoolExperienceID)
	docSnapshot, err := docRef.Get(context.Background())
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "School experience not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch school experience",
		})
	}

	data := docSnapshot.Data()
	if data["UID"] != uid {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You are not authorized to delete this school experience",
		})
	}

	_, err = docRef.Delete(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete school experience",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "School experience deleted successfully",
	})
}
