package mappers

import (
	"github.com/rogerjeasy/go-letusconnect/models"
)

// MapUserFrontendToBackend maps frontend User data to Firestore-compatible format (snake_case)
func MapUserFrontendToBackend(user *models.User) map[string]interface{} {
	return map[string]interface{}{
		"uid":                   user.UID,
		"username":              user.Username,
		"first_name":            user.FirstName,
		"last_name":             user.LastName,
		"email":                 user.Email,
		"phone_number":          user.PhoneNumber,
		"profile_picture":       user.ProfilePicture,
		"bio":                   user.Bio,
		"role":                  user.Role,
		"graduation_year":       user.GraduationYear,
		"current_job_title":     user.CurrentJobTitle,
		"areas_of_expertise":    user.AreasOfExpertise,
		"interests":             user.Interests,
		"looking_for_mentor":    user.LookingForMentor,
		"willing_to_mentor":     user.WillingToMentor,
		"connections_made":      user.ConnectionsMade,
		"account_creation_date": user.AccountCreatedAt,
		"is_active":             user.IsActive,
		"is_verified":           user.IsVerified,
		"program":               user.Program,
		"date_of_birth":         user.DateOfBirth,
		"phone_code":            user.PhoneCode,
		"languages":             user.Languages,
		"skills":                user.Skills,
		"certifications":        user.Certifications,
		"projects":              user.Projects,
	}
}

// MapUserBackendToFrontend maps Firestore User data to frontend format (camelCase)
func MapUserBackendToFrontend(backendUser map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"uid":              getStringValue(backendUser, "uid"),
		"username":         getStringValue(backendUser, "username"),
		"firstName":        getStringValue(backendUser, "first_name"),
		"lastName":         getStringValue(backendUser, "last_name"),
		"email":            getStringValue(backendUser, "email"),
		"phoneNumber":      getStringValue(backendUser, "phone_number"),
		"profilePicture":   getStringValue(backendUser, "profile_picture"),
		"bio":              getStringValue(backendUser, "bio"),
		"role":             backendUser["role"],
		"graduationYear":   getIntValueSafe(backendUser, "graduation_year"),
		"currentJobTitle":  getStringValue(backendUser, "current_job_title"),
		"areasOfExpertise": backendUser["areas_of_expertise"],
		"interests":        backendUser["interests"],
		"lookingForMentor": backendUser["looking_for_mentor"],
		"willingToMentor":  backendUser["willing_to_mentor"],
		"connectionsMade":  getIntValueSafe(backendUser, "connections_made"),
		"accountCreatedAt": getStringValue(backendUser, "account_creation_date"),
		"isActive":         backendUser["is_active"],
		"isVerified":       backendUser["is_verified"],
		"program":          getStringValue(backendUser, "program"),
		"dateOfBirth":      getStringValue(backendUser, "date_of_birth"),
		"phoneCode":        getStringValue(backendUser, "phone_code"),
		"languages":        backendUser["languages"],
		"skills":           backendUser["skills"],
		"certifications":   backendUser["certifications"],
		"projects":         backendUser["projects"],
	}
}

// MapFrontendToUser maps frontend data to the User model
func MapFrontendToUser(data map[string]interface{}) models.User {
	return models.User{
		UID:              getStringValue(data, "uid"),
		Username:         getStringValue(data, "username"),
		FirstName:        getStringValue(data, "firstName"),
		LastName:         getStringValue(data, "lastName"),
		Email:            getStringValue(data, "email"),
		PhoneNumber:      getStringValue(data, "phoneNumber"),
		ProfilePicture:   getStringValue(data, "profilePicture"),
		Bio:              getStringValue(data, "bio"),
		Role:             getStringArrayValue(data, "role"),
		GraduationYear:   getIntValueSafe(data, "graduationYear"),
		CurrentJobTitle:  getStringValue(data, "currentJobTitle"),
		AreasOfExpertise: getStringArrayValue(data, "areasOfExpertise"),
		Interests:        getStringArrayValue(data, "interests"),
		LookingForMentor: getBoolValueSafe(data, "lookingForMentor"),
		WillingToMentor:  getBoolValueSafe(data, "willingToMentor"),
		ConnectionsMade:  getIntValueSafe(data, "connectionsMade"),
		AccountCreatedAt: getStringValue(data, "accountCreatedAt"),
		IsActive:         getBoolValueSafe(data, "isActive"),
		IsVerified:       getBoolValueSafe(data, "isVerified"),
		Password:         getStringValue(data, "password"),
		Program:          getStringValue(data, "program"),
		DateOfBirth:      getStringValue(data, "dateOfBirth"),
		PhoneCode:        getStringValue(data, "phoneCode"),
		Languages:        getStringArrayValue(data, "languages"),
		Skills:           getStringArrayValue(data, "skills"),
		Certifications:   getStringArrayValue(data, "certifications"),
		Projects:         getStringArrayValue(data, "projects"),
	}
}

// MapBackendToUser maps Firestore User data to the User model struct
func MapBackendToUser(data map[string]interface{}) models.User {
	return models.User{
		UID:              getStringValue(data, "uid"),
		Username:         getStringValue(data, "username"),
		FirstName:        getStringValue(data, "first_name"),
		LastName:         getStringValue(data, "last_name"),
		Email:            getStringValue(data, "email"),
		PhoneNumber:      getStringValue(data, "phone_number"),
		ProfilePicture:   getStringValue(data, "profile_picture"),
		Bio:              getStringValue(data, "bio"),
		Role:             getStringArrayValue(data, "role"),
		GraduationYear:   getIntValueSafe(data, "graduation_year"),
		CurrentJobTitle:  getStringValue(data, "current_job_title"),
		AreasOfExpertise: getStringArrayValue(data, "areas_of_expertise"),
		Interests:        getStringArrayValue(data, "interests"),
		LookingForMentor: getBoolValueSafe(data, "looking_for_mentor"),
		WillingToMentor:  getBoolValueSafe(data, "willing_to_mentor"),
		ConnectionsMade:  getIntValueSafe(data, "connections_made"),
		AccountCreatedAt: getStringValue(data, "account_creation_date"),
		IsActive:         getBoolValueSafe(data, "is_active"),
		IsVerified:       getBoolValueSafe(data, "is_verified"),
		Program:          getStringValue(data, "program"),
		DateOfBirth:      getStringValue(data, "date_of_birth"),
		PhoneCode:        getStringValue(data, "phone_code"),
		Languages:        getStringArrayValue(data, "languages"),
		Skills:           getStringArrayValue(data, "skills"),
		Certifications:   getStringArrayValue(data, "certifications"),
		Projects:         getStringArrayValue(data, "projects"),
	}
}

// MapUserToFrontend maps a User model struct to frontend format (camelCase)
func MapUserToFrontend(user *models.User) map[string]interface{} {
	return map[string]interface{}{
		"uid":              user.UID,
		"username":         user.Username,
		"firstName":        user.FirstName,
		"lastName":         user.LastName,
		"email":            user.Email,
		"phoneNumber":      user.PhoneNumber,
		"profilePicture":   user.ProfilePicture,
		"bio":              user.Bio,
		"role":             user.Role,
		"graduationYear":   user.GraduationYear,
		"currentJobTitle":  user.CurrentJobTitle,
		"areasOfExpertise": user.AreasOfExpertise,
		"interests":        user.Interests,
		"lookingForMentor": user.LookingForMentor,
		"willingToMentor":  user.WillingToMentor,
		"connectionsMade":  user.ConnectionsMade,
		"accountCreatedAt": user.AccountCreatedAt,
		"isActive":         user.IsActive,
		"isVerified":       user.IsVerified,
		"program":          user.Program,
		"dateOfBirth":      user.DateOfBirth,
		"phoneCode":        user.PhoneCode,
		"languages":        user.Languages,
		"skills":           user.Skills,
		"certifications":   user.Certifications,
		"projects":         user.Projects,
	}
}

// getBoolValueSafe safely retrieves a boolean value from a map
func getBoolValueSafe(data map[string]interface{}, key string) bool {
	if value, ok := data[key]; ok {
		if boolVal, isBool := value.(bool); isBool {
			return boolVal
		}
	}
	return false
}
