package mappers

import (
	"time"

	"github.com/rogerjeasy/go-letusconnect/models"
)

// FrontendToContactUs maps frontend JSON input to the ContactUs struct
func FrontendToContactUs(data map[string]interface{}) *models.ContactUs {
	return &models.ContactUs{
		Name:        data["name"].(string),
		Email:       data["email"].(string),
		Subject:     data["subject"].(string),
		Message:     data["message"].(string),
		Attachments: toStringSlice(data["attachments"]),
		Status:      models.StatusUnread,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// FirestoreToContactUs maps Firestore data to the ContactUs struct
func FirestoreToContactUs(doc map[string]interface{}) *models.ContactUs {
	return &models.ContactUs{
		ID:           doc["id"].(string),
		Name:         doc["name"].(string),
		Email:        doc["email"].(string),
		Subject:      doc["subject"].(string),
		Message:      doc["message"].(string),
		Attachments:  toStringSlice(doc["attachments"]),
		Status:       doc["status"].(string),
		CreatedAt:    doc["created_at"].(time.Time),
		UpdatedAt:    doc["updated_at"].(time.Time),
		RepliedBy:    doc["replied_by"].(string),
		ReplyMessage: doc["reply_message"].(string),
	}
}

// ContactUsToFirestore maps the ContactUs struct to a Firestore-compatible map
func ContactUsToFirestore(contact *models.ContactUs) map[string]interface{} {
	return map[string]interface{}{
		"id":            contact.ID,
		"name":          contact.Name,
		"email":         contact.Email,
		"subject":       contact.Subject,
		"message":       contact.Message,
		"attachments":   contact.Attachments,
		"status":        contact.Status,
		"created_at":    contact.CreatedAt,
		"updated_at":    contact.UpdatedAt,
		"replied_by":    contact.RepliedBy,
		"reply_message": contact.ReplyMessage,
	}
}

// Helper function to convert interface{} to []string
func toStringSlice(data interface{}) []string {
	if data == nil {
		return []string{}
	}
	rawSlice := data.([]interface{})
	result := make([]string, len(rawSlice))
	for i, v := range rawSlice {
		result[i] = v.(string)
	}
	return result
}
