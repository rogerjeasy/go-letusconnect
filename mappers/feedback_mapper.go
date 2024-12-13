package mappers

import (
	"github.com/rogerjeasy/go-letusconnect/models"
)

// MapFeedbackFrontendToGo maps frontend Feedback data to Go struct format
func MapFeedbackFrontendToGo(data map[string]interface{}) models.Feedback {
	return models.Feedback{
		UserID:    getStringValue(data, "userId"),
		Rating:    getIntValueSafe(data, "rating"),
		Comment:   getStringValue(data, "comment"),
		CreatedAt: getTimeValue(data, "createdAt"),
	}
}

// MapFeedbackGoToFirestore maps Go struct Feedback data to Firestore format
func MapFeedbackGoToFirestore(feedback models.Feedback) map[string]interface{} {
	return map[string]interface{}{
		"user_id":    feedback.UserID,
		"rating":     feedback.Rating,
		"comment":    feedback.Comment,
		"created_at": feedback.CreatedAt,
	}
}

// MapFeedbackFirestoreToFrontend maps Firestore Feedback data to frontend format
func MapFeedbackFirestoreToFrontend(data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"userId":    getStringValue(data, "user_id"),
		"rating":    getIntValueSafe(data, "rating"),
		"comment":   getStringValue(data, "comment"),
		"createdAt": getTimeValue(data, "created_at"),
	}
}

// 4. Map Feedback from Firestore format to Go struct
func MapFeedbackFirestoreToGo(data map[string]interface{}) models.Feedback {
	return models.Feedback{
		UserID:    getStringValue(data, "user_id"),
		Rating:    getIntValueSafe(data, "rating"),
		Comment:   getStringValue(data, "comment"),
		CreatedAt: getTimeValue(data, "created_at"),
	}
}
