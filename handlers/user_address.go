package handlers

import (
	"context"
	"strings"
	"time"

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

func (a *AddressHandler) CreateUserAddress(c *fiber.Ctx) error {
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

	createdAddress, err := a.AddressService.CreateUserAddress(uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Map the address to frontend format
	// frontendAddress := mappers.MapBackendToFrontend(mappers.MapUserAddressToBackend(createdAddress))
	addressUser := mappers.MapUserAddressToFrontend(createdAddress)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Address created successfully",
		"data":    addressUser,
	})
}

func (a *AddressHandler) UpdateUserAddress(c *fiber.Ctx) error {
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

	var requestData map[string]interface{}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	updatedAddress := mappers.MapFrontendToUserAddress(requestData)

	result, err := a.AddressService.UpdateUserAddress(addressID, uid, updatedAddress)
	if err != nil {
		if err.Error() == "unauthorized to update this address" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if strings.Contains(err.Error(), "failed to fetch address") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Address not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	frontendAddress := mappers.MapUserAddressToFrontend(result)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Address updated successfully",
		"data":    frontendAddress,
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

	if a.AddressService == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Address service is not initialized",
		})
	}

	// Get addresses using the service function
	addresses, err := a.AddressService.GetUserAddresses(uid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch addresses: " + err.Error(),
		})
	}

	// Handle empty results
	if len(addresses) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":   "No addresses found for the user",
			"addresses": []models.UserAddress{},
		})
	}

	// Map addresses to frontend format
	var frontendAddresses []map[string]interface{}
	for _, address := range addresses {
		frontendAddress := mappers.MapUserAddressToFrontend(address)
		frontendAddresses = append(frontendAddresses, frontendAddress)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Addresses retrieved successfully",
		"data":    frontendAddresses,
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
	docRef := services.Firestore.Collection("user_addresses").Doc(addressID)
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

	_, _, err = services.Firestore.Collection("user_work_experiences").Add(context.Background(), workExperience)
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

	docRef := services.Firestore.Collection("user_work_experiences").Doc(workExperienceID)
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

	docRef := services.Firestore.Collection("user_work_experiences").Doc(workExperienceID)
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

	docRef := services.Firestore.Collection("user_work_experiences").Doc(workExperienceID)
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

	_, _, err = services.Firestore.Collection("user_school_experiences").Add(context.Background(), schoolExperience)
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

	docRef := services.Firestore.Collection("user_school_experiences").Doc(schoolExperienceID)
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

	query := services.Firestore.Collection("user_school_experiences").Where("UID", "==", uid).Documents(context.Background())
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

	docRef := services.Firestore.Collection("user_school_experiences").Doc(schoolExperienceID)
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
