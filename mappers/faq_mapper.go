package mappers

import (
	"github.com/rogerjeasy/go-letusconnect/models"
)

// MapFAQFrontendToGo maps frontend FAQ data to Go struct format
func MapFAQFrontendToGo(data map[string]interface{}) models.FAQ {
	return models.FAQ{
		ID:        getStringValue(data, "id"),
		Username:  getStringValue(data, "username"),
		Question:  getStringValue(data, "question"),
		Response:  getStringValue(data, "response"),
		CreatedAt: getTimeValue(data, "createdAt"),
		UpdatedAt: getTimeValue(data, "updatedAt"),
		CreatedBy: getStringValue(data, "createdBy"),
		UpdatedBy: getStringValue(data, "updatedBy"),
		Status:    getStringValue(data, "status"),
		Category:  getStringValue(data, "category"),
	}
}

// MapFAQGoToFirestore maps Go struct FAQ data to Firestore format
func MapFAQGoToFirestore(faq models.FAQ) map[string]interface{} {
	return map[string]interface{}{
		"id":         faq.ID,
		"username":   faq.Username,
		"question":   faq.Question,
		"response":   faq.Response,
		"created_at": faq.CreatedAt,
		"updated_at": faq.UpdatedAt,
		"created_by": faq.CreatedBy,
		"updated_by": faq.UpdatedBy,
		"status":     faq.Status,
		"category":   faq.Category,
	}
}

// MapFAQFirestoreToGo maps Firestore FAQ data to Go struct format
func MapFAQFirestoreToGo(data map[string]interface{}) models.FAQ {
	return models.FAQ{
		ID:        getStringValue(data, "id"),
		Username:  getStringValue(data, "username"),
		Question:  getStringValue(data, "question"),
		Response:  getStringValue(data, "response"),
		CreatedAt: getFirestoreTimeToGoTime(data["created_at"]),
		UpdatedAt: getFirestoreTimeToGoTime(data["updated_at"]),
		CreatedBy: getStringValue(data, "created_by"),
		UpdatedBy: getStringValue(data, "updated_by"),
		Status:    getStringValue(data, "status"),
		Category:  getStringValue(data, "category"),
	}
}

// MapFAQGoToFrontend maps Go struct FAQ data to frontend format
func MapFAQGoToFrontend(faq models.FAQ) map[string]interface{} {
	return map[string]interface{}{
		"id":        faq.ID,
		"username":  faq.Username,
		"question":  faq.Question,
		"response":  faq.Response,
		"createdAt": faq.CreatedAt,
		"updatedAt": faq.UpdatedAt,
		"createdBy": faq.CreatedBy,
		"updatedBy": faq.UpdatedBy,
		"status":    faq.Status,
		"category":  faq.Category,
	}
}
