package mappers

import (
	"github.com/rogerjeasy/go-letusconnect/models"
)

// 1. Map Attachment from frontend format to Go struct
func MapAttachmentFrontendToGo(data map[string]interface{}) models.Attachment {
	return models.Attachment{
		FileName:   getStringValue(data, "fileName"),
		URL:        getStringValue(data, "url"),
		UploadedAt: getTimeValue(data, "uploadedAt"),
	}
}

// 2. Map Attachment from Go struct to Firestore format
func MapAttachmentGoToFirestore(attachment models.Attachment) map[string]interface{} {
	return map[string]interface{}{
		"file_name":   attachment.FileName,
		"url":         attachment.URL,
		"uploaded_at": attachment.UploadedAt,
	}
}

// 3. Map Attachment from Firestore format to frontend format
func MapAttachmentFirestoreToFrontend(data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"fileName":   getStringValue(data, "file_name"),
		"url":        getStringValue(data, "url"),
		"uploadedAt": getTimeValue(data, "uploaded_at"),
	}
}

// 4. Map Attachment from Firestore format to Go struct
func MapAttachmentFirestoreToGo(data map[string]interface{}) models.Attachment {
	return models.Attachment{
		FileName:   getStringValue(data, "file_name"),
		URL:        getStringValue(data, "url"),
		UploadedAt: getTimeValue(data, "uploaded_at"),
	}
}
