package mappers

import (
	"log"
	"time"

	"github.com/rogerjeasy/go-letusconnect/models"
)

// MapInvitedUserFrontendToGo maps frontend InvitedUser data to Go struct format
func MapInvitedUserFrontendToGo(data map[string]interface{}) models.InvitedUser {
	return models.InvitedUser{
		UserID:         getStringValue(data, "userId"),
		Username:       getStringValue(data, "username"),
		Email:          getStringValue(data, "email"),
		ProfilePicture: getStringValue(data, "profilePicture"),
		Role:           getStringValue(data, "role"),
		JoinedAt:       getTimeValue(data, "joinedAt"),
	}
}

// MapInvitedUserGoToFirestore maps Go struct InvitedUser data to Firestore format
func MapInvitedUserGoToFirestore(user models.InvitedUser) map[string]interface{} {
	return map[string]interface{}{
		"user_id":         user.UserID,
		"username":        user.Username,
		"email":           user.Email,
		"profile_picture": user.ProfilePicture,
		"role":            user.Role,
		"joined_at":       user.JoinedAt,
	}
}

func MapInvitedUserGoToFrontend(user models.InvitedUser) map[string]interface{} {
	return map[string]interface{}{
		"userId":         user.UserID,
		"username":       user.Username,
		"email":          user.Email,
		"profilePicture": user.ProfilePicture,
		"role":           user.Role,
		"joinedAt":       user.JoinedAt.Format(time.RFC3339),
	}
}

// MapInvitedUserFirestoreToFrontend maps Firestore InvitedUser data to frontend format
func MapInvitedUserFirestoreToFrontend(data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"userId":         getStringValue(data, "user_id"),
		"username":       getStringValue(data, "username"),
		"email":          getStringValue(data, "email"),
		"profilePicture": getStringValue(data, "profile_picture"),
		"role":           getStringValue(data, "role"),
		"joinedAt":       getTimeValue(data, "joined_at").Format(time.RFC3339),
	}
}

// MapInvitedUserFirestoreToGo maps Firestore InvitedUser data to Go struct format
func MapInvitedUserFirestoreToGo(data map[string]interface{}) models.InvitedUser {
	return models.InvitedUser{
		UserID:         getStringValue(data, "user_id"),
		Username:       getStringValue(data, "username"),
		Email:          getStringValue(data, "email"),
		ProfilePicture: getStringValue(data, "profile_picture"),
		Role:           getStringValue(data, "role"),
		JoinedAt:       getTimeValue(data, "joined_at"),
	}
}

func mapInvitedUsersArrayToFrontend(data interface{}) []map[string]interface{} {
	var result []map[string]interface{}

	switch users := data.(type) {
	case []interface{}:
		// Step 1: Convert Firestore data to Go structs
		var goUsers []models.InvitedUser
		for _, u := range users {
			if userMap, ok := u.(map[string]interface{}); ok {
				goUser := MapInvitedUserFirestoreToGo(userMap)
				goUsers = append(goUsers, goUser)
			}
		}
		// Step 2: Convert Go structs to frontend format
		for _, goUser := range goUsers {
			result = append(result, MapInvitedUserGoToFrontend(goUser))
		}

	case []map[string]interface{}:
		// Handle Firestore data returned as []map[string]interface{}
		for _, userMap := range users {
			result = append(result, MapInvitedUserFirestoreToFrontend(userMap))
		}

	case []models.InvitedUser:
		// Handle data returned as []models.InvitedUser
		for _, user := range users {
			userMap := map[string]interface{}{
				"userId":         user.UserID,
				"username":       user.Username,
				"email":          user.Email,
				"profilePicture": user.ProfilePicture,
				"role":           user.Role,
				"joinedAt":       user.JoinedAt.Format(time.RFC3339),
			}
			result = append(result, userMap)
		}

	default:
		log.Printf("Unsupported data type: %T\n", data)
	}

	return result
}
