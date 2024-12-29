package mappers

import (
	"github.com/rogerjeasy/go-letusconnect/models"
)

// 1. Frontend (JSON) to Go struct
// Input: map with camelCase keys (e.g., "houseNumber")
// Output: Go struct with PascalCase fields (e.g., HouseNumber)
func MapFrontendToUserAddress(data map[string]interface{}) models.UserAddress {
	return models.UserAddress{
		ID:          getStringValue(data, "id"),
		UID:         getStringValue(data, "uid"),
		Street:      getStringValue(data, "street"),
		City:        getStringValue(data, "city"),
		State:       getStringValue(data, "state"),
		Country:     getStringValue(data, "country"),
		PostalCode:  getIntValueSafe(data, "postalCode"),
		HouseNumber: getIntValueSafe(data, "houseNumber"),
		Apartment:   getStringValue(data, "apartment"),
		Region:      getStringValue(data, "region"),
	}
}

// 2. Go struct to Firestore
// Input: Go struct with PascalCase fields (e.g., HouseNumber)
// Output: map with snake_case keys (e.g., "house_number")
func MapUserAddressToBackend(address models.UserAddress) map[string]interface{} {
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

// 3. Firestore to Frontend
// Input: map with snake_case keys (e.g., "house_number")
// Output: map with camelCase keys (e.g., "houseNumber")
func MapBackendToFrontend(data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":          data["id"],
		"uid":         data["uid"],
		"street":      data["street"],
		"city":        data["city"],
		"state":       data["state"],
		"country":     data["country"],
		"postalCode":  data["postal_code"],
		"houseNumber": data["house_number"],
		"apartment":   data["apartment"],
		"region":      data["region"],
	}
}

// 4. Firestore to Go struct
// Input: map with snake_case keys (e.g., "house_number")
// Output: Go struct with PascalCase fields (e.g., HouseNumber)
func MapBackendToUserAddress(data map[string]interface{}) models.UserAddress {
	return models.UserAddress{
		ID:          getStringValue(data, "id"),
		UID:         getStringValue(data, "uid"),
		Street:      getStringValue(data, "street"),
		City:        getStringValue(data, "city"),
		State:       getStringValue(data, "state"),
		Country:     getStringValue(data, "country"),
		PostalCode:  getIntValueSafe(data, "postal_code"),
		HouseNumber: getIntValueSafe(data, "house_number"),
		Apartment:   getStringValue(data, "apartment"),
		Region:      getStringValue(data, "region"),
	}
}

// Go struct (PascalCase) -> Frontend (camelCase)
func MapUserAddressToFrontend(address models.UserAddress) map[string]interface{} {
	return map[string]interface{}{
		"id":          address.ID,
		"uid":         address.UID,
		"street":      address.Street,
		"city":        address.City,
		"state":       address.State,
		"country":     address.Country,
		"postalCode":  address.PostalCode,
		"houseNumber": address.HouseNumber,
		"apartment":   address.Apartment,
		"region":      address.Region,
	}
}
