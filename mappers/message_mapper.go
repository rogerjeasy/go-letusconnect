package mappers

import (
	"log"
	"time"

	"github.com/rogerjeasy/go-letusconnect/models"
)

// MapMessageFrontendToGo converts client JSON data to a Message struct
func MapMessageFrontendToGo(data map[string]interface{}) models.Message {
	return models.Message{
		ID:         getStringValue(data, "id"),
		SenderID:   getStringValue(data, "senderId"),
		ReceiverID: getStringValue(data, "receiverId"),
		Content:    getStringValue(data, "content"),
		CreatedAt:  getTimeValue(data, "createdAt"),
	}
}

// MapMessageGoToFirestore maps a Message struct to Firestore format
func MapMessageGoToFirestore(msg models.Message) map[string]interface{} {
	return map[string]interface{}{
		"id":          msg.ID,
		"sender_id":   msg.SenderID,
		"receiver_id": msg.ReceiverID,
		"content":     msg.Content,
		"created_at":  msg.CreatedAt,
	}
}

// MapMessageFirestoreToGo converts Firestore data to a Message struct
func MapMessageFirestoreToGo(data map[string]interface{}) models.Message {
	return models.Message{
		ID:         getStringValue(data, "id"),
		SenderID:   getStringValue(data, "sender_id"),
		ReceiverID: getStringValue(data, "receiver_id"),
		Content:    getStringValue(data, "content"),
		CreatedAt:  getTimeValue(data, "created_at"),
	}
}

// MapMessageGoToFrontend maps a Message struct to frontend format
func MapMessageGoToFrontend(msg models.Message) map[string]interface{} {
	return map[string]interface{}{
		"id":         msg.ID,
		"senderId":   msg.SenderID,
		"receiverId": msg.ReceiverID,
		"content":    msg.Content,
		"createdAt":  msg.CreatedAt.Format(time.RFC3339),
	}
}

// MapMessageFirestoreToFrontend maps Firestore Message data to frontend format
func MapMessageFirestoreToFrontend(data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":         getStringValue(data, "id"),
		"senderId":   getStringValue(data, "sender_id"),
		"receiverId": getStringValue(data, "receiver_id"),
		"content":    getStringValue(data, "content"),
		"createdAt":  getTimeValue(data, "created_at").Format(time.RFC3339),
	}
}

// MapMessagesArrayToFrontend maps an array of messages to frontend format
func MapMessagesArrayToFrontend(data interface{}) []map[string]interface{} {
	var result []map[string]interface{}

	switch messages := data.(type) {
	case []interface{}:
		// Convert Firestore data to Go structs
		var goMessages []models.Message
		for _, m := range messages {
			if msgMap, ok := m.(map[string]interface{}); ok {
				goMessage := MapMessageFirestoreToGo(msgMap)
				goMessages = append(goMessages, goMessage)
			}
		}
		// Convert Go structs to frontend format
		for _, goMessage := range goMessages {
			result = append(result, MapMessageGoToFrontend(goMessage))
		}

	case []map[string]interface{}:
		// Handle Firestore data returned as []map[string]interface{}
		for _, msgMap := range messages {
			result = append(result, MapMessageFirestoreToFrontend(msgMap))
		}

	case []models.Message:
		// Handle data returned as []models.Message
		for _, msg := range messages {
			msgMap := map[string]interface{}{
				"id":         msg.ID,
				"senderId":   msg.SenderID,
				"receiverId": msg.ReceiverID,
				"content":    msg.Content,
				"createdAt":  msg.CreatedAt.Format(time.RFC3339),
			}
			result = append(result, msgMap)
		}

	default:
		log.Printf("Unsupported data type: %T\n", data)
	}

	return result
}
