package mappers

import (
	"testing"

	"github.com/rogerjeasy/go-letusconnect/models"
	"github.com/stretchr/testify/assert"
)

func TestMapUserFrontendToBackend(t *testing.T) {
	user := &models.User{
		UID:              "123",
		Username:         "testuser",
		FirstName:        "John",
		LastName:         "Doe",
		Email:            "john.doe@example.com",
		PhoneNumber:      "1234567890",
		ProfilePicture:   "profile.jpg",
		Bio:              "Hello, I'm John!",
		Role:             []string{"user"},
		GraduationYear:   2020,
		CurrentJobTitle:  "Software Engineer",
		AreasOfExpertise: []string{"Go", "JavaScript"},
		Interests:        []string{"Programming", "Reading"},
		LookingForMentor: true,
		WillingToMentor:  false,
		ConnectionsMade:  10,
		AccountCreatedAt: "2020-01-01",
		IsActive:         true,
		IsVerified:       true,
		Program:          "Computer Science",
		DateOfBirth:      "1990-01-01",
		PhoneCode:        "+1",
		Languages:        []string{"English", "Spanish"},
		Skills:           []string{"Go", "JavaScript"},
		Certifications:   []string{"AWS Certified"},
		Projects:         []string{"Project 1", "Project 2"},
		IsOnline:         true,
		IsPrivate:        false,
	}

	expected := map[string]interface{}{
		"uid":                   "123",
		"username":              "testuser",
		"first_name":            "John",
		"last_name":             "Doe",
		"email":                 "john.doe@example.com",
		"phone_number":          "1234567890",
		"profile_picture":       "profile.jpg",
		"bio":                   "Hello, I'm John!",
		"role":                  []string{"user"},
		"graduation_year":       2020,
		"current_job_title":     "Software Engineer",
		"areas_of_expertise":    []string{"Go", "JavaScript"},
		"interests":             []string{"Programming", "Reading"},
		"looking_for_mentor":    true,
		"willing_to_mentor":     false,
		"connections_made":      10,
		"account_creation_date": "2020-01-01",
		"is_active":             true,
		"is_verified":           true,
		"program":               "Computer Science",
		"date_of_birth":         "1990-01-01",
		"phone_code":            "+1",
		"languages":             []string{"English", "Spanish"},
		"skills":                []string{"Go", "JavaScript"},
		"certifications":        []string{"AWS Certified"},
		"projects":              []string{"Project 1", "Project 2"},
		"is_online":             true,
		"is_private":            false,
	}

	result := MapUserFrontendToBackend(user)
	assert.Equal(t, expected, result)
}

func TestMapUserBackendToFrontend(t *testing.T) {
	backendUser := map[string]interface{}{
		"uid":                   "123",
		"username":              "testuser",
		"first_name":            "John",
		"last_name":             "Doe",
		"email":                 "john.doe@example.com",
		"phone_number":          "1234567890",
		"profile_picture":       "profile.jpg",
		"bio":                   "Hello, I'm John!",
		"role":                  []interface{}{"user"},
		"graduation_year":       2020,
		"current_job_title":     "Software Engineer",
		"areas_of_expertise":    []interface{}{"Go", "JavaScript"},
		"interests":             []interface{}{"Programming", "Reading"},
		"looking_for_mentor":    true,
		"willing_to_mentor":     false,
		"connections_made":      10,
		"account_creation_date": "2020-01-01",
		"is_active":             true,
		"is_verified":           true,
		"program":               "Computer Science",
		"date_of_birth":         "1990-01-01",
		"phone_code":            "+1",
		"languages":             []interface{}{"English", "Spanish"},
		"skills":                []interface{}{"Go", "JavaScript"},
		"certifications":        []interface{}{"AWS Certified"},
		"projects":              []interface{}{"Project 1", "Project 2"},
		"is_online":             true,
		"is_private":            false,
	}

	expected := map[string]interface{}{
		"uid":              "123",
		"username":         "testuser",
		"firstName":        "John",
		"lastName":         "Doe",
		"email":            "john.doe@example.com",
		"phoneNumber":      "1234567890",
		"profilePicture":   "profile.jpg",
		"bio":              "Hello, I'm John!",
		"role":             []interface{}{"user"},
		"graduationYear":   2020,
		"currentJobTitle":  "Software Engineer",
		"areasOfExpertise": []interface{}{"Go", "JavaScript"},
		"interests":        []interface{}{"Programming", "Reading"},
		"lookingForMentor": true,
		"willingToMentor":  false,
		"connectionsMade":  10,
		"accountCreatedAt": "2020-01-01",
		"isActive":         true,
		"isVerified":       true,
		"program":          "Computer Science",
		"dateOfBirth":      "1990-01-01",
		"phoneCode":        "+1",
		"languages":        []interface{}{"English", "Spanish"},
		"skills":           []interface{}{"Go", "JavaScript"},
		"certifications":   []interface{}{"AWS Certified"},
		"projects":         []interface{}{"Project 1", "Project 2"},
		"isOnline":         true,
		"isPrivate":        false,
	}

	result := MapUserBackendToFrontend(backendUser)
	assert.Equal(t, expected, result)
}

func TestMapFrontendToUser(t *testing.T) {
	data := map[string]interface{}{
		"uid":              "123",
		"username":         "testuser",
		"firstName":        "John",
		"lastName":         "Doe",
		"email":            "john.doe@example.com",
		"phoneNumber":      "1234567890",
		"profilePicture":   "profile.jpg",
		"bio":              "Hello, I'm John!",
		"role":             []interface{}{"user"},
		"graduationYear":   2020,
		"currentJobTitle":  "Software Engineer",
		"areasOfExpertise": []interface{}{"Go", "JavaScript"},
		"interests":        []interface{}{"Programming", "Reading"},
		"lookingForMentor": true,
		"willingToMentor":  false,
		"connectionsMade":  10,
		"accountCreatedAt": "2020-01-01",
		"isActive":         true,
		"isVerified":       true,
		"password":         "password",
		"program":          "Computer Science",
		"dateOfBirth":      "1990-01-01",
		"phoneCode":        "+1",
		"languages":        []interface{}{"English", "Spanish"},
		"skills":           []interface{}{"Go", "JavaScript"},
		"certifications":   []interface{}{"AWS Certified"},
		"projects":         []interface{}{"Project 1", "Project 2"},
		"isOnline":         true,
		"isPrivate":        false,
	}

	expected := models.User{
		UID:              "123",
		Username:         "testuser",
		FirstName:        "John",
		LastName:         "Doe",
		Email:            "john.doe@example.com",
		PhoneNumber:      "1234567890",
		ProfilePicture:   "profile.jpg",
		Bio:              "Hello, I'm John!",
		Role:             []string{"user"},
		GraduationYear:   2020,
		CurrentJobTitle:  "Software Engineer",
		AreasOfExpertise: []string{"Go", "JavaScript"},
		Interests:        []string{"Programming", "Reading"},
		LookingForMentor: true,
		WillingToMentor:  false,
		ConnectionsMade:  10,
		AccountCreatedAt: "2020-01-01",
		IsActive:         true,
		IsVerified:       true,
		Password:         "password",
		Program:          "Computer Science",
		DateOfBirth:      "1990-01-01",
		PhoneCode:        "+1",
		Languages:        []string{"English", "Spanish"},
		Skills:           []string{"Go", "JavaScript"},
		Certifications:   []string{"AWS Certified"},
		Projects:         []string{"Project 1", "Project 2"},
		IsOnline:         true,
		IsPrivate:        false,
	}

	result := MapFrontendToUser(data)
	assert.Equal(t, expected, result)
}

func TestMapBackendToUser(t *testing.T) {
	data := map[string]interface{}{
		"uid":                   "123",
		"username":              "testuser",
		"first_name":            "John",
		"last_name":             "Doe",
		"email":                 "john.doe@example.com",
		"phone_number":          "1234567890",
		"profile_picture":       "profile.jpg",
		"bio":                   "Hello, I'm John!",
		"role":                  []interface{}{"user"},
		"graduation_year":       2020,
		"current_job_title":     "Software Engineer",
		"areas_of_expertise":    []interface{}{"Go", "JavaScript"},
		"interests":             []interface{}{"Programming", "Reading"},
		"looking_for_mentor":    true,
		"willing_to_mentor":     false,
		"connections_made":      10,
		"account_creation_date": "2020-01-01",
		"is_active":             true,
		"is_verified":           true,
		"program":               "Computer Science",
		"date_of_birth":         "1990-01-01",
		"phone_code":            "+1",
		"languages":             []interface{}{"English", "Spanish"},
		"skills":                []interface{}{"Go", "JavaScript"},
		"certifications":        []interface{}{"AWS Certified"},
		"projects":              []interface{}{"Project 1", "Project 2"},
		"is_online":             true,
		"is_private":            false,
	}

	expected := models.User{
		UID:              "123",
		Username:         "testuser",
		FirstName:        "John",
		LastName:         "Doe",
		Email:            "john.doe@example.com",
		PhoneNumber:      "1234567890",
		ProfilePicture:   "profile.jpg",
		Bio:              "Hello, I'm John!",
		Role:             []string{"user"},
		GraduationYear:   2020,
		CurrentJobTitle:  "Software Engineer",
		AreasOfExpertise: []string{"Go", "JavaScript"},
		Interests:        []string{"Programming", "Reading"},
		LookingForMentor: true,
		WillingToMentor:  false,
		ConnectionsMade:  10,
		AccountCreatedAt: "2020-01-01",
		IsActive:         true,
		IsVerified:       true,
		Program:          "Computer Science",
		DateOfBirth:      "1990-01-01",
		PhoneCode:        "+1",
		Languages:        []string{"English", "Spanish"},
		Skills:           []string{"Go", "JavaScript"},
		Certifications:   []string{"AWS Certified"},
		Projects:         []string{"Project 1", "Project 2"},
		IsOnline:         true,
		IsPrivate:        false,
	}

	result := MapBackendToUser(data)
	assert.Equal(t, expected, result)
}

func TestMapUserToFrontend(t *testing.T) {
	user := &models.User{
		UID:              "123",
		Username:         "testuser",
		FirstName:        "John",
		LastName:         "Doe",
		Email:            "john.doe@example.com",
		PhoneNumber:      "1234567890",
		ProfilePicture:   "profile.jpg",
		Bio:              "Hello, I'm John!",
		Role:             []string{"user"},
		GraduationYear:   2020,
		CurrentJobTitle:  "Software Engineer",
		AreasOfExpertise: []string{"Go", "JavaScript"},
		Interests:        []string{"Programming", "Reading"},
		LookingForMentor: true,
		WillingToMentor:  false,
		ConnectionsMade:  10,
		AccountCreatedAt: "2020-01-01",
		IsActive:         true,
		IsVerified:       true,
		Program:          "Computer Science",
		DateOfBirth:      "1990-01-01",
		PhoneCode:        "+1",
		Languages:        []string{"English", "Spanish"},
		Skills:           []string{"Go", "JavaScript"},
		Certifications:   []string{"AWS Certified"},
		Projects:         []string{"Project 1", "Project 2"},
		IsOnline:         true,
		IsPrivate:        false,
	}

	expected := map[string]interface{}{
		"uid":              "123",
		"username":         "testuser",
		"firstName":        "John",
		"lastName":         "Doe",
		"email":            "john.doe@example.com",
		"phoneNumber":      "1234567890",
		"profilePicture":   "profile.jpg",
		"bio":              "Hello, I'm John!",
		"role":             []string{"user"},
		"graduationYear":   2020,
		"currentJobTitle":  "Software Engineer",
		"areasOfExpertise": []string{"Go", "JavaScript"},
		"interests":        []string{"Programming", "Reading"},
		"lookingForMentor": true,
		"willingToMentor":  false,
		"connectionsMade":  10,
		"accountCreatedAt": "2020-01-01",
		"isActive":         true,
		"isVerified":       true,
		"program":          "Computer Science",
		"dateOfBirth":      "1990-01-01",
		"phoneCode":        "+1",
		"languages":        []string{"English", "Spanish"},
		"skills":           []string{"Go", "JavaScript"},
		"certifications":   []string{"AWS Certified"},
		"projects":         []string{"Project 1", "Project 2"},
		"isOnline":         true,
		"isPrivate":        false,
	}

	result := MapUserToFrontend(user)
	assert.Equal(t, expected, result)
}

func TestGetBoolValueSafe(t *testing.T) {
	data := map[string]interface{}{
		"trueValue":  true,
		"falseValue": false,
		"nonBool":    "not a bool",
	}

	assert.True(t, getBoolValueSafe(data, "trueValue"))
	assert.False(t, getBoolValueSafe(data, "falseValue"))
	assert.False(t, getBoolValueSafe(data, "nonBool"))
	assert.False(t, getBoolValueSafe(data, "nonExistentKey"))
}
