package mappers

import (
	"time"

	"github.com/rogerjeasy/go-letusconnect/models"
)

// Convert frontend Notification data to Go struct format
func MapNotificationFrontendToGo(data map[string]interface{}) models.Notification {
	createdAt, _ := time.Parse(time.RFC3339, data["createdAt"].(string))
	updatedAt, _ := time.Parse(time.RFC3339, data["updatedAt"].(string))

	var expiresAt, sentAt, readAt *time.Time
	if expiresAtStr, ok := data["expiresAt"].(string); ok && expiresAtStr != "" {
		t, _ := time.Parse(time.RFC3339, expiresAtStr)
		expiresAt = &t
	}
	if sentAtStr, ok := data["sentAt"].(string); ok && sentAtStr != "" {
		t, _ := time.Parse(time.RFC3339, sentAtStr)
		sentAt = &t
	}
	if readAtStr, ok := data["readAt"].(string); ok && readAtStr != "" {
		t, _ := time.Parse(time.RFC3339, readAtStr)
		readAt = &t
	}

	return models.Notification{
		ID:              data["id"].(string),
		UserID:          data["userId"].(string),
		ActorID:         getStringValue(data, "actorId"),
		ActorName:       getStringValue(data, "actorName"),
		ActorType:       getStringValue(data, "actorType"),
		Type:            models.NotificationType(data["type"].(string)),
		Title:           data["title"].(string),
		Content:         data["content"].(string),
		Category:        data["category"].(string),
		Priority:        models.NotificationPriority(data["priority"].(string)),
		Status:          models.NotificationStatus(data["status"].(string)),
		RelatedEntities: MapEntityReferencesFrontendToGo(getArrayValue(data, "relatedEntities")),
		Metadata:        getMapValue(data, "metadata"),
		Actions:         MapActionsFrontendToGo(getArrayValue(data, "actions")),
		ReadStatus:      getReadStatusMap(data, "readStatus"),
		IsArchived:      getReadStatusMap(data, "isArchived"),
		IsImportant:     data["isImportant"].(bool),
		ExpiresAt:       expiresAt,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
		SentAt:          sentAt,
		ReadAt:          readAt,
		Source:          getStringValue(data, "source"),
		Tags:            getStringArrayValue(data, "tags"),
		GroupID:         getStringValue(data, "groupId"),
		DeliveryChannel: getStringValue(data, "deliveryChannel"),
		TargetedUsers:   getStringArrayValue(data, "targetedUsers"),
		IsRead:          data["isRead"].(bool),
	}
}

// MapNotificationGoToFrontend converts a Go Notification struct to frontend format
func MapNotificationGoToFrontend(notification models.Notification) map[string]interface{} {
	result := map[string]interface{}{
		"id":              notification.ID,
		"userId":          notification.UserID,
		"actorId":         notification.ActorID,
		"actorName":       notification.ActorName,
		"actorType":       notification.ActorType,
		"type":            string(notification.Type),
		"title":           notification.Title,
		"content":         notification.Content,
		"category":        notification.Category,
		"priority":        string(notification.Priority),
		"status":          string(notification.Status),
		"relatedEntities": MapEntityReferencesGoToFrontend(notification.RelatedEntities),
		"metadata":        notification.Metadata,
		"actions":         MapActionsGoToFrontend(notification.Actions),
		"readStatus":      notification.ReadStatus,
		"isArchived":      notification.IsArchived,
		"isImportant":     notification.IsImportant,
		"createdAt":       notification.CreatedAt.Format(time.RFC3339),
		"updatedAt":       notification.UpdatedAt.Format(time.RFC3339),
		"source":          notification.Source,
		"tags":            notification.Tags,
		"groupId":         notification.GroupID,
		"deliveryChannel": notification.DeliveryChannel,
		"targetedUsers":   notification.TargetedUsers,
		"isRead":          notification.IsRead,
	}

	if notification.ExpiresAt != nil {
		result["expiresAt"] = notification.ExpiresAt.Format(time.RFC3339)
	}
	if notification.SentAt != nil {
		result["sentAt"] = notification.SentAt.Format(time.RFC3339)
	}
	if notification.ReadAt != nil {
		result["readAt"] = notification.ReadAt.Format(time.RFC3339)
	}

	return result
}

// Convert Go struct Notification data to Firestore format
func MapNotificationGoToFirestore(notification models.Notification) map[string]interface{} {
	data := map[string]interface{}{
		"id":               notification.ID,
		"user_id":          notification.UserID,
		"actor_id":         notification.ActorID,
		"actor_name":       notification.ActorName,
		"actor_type":       notification.ActorType,
		"type":             string(notification.Type),
		"title":            notification.Title,
		"content":          notification.Content,
		"category":         notification.Category,
		"priority":         string(notification.Priority),
		"status":           string(notification.Status),
		"related_entities": MapEntityReferencesGoToFirestore(notification.RelatedEntities),
		"metadata":         notification.Metadata,
		"actions":          MapActionsGoToFirestore(notification.Actions),
		"read_status":      notification.ReadStatus,
		"is_archived":      notification.IsArchived,
		"is_important":     notification.IsImportant,
		"created_at":       notification.CreatedAt,
		"updated_at":       notification.UpdatedAt,
		"source":           notification.Source,
		"tags":             notification.Tags,
		"group_id":         notification.GroupID,
		"delivery_channel": notification.DeliveryChannel,
		"targeted_users":   notification.TargetedUsers,
		"is_read":          notification.IsRead,
	}

	if notification.ExpiresAt != nil {
		data["expires_at"] = *notification.ExpiresAt
	}
	if notification.SentAt != nil {
		data["sent_at"] = *notification.SentAt
	}
	if notification.ReadAt != nil {
		data["read_at"] = *notification.ReadAt
	}

	return data
}

// Convert Firestore Notification data to frontend format
func MapNotificationFirestoreToFrontend(data map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{
		"id":              data["id"].(string),
		"userId":          data["user_id"].(string),
		"actorId":         getStringValue(data, "actor_id"),
		"actorName":       getStringValue(data, "actor_name"),
		"actorType":       getStringValue(data, "actor_type"),
		"type":            data["type"].(string),
		"title":           data["title"].(string),
		"content":         data["content"].(string),
		"category":        data["category"].(string),
		"priority":        data["priority"].(string),
		"status":          data["status"].(string),
		"relatedEntities": MapEntityReferencesFirestoreToFrontend(getArrayValue(data, "related_entities")),
		"metadata":        getMapValue(data, "metadata"),
		"actions":         MapActionsFirestoreToFrontend(getArrayValue(data, "actions")),
		"isRead":          data["is_read"].(bool),
		"isArchived":      data["is_archived"].(bool),
		"isImportant":     data["is_important"].(bool),
		"createdAt":       data["created_at"].(time.Time).Format(time.RFC3339),
		"updatedAt":       data["updated_at"].(time.Time).Format(time.RFC3339),
		"source":          getStringValue(data, "source"),
		"tags":            getStringArrayValue(data, "tags"),
		"groupId":         getStringValue(data, "group_id"),
		"deliveryChannel": getStringValue(data, "delivery_channel"),
		"targetedUsers":   getStringArrayValue(data, "targeted_users"),
	}

	if expiresAt, ok := data["expires_at"].(time.Time); ok {
		result["expiresAt"] = expiresAt.Format(time.RFC3339)
	}
	if sentAt, ok := data["sent_at"].(time.Time); ok {
		result["sentAt"] = sentAt.Format(time.RFC3339)
	}
	if readAt, ok := data["read_at"].(time.Time); ok {
		result["readAt"] = readAt.Format(time.RFC3339)
	}

	return result
}

// Convert Firestore Notification data to Go struct format
func MapNotificationFirestoreToGo(data map[string]interface{}) models.Notification {
	notification := models.Notification{
		ID:              data["id"].(string),
		UserID:          data["user_id"].(string),
		ActorID:         getStringValue(data, "actor_id"),
		ActorName:       getStringValue(data, "actor_name"),
		ActorType:       getStringValue(data, "actor_type"),
		Type:            models.NotificationType(data["type"].(string)),
		Title:           data["title"].(string),
		Content:         data["content"].(string),
		Category:        data["category"].(string),
		Priority:        models.NotificationPriority(data["priority"].(string)),
		Status:          models.NotificationStatus(data["status"].(string)),
		RelatedEntities: MapEntityReferencesFirestoreToGo(getArrayValue(data, "related_entities")),
		Metadata:        getMapValue(data, "metadata"),
		Actions:         MapActionsFirestoreToGo(getArrayValue(data, "actions")),
		ReadStatus:      getReadStatusMap(data, "read_status"),
		IsArchived:      getReadStatusMap(data, "is_archived"),
		IsImportant:     data["is_important"].(bool),
		CreatedAt:       data["created_at"].(time.Time),
		UpdatedAt:       data["updated_at"].(time.Time),
		Source:          getStringValue(data, "source"),
		Tags:            getStringArrayValue(data, "tags"),
		GroupID:         getStringValue(data, "group_id"),
		DeliveryChannel: getStringValue(data, "delivery_channel"),
		TargetedUsers:   getStringArrayValue(data, "targeted_users"),
		IsRead:          data["is_read"].(bool),
	}

	if expiresAt, ok := data["expires_at"].(time.Time); ok {
		notification.ExpiresAt = &expiresAt
	}
	if sentAt, ok := data["sent_at"].(time.Time); ok {
		notification.SentAt = &sentAt
	}
	if readAt, ok := data["read_at"].(time.Time); ok {
		notification.ReadAt = &readAt
	}

	return notification
}

// MapEntityReferencesGoToFrontend converts a slice of EntityReference structs to frontend format
func MapEntityReferencesGoToFrontend(entities []models.EntityReference) []map[string]interface{} {
	result := make([]map[string]interface{}, len(entities))
	for i, entity := range entities {
		result[i] = map[string]interface{}{
			"id":   entity.ID,
			"type": entity.Type,
		}
	}
	return result
}

// MapActionsGoToFrontend converts a slice of NotificationAction structs to frontend format
func MapActionsGoToFrontend(actions []models.NotificationAction) []map[string]interface{} {
	result := make([]map[string]interface{}, len(actions))
	for i, action := range actions {
		result[i] = map[string]interface{}{
			"label":      action.Label,
			"url":        action.URL,
			"isPrimary":  action.IsPrimary,
			"actionType": action.ActionType,
		}
	}
	return result
}
