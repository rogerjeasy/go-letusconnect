package mappers

import "github.com/rogerjeasy/go-letusconnect/models"

func MapFrontendToBackend(user *models.User) map[string]interface{} {
	backend := make(map[string]interface{})

	if user.UID != "" {
		backend["uid"] = user.UID
	}
	if user.Username != "" {
		backend["username"] = user.Username
	}
	if user.FirstName != "" {
		backend["first_name"] = user.FirstName
	}
	if user.LastName != "" {
		backend["last_name"] = user.LastName
	}
	if user.Email != "" {
		backend["email"] = user.Email
	}
	if user.PhoneNumber != "" {
		backend["phone_number"] = user.PhoneNumber
	}
	if user.ProfilePicture != "" {
		backend["profile_picture"] = user.ProfilePicture
	}
	if user.Bio != "" {
		backend["bio"] = user.Bio
	}
	if len(user.Role) > 0 {
		backend["role"] = user.Role
	}
	if user.GraduationYear != 0 {
		backend["graduation_year"] = user.GraduationYear
	}
	if user.CurrentJobTitle != "" {
		backend["current_job_title"] = user.CurrentJobTitle
	}
	if len(user.AreasOfExpertise) > 0 {
		backend["areas_of_expertise"] = user.AreasOfExpertise
	}
	if len(user.Interests) > 0 {
		backend["interests"] = user.Interests
	}
	if user.LookingForMentor {
		backend["looking_for_mentor"] = user.LookingForMentor
	}
	if user.WillingToMentor {
		backend["willing_to_mentor"] = user.WillingToMentor
	}
	if user.ConnectionsMade != 0 {
		backend["connections_made"] = user.ConnectionsMade
	}
	if user.AccountCreatedAt != "" {
		backend["account_creation_date"] = user.AccountCreatedAt
	}
	if user.IsActive {
		backend["is_active"] = user.IsActive
	}
	if user.IsVerified {
		backend["is_verified"] = user.IsVerified
	}
	if user.Program != "" {
		backend["program"] = user.Program
	}
	if user.DateOfBirth != "" {
		backend["date_of_birth"] = user.DateOfBirth
	}
	if user.PhoneCode != "" {
		backend["phone_code"] = user.PhoneCode
	}
	if len(user.Languages) > 0 {
		backend["languages"] = user.Languages
	}
	if len(user.Skills) > 0 {
		backend["skills"] = user.Skills
	}
	if len(user.Certifications) > 0 {
		backend["certifications"] = user.Certifications
	}
	if len(user.Projects) > 0 {
		backend["projects"] = user.Projects
	}

	return backend
}

func MapBackendToFrontend(backendUser models.User) map[string]interface{} {
	return map[string]interface{}{
		"uid":              backendUser.UID,
		"username":         backendUser.Username,
		"firstName":        backendUser.FirstName,
		"lastName":         backendUser.LastName,
		"email":            backendUser.Email,
		"phoneNumber":      backendUser.PhoneNumber,
		"profilePicture":   backendUser.ProfilePicture,
		"bio":              backendUser.Bio,
		"role":             backendUser.Role,
		"graduationYear":   backendUser.GraduationYear,
		"currentJobTitle":  backendUser.CurrentJobTitle,
		"areasOfExpertise": backendUser.AreasOfExpertise,
		"interests":        backendUser.Interests,
		"lookingForMentor": backendUser.LookingForMentor,
		"willingToMentor":  backendUser.WillingToMentor,
		"connectionsMade":  backendUser.ConnectionsMade,
		"accountCreatedAt": backendUser.AccountCreatedAt,
		"isActive":         backendUser.IsActive,
		"isVerified":       backendUser.IsVerified,
		"program":          backendUser.Program,
		"dateOfBirth":      backendUser.DateOfBirth,
		"phoneCode":        backendUser.PhoneCode,
		"languages":        backendUser.Languages,
		"skills":           backendUser.Skills,
		"certifications":   backendUser.Certifications,
		"projects":         backendUser.Projects,
	}
}
