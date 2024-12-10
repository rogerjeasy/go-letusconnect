package mappers

import (
	"time"

	"github.com/rogerjeasy/go-letusconnect/models"
)

// FrontendToFAQ maps frontend data to the FAQ struct
func FrontendToFAQ(data map[string]interface{}, username, createdBy string) *models.FAQ {
	return &models.FAQ{
		Username:  username,
		Question:  data["question"].(string),
		Response:  data["response"].(string),
		Category:  data["category"].(string),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: createdBy,
		Status:    "active",
	}
}

// FAQToFirestore maps the FAQ struct to Firestore-compatible format
func FAQToFirestore(faq *models.FAQ) map[string]interface{} {
	return map[string]interface{}{
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

// FirestoreToFAQ maps Firestore data to the FAQ struct
func FirestoreToFAQ(data map[string]interface{}) *models.FAQ {
	// Safe type assertions with fallback values to prevent panic
	id, _ := data["id"].(string)
	username, _ := data["username"].(string)
	question, _ := data["question"].(string)
	response, _ := data["response"].(string)
	createdBy, _ := data["created_by"].(string)
	updatedBy, _ := data["updated_by"].(string)
	status, _ := data["status"].(string)
	category, _ := data["category"].(string)

	// Safe handling for time fields
	var createdAt, updatedAt time.Time

	if createdAtRaw, ok := data["created_at"].(time.Time); ok {
		createdAt = createdAtRaw
	} else {
		createdAt = time.Now()
	}

	if updatedAtRaw, ok := data["updated_at"].(time.Time); ok {
		updatedAt = updatedAtRaw
	} else {
		updatedAt = time.Now()
	}

	return &models.FAQ{
		ID:        id,
		Username:  username,
		Question:  question,
		Response:  response,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		CreatedBy: createdBy,
		UpdatedBy: updatedBy,
		Status:    status,
		Category:  category,
	}
}
