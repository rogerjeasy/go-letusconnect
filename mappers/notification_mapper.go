package mappers

import (
	"time"

	"github.com/rogerjeasy/go-letusconnect/models"
)

// Convert frontend Notification data to Go struct format
func MapNotificationFrontendToGo(data map[string]interface{}) models.Notification {
	createdAt, _ := time.Parse(time.RFC3339, data["createdAt"].(string))
	updatedAt, _ := time.Parse(time.RFC3339, data["updatedAt"].(string))

	return models.Notification{
		ID:              data["id"].(string),
		UserID:          data["userId"].(string),
		ActorID:         getStringValue(data, "actorId"),
		ActorName:       getStringValue(data, "actorName"),
		Type:            models.NotificationType(data["type"].(string)),
		Title:           data["title"].(string),
		Content:         data["content"].(string),
		Category:        data["category"].(string),
		Priority:        data["priority"].(string),
		Status:          models.NotificationStatus(data["status"].(string)),
		RelatedEntityID: getStringValue(data, "relatedEntityId"),
		RelatedEntity:   getMapValue(data, "relatedEntity"),
		Metadata:        getStringMapValue(data, "metadata"),
		Actions:         MapActionsFrontendToGo(getArrayValue(data, "actions")),
		IsRead:          data["isRead"].(bool),
		IsArchived:      data["isArchived"].(bool),
		IsImportant:     data["isImportant"].(bool),
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
	}
}

// Convert Go struct Notification data to Firestore format
func MapNotificationGoToFirestore(notification models.Notification) map[string]interface{} {
	return map[string]interface{}{
		"id":                notification.ID,
		"user_id":           notification.UserID,
		"actor_id":          notification.ActorID,
		"actor_name":        notification.ActorName,
		"type":              string(notification.Type),
		"title":             notification.Title,
		"content":           notification.Content,
		"category":          notification.Category,
		"priority":          notification.Priority,
		"status":            string(notification.Status),
		"related_entity_id": notification.RelatedEntityID,
		"related_entity":    notification.RelatedEntity,
		"metadata":          notification.Metadata,
		"actions":           MapActionsGoToFirestore(notification.Actions),
		"is_read":           notification.IsRead,
		"is_archived":       notification.IsArchived,
		"is_important":      notification.IsImportant,
		"created_at":        notification.CreatedAt,
		"updated_at":        notification.UpdatedAt,
	}
}

// Convert Firestore Notification data to frontend format
func MapNotificationFirestoreToFrontend(data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":              data["id"].(string),
		"userId":          data["user_id"].(string),
		"actorId":         getStringValue(data, "actor_id"),
		"actorName":       getStringValue(data, "actor_name"),
		"type":            data["type"].(string),
		"title":           data["title"].(string),
		"content":         data["content"].(string),
		"category":        data["category"].(string),
		"priority":        data["priority"].(string),
		"status":          data["status"].(string),
		"relatedEntityId": getStringValue(data, "related_entity_id"),
		"relatedEntity":   getMapValue(data, "related_entity"),
		"metadata":        getStringMapValue(data, "metadata"),
		"actions":         MapActionsFirestoreToFrontend(getArrayValue(data, "actions")),
		"isRead":          data["is_read"].(bool),
		"isArchived":      data["is_archived"].(bool),
		"isImportant":     data["is_important"].(bool),
		"createdAt":       data["created_at"].(time.Time).Format(time.RFC3339),
		"updatedAt":       data["updated_at"].(time.Time).Format(time.RFC3339),
	}
}

// Convert Firestore Notification data to Go struct format
func MapNotificationFirestoreToGo(data map[string]interface{}) models.Notification {
	return models.Notification{
		ID:              data["id"].(string),
		UserID:          data["user_id"].(string),
		ActorID:         getStringValue(data, "actor_id"),
		ActorName:       getStringValue(data, "actor_name"),
		Type:            models.NotificationType(data["type"].(string)),
		Title:           data["title"].(string),
		Content:         data["content"].(string),
		Category:        data["category"].(string),
		Priority:        data["priority"].(string),
		Status:          models.NotificationStatus(data["status"].(string)),
		RelatedEntityID: getStringValue(data, "related_entity_id"),
		RelatedEntity:   getMapValue(data, "related_entity"),
		Metadata:        getStringMapValue(data, "metadata"),
		Actions:         MapActionsFirestoreToGo(getArrayValue(data, "actions")),
		IsRead:          data["is_read"].(bool),
		IsArchived:      data["is_archived"].(bool),
		IsImportant:     data["is_important"].(bool),
		CreatedAt:       data["created_at"].(time.Time),
		UpdatedAt:       data["updated_at"].(time.Time),
	}
}
