package mappers

import (
	"strconv"

	"github.com/rogerjeasy/go-letusconnect/models"
)

// MapUserAddressFrontendToBackend maps frontend data to Firestore-compatible format (snake_case)
func MapUserAddressFrontendToBackend(address *models.UserAddress) map[string]interface{} {
	return map[string]interface{}{
		"id":           address.ID,
		"uid":          address.UID,
		"street":       address.Street,
		"city":         address.City,
		"state":        address.State,
		"country":      address.Country,
		"postal_code":  address.PostalCode,
		"house_number": address.HouseNumber,
		"apartment":    address.Apartment,
		"region":       address.Region,
	}
}

// MapUserAddressBackendToFrontend maps Firestore data to frontend format (camelCase)
func MapUserAddressBackendToFrontend(backendAddress map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":          backendAddress["id"],
		"uid":         backendAddress["uid"],
		"street":      backendAddress["street"],
		"city":        backendAddress["city"],
		"state":       backendAddress["state"],
		"country":     backendAddress["country"],
		"postalCode":  backendAddress["postal_code"],
		"houseNumber": backendAddress["house_number"],
		"apartment":   backendAddress["apartment"],
		"region":      backendAddress["region"],
	}
}

func MapFrontendToUserAddress(data map[string]interface{}) models.UserAddress {
	address := models.UserAddress{
		Street:      getStringValue(data, "street"),
		City:        getStringValue(data, "city"),
		State:       getStringValue(data, "state"),
		Country:     getStringValue(data, "country"),
		PostalCode:  getIntValueSafe(data, "postalCode"),
		HouseNumber: getIntValueSafe(data, "houseNumber"),
		Apartment:   getStringValue(data, "apartment"),
		Region:      getStringValue(data, "region"),
		ID:          getStringValue(data, "id"),
	}

	return address
}

// Helper function to safely get string values
func getStringValue(data map[string]interface{}, key string) string {
	if value, ok := data[key]; ok {
		if strVal, isString := value.(string); isString {
			return strVal
		}
	}
	return ""
}

// Helper function to safely get integer values
func getIntValueSafe(data map[string]interface{}, key string) int {
	if value, ok := data[key]; ok {
		switch v := value.(type) {
		case float64:
			return int(v) // Convert float64 to int for numeric fields from JSON
		case int:
			return v
		case string:
			// Attempt to convert string to int
			if intValue, err := strconv.Atoi(v); err == nil {
				return intValue
			}
		}
	}
	return 0
}
