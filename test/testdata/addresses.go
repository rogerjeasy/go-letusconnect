package testdata

import (
	"github.com/rogerjeasy/go-letusconnect/models"
)

var (
	// TestAddresses contains sample address data for testing
	TestAddresses = []models.UserAddress{
		{
			ID:          "addr1",
			UID:         "user1",
			Street:      "123 Main St",
			City:        "New York",
			State:       "NY",
			Country:     "USA",
			PostalCode:  10001,
			HouseNumber: 123,
			Apartment:   "4B",
			Region:      "Manhattan",
		},
		{
			ID:          "addr2",
			UID:         "user1",
			Street:      "456 Oak Ave",
			City:        "Los Angeles",
			State:       "CA",
			Country:     "USA",
			PostalCode:  90001,
			HouseNumber: 456,
			Apartment:   "",
			Region:      "Downtown",
		},
	}

	// TestBackendAddresses contains sample Firestore data format
	TestBackendAddresses = []map[string]interface{}{
		{
			"id":           "addr1",
			"uid":          "user1",
			"street":       "123 Main St",
			"city":         "New York",
			"state":        "NY",
			"country":      "USA",
			"postal_code":  10001,
			"house_number": 123,
			"apartment":    "4B",
			"region":       "Manhattan",
		},
		{
			"id":           "addr2",
			"uid":          "user1",
			"street":       "456 Oak Ave",
			"city":         "Los Angeles",
			"state":        "CA",
			"country":      "USA",
			"postal_code":  90001,
			"house_number": 456,
			"apartment":    "",
			"region":       "Downtown",
		},
	}

	// Test user IDs
	TestUIDs = []string{
		"user1",
		"user2",
		"user3",
	}
)
