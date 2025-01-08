package mappers

import (
	"time"

	"github.com/rogerjeasy/go-letusconnect/models"
)

// ========================= MessageConversation Mappers =========================

// 1. Frontend to Go
func MapMessageConversationFrontendToGo(data map[string]interface{}) models.MessageConversation {
	return models.MessageConversation{
		ID:        getStringValue(data, "id"),
		CreatedAt: getTimeValue(data, "createdAt"),
		Message:   getStringValue(data, "message"),
		Response:  getStringValue(data, "response"),
		Role:      getStringValue(data, "role"),
	}
}

// 2. Go to Firestore
func MapMessageConversationGoToFirestore(msg models.MessageConversation) map[string]interface{} {
	return map[string]interface{}{
		"id":         msg.ID,
		"created_at": msg.CreatedAt,
		"message":    msg.Message,
		"response":   msg.Response,
		"role":       msg.Role,
	}
}

// 3. Firestore to Frontend
func MapMessageConversationFirestoreToFrontend(data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":        getStringValue(data, "id"),
		"createdAt": getTimeValue(data, "created_at").Format(time.RFC3339),
		"message":   getStringValue(data, "message"),
		"response":  getStringValue(data, "response"),
		"role":      getStringValue(data, "role"),
	}
}

// 4. Firestore to Go
func MapMessageConversationFirestoreToGo(data map[string]interface{}) models.MessageConversation {
	return models.MessageConversation{
		ID:        getStringValue(data, "id"),
		CreatedAt: getTimeValue(data, "created_at"),
		Message:   getStringValue(data, "message"),
		Response:  getStringValue(data, "response"),
		Role:      getStringValue(data, "role"),
	}
}

// MapConversationGoToFrontend maps Go struct Conversation data to frontend format
func MapConversationGoToFrontend(conv models.Conversation) map[string]interface{} {
	return map[string]interface{}{
		"id":        conv.ID,
		"userId":    conv.UserID,
		"title":     conv.Title,
		"createdAt": conv.CreatedAt.Format(time.RFC3339),
		"updatedAt": conv.UpdatedAt.Format(time.RFC3339),
		"messages":  mapMessageConversationArrayGoToFrontend(conv.Messages),
	}
}

// ========================= Conversation Mappers =========================

// 1. Frontend to Go
func MapConversationFrontendToGo(data map[string]interface{}) models.Conversation {
	return models.Conversation{
		ID:        getStringValue(data, "id"),
		UserID:    getStringValue(data, "userId"),
		Title:     getStringValue(data, "title"),
		CreatedAt: getTimeValue(data, "createdAt"),
		UpdatedAt: getTimeValue(data, "updatedAt"),
		Messages:  getMessageConversationArrayFromFrontend(data, "messages"),
	}
}

// 2. Go to Firestore
func MapConversationGoToFirestore(conv models.Conversation) map[string]interface{} {
	return map[string]interface{}{
		"id":         conv.ID,
		"user_id":    conv.UserID,
		"title":      conv.Title,
		"created_at": conv.CreatedAt,
		"updated_at": conv.UpdatedAt,
		"messages":   mapMessageConversationArrayToFirestore(conv.Messages),
	}
}

// 3. Firestore to Frontend
func MapConversationFirestoreToFrontend(data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":        getStringValue(data, "id"),
		"userId":    getStringValue(data, "user_id"),
		"title":     getStringValue(data, "title"),
		"createdAt": getTimeValue(data, "created_at").Format(time.RFC3339),
		"updatedAt": getTimeValue(data, "updated_at").Format(time.RFC3339),
		"messages":  mapMessageConversationArrayToFrontend(getMessageConversationArrayFromFirestore(data, "messages")),
	}
}

// 4. Firestore to Go
func MapConversationFirestoreToGo(data map[string]interface{}) models.Conversation {
	return models.Conversation{
		ID:        getStringValue(data, "id"),
		UserID:    getStringValue(data, "user_id"),
		Title:     getStringValue(data, "title"),
		CreatedAt: getTimeValue(data, "created_at"),
		UpdatedAt: getTimeValue(data, "updated_at"),
		Messages:  getMessageConversationArrayFromFirestore(data, "messages"),
	}
}

// Helper functions for handling arrays of MessageConversation

func getMessageConversationArrayFromFrontend(data map[string]interface{}, key string) []models.MessageConversation {
	messagesData, ok := data[key].([]interface{})
	if !ok {
		return []models.MessageConversation{}
	}

	messages := make([]models.MessageConversation, len(messagesData))
	for i, msg := range messagesData {
		if msgMap, ok := msg.(map[string]interface{}); ok {
			messages[i] = MapMessageConversationFrontendToGo(msgMap)
		}
	}
	return messages
}

func mapMessageConversationArrayToFirestore(messages []models.MessageConversation) []map[string]interface{} {
	result := make([]map[string]interface{}, len(messages))
	for i, msg := range messages {
		result[i] = MapMessageConversationGoToFirestore(msg)
	}
	return result
}

func getMessageConversationArrayFromFirestore(data map[string]interface{}, key string) []models.MessageConversation {
	messagesData, ok := data[key].([]interface{})
	if !ok {
		return []models.MessageConversation{}
	}

	messages := make([]models.MessageConversation, len(messagesData))
	for i, msg := range messagesData {
		if msgMap, ok := msg.(map[string]interface{}); ok {
			messages[i] = MapMessageConversationFirestoreToGo(msgMap)
		}
	}
	return messages
}

func mapMessageConversationArrayToFrontend(messages []models.MessageConversation) []map[string]interface{} {
	result := make([]map[string]interface{}, len(messages))
	for i, msg := range messages {
		result[i] = MapMessageConversationGoToFirestore(msg)
	}
	return result
}

func mapMessageConversationArrayGoToFrontend(messages []models.MessageConversation) []map[string]interface{} {
	result := make([]map[string]interface{}, len(messages))
	for i, msg := range messages {
		result[i] = mapMessageConversationGoToFrontend(msg)
	}
	return result
}

// Helper function to convert a single MessageConversation to frontend format
func mapMessageConversationGoToFrontend(msg models.MessageConversation) map[string]interface{} {
	return map[string]interface{}{
		"id":        msg.ID,
		"createdAt": msg.CreatedAt.Format(time.RFC3339),
		"message":   msg.Message,
		"response":  msg.Response,
		"role":      msg.Role,
	}
}

// Assuming these helper functions exist in your codebase:
// getStringValue(data map[string]interface{}, key string) string
// getTimeValue(data map[string]interface{}, key string) time.Time
// getBoolValue(data map[string]interface{}, key string) bool
