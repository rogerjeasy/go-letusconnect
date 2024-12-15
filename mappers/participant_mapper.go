package mappers

import (
	"log"
	"time"

	"github.com/rogerjeasy/go-letusconnect/models"
)

// MapParticipantFrontendToGo maps frontend Participant data to Go struct format
func MapParticipantFrontendToGo(data map[string]interface{}) models.Participant {
	return models.Participant{
		UserID:         getStringValue(data, "userId"),
		Role:           getStringValue(data, "role"),
		ProfilePicture: getStringValue(data, "profilePicture"),
		Username:       getStringValue(data, "username"),
		Email:          getStringValue(data, "email"),
		JoinedAt:       getTimeValue(data, "joinedAt"),
	}
}

// MapParticipantGoToFirestore maps Go struct Participant data to Firestore format
func MapParticipantGoToFirestore(user models.Participant) map[string]interface{} {
	return map[string]interface{}{
		"user_id":         user.UserID,
		"role":            user.Role,
		"profile_picture": user.ProfilePicture,
		"username":        user.Username,
		"email":           user.Email,
		"joined_at":       user.JoinedAt,
	}
}

func MapParticipantGoToFrontend(user models.Participant) map[string]interface{} {
	return map[string]interface{}{
		"userId":         user.UserID,
		"role":           user.Role,
		"profilePicture": user.ProfilePicture,
		"username":       user.Username,
		"email":          user.Email,
		"joinedAt":       user.JoinedAt.Format(time.RFC3339),
	}
}

// MapParticipantFirestoreToFrontend maps Firestore Participant data to frontend format
func MapParticipantFirestoreToFrontend(data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"userId":         getStringValue(data, "user_id"),
		"role":           getStringValue(data, "role"),
		"profilePicture": getStringValue(data, "profile_picture"),
		"username":       getStringValue(data, "username"),
		"email":          getStringValue(data, "email"),
		"joinedAt":       getTimeValue(data, "joined_at").Format(time.RFC3339),
	}
}

// MapParticipantFirestoreToGo maps Firestore Participant data to Go struct format
func MapParticipantFirestoreToGo(data map[string]interface{}) models.Participant {
	return models.Participant{
		UserID:         getStringValue(data, "user_id"),
		Role:           getStringValue(data, "role"),
		ProfilePicture: getStringValue(data, "profile_picture"),
		Username:       getStringValue(data, "username"),
		Email:          getStringValue(data, "email"),
		JoinedAt:       getTimeValue(data, "joined_at"),
	}
}

func mapParticipantsArrayToFrontend(data interface{}) []map[string]interface{} {
	var result []map[string]interface{}

	switch users := data.(type) {
	case []interface{}:
		// Step 1: Convert Firestore data to Go structs
		var goUsers []models.Participant
		for _, u := range users {
			if userMap, ok := u.(map[string]interface{}); ok {
				goUser := MapParticipantFirestoreToGo(userMap)
				goUsers = append(goUsers, goUser)
			}
		}
		// Step 2: Convert Go structs to frontend format
		for _, goUser := range goUsers {
			result = append(result, MapParticipantGoToFrontend(goUser))
		}

	case []map[string]interface{}:
		// Handle Firestore data returned as []map[string]interface{}
		for _, userMap := range users {
			result = append(result, MapParticipantFirestoreToFrontend(userMap))
		}

	case []models.Participant:
		// Handle data returned as []models.Participant
		for _, user := range users {
			userMap := map[string]interface{}{
				"userId":         user.UserID,
				"role":           user.Role,
				"profilePicture": user.ProfilePicture,
				"username":       user.Username,
				"email":          user.Email,
				"joinedAt":       user.JoinedAt.Format(time.RFC3339),
			}
			result = append(result, userMap)
		}

	default:
		log.Printf("Unsupported data type: %T\n", data)
	}

	return result
}
