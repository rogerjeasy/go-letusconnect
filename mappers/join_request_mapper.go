package mappers

import (
	"github.com/rogerjeasy/go-letusconnect/models"
)

// 1. Map JoinRequest from frontend format to Go struct
func MapJoinRequestFrontendToGo(data map[string]interface{}) models.JoinRequest {
	return models.JoinRequest{
		UserID:      getStringValue(data, "userId"),
		UserName:    getStringValue(data, "userName"),
		Message:     getStringValue(data, "message"),
		RequestedAt: getTimeValue(data, "requestedAt"),
		Status:      getStringValue(data, "status"),
	}
}

// 2. Map JoinRequest from Go struct to Firestore format
func MapJoinRequestGoToFirestore(request models.JoinRequest) map[string]interface{} {
	return map[string]interface{}{
		"user_id":      request.UserID,
		"user_name":    request.UserName,
		"message":      request.Message,
		"requested_at": request.RequestedAt,
		"status":       request.Status,
	}
}

// 3. Map JoinRequest from Firestore format to frontend format
func MapJoinRequestFirestoreToFrontend(data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"userId":      getStringValue(data, "user_id"),
		"userName":    getStringValue(data, "user_name"),
		"message":     getStringValue(data, "message"),
		"requestedAt": getTimeValue(data, "requested_at"),
		"status":      getStringValue(data, "status"),
	}
}

// 4. Map JoinRequest from Firestore format to Go struct
func MapJoinRequestFirestoreToGo(data map[string]interface{}) models.JoinRequest {
	return models.JoinRequest{
		UserID:      getStringValue(data, "user_id"),
		UserName:    getStringValue(data, "user_name"),
		Message:     getStringValue(data, "message"),
		RequestedAt: getTimeValue(data, "requested_at"),
		Status:      getStringValue(data, "status"),
	}
}

func mapJoinRequestsArrayToFrontend(data interface{}) []map[string]interface{} {
	if joinRequests, ok := data.([]interface{}); ok {
		var result []map[string]interface{}
		for _, jr := range joinRequests {
			if jrMap, ok := jr.(map[string]interface{}); ok {
				result = append(result, MapJoinRequestFirestoreToFrontend(jrMap))
			}
		}
		return result
	}
	return []map[string]interface{}{}
}

func mapTasksArrayToFrontend(data interface{}) []map[string]interface{} {
	if tasks, ok := data.([]interface{}); ok {
		var result []map[string]interface{}
		for _, task := range tasks {
			if taskMap, ok := task.(map[string]interface{}); ok {
				result = append(result, MapTaskFirestoreToFrontend(taskMap))
			}
		}
		return result
	}
	return []map[string]interface{}{}
}

func mapAttachmentsArrayToFrontend(data interface{}) []map[string]interface{} {
	if attachments, ok := data.([]interface{}); ok {
		var result []map[string]interface{}
		for _, attachment := range attachments {
			if attachmentMap, ok := attachment.(map[string]interface{}); ok {
				result = append(result, MapAttachmentFirestoreToFrontend(attachmentMap))
			}
		}
		return result
	}
	return []map[string]interface{}{}
}

func mapFeedbacksArrayToFrontend(data interface{}) []map[string]interface{} {
	if feedbacks, ok := data.([]interface{}); ok {
		var result []map[string]interface{}
		for _, feedback := range feedbacks {
			if feedbackMap, ok := feedback.(map[string]interface{}); ok {
				result = append(result, MapFeedbackFirestoreToFrontend(feedbackMap))
			}
		}
		return result
	}
	return []map[string]interface{}{}
}

// GetJoinRequestsArray extracts and maps the "join_requests" field to a slice of models.JoinRequest
func GetJoinRequestsArray(data map[string]interface{}, key string) []models.JoinRequest {
	if joinRequests, ok := data[key].([]interface{}); ok {
		var result []models.JoinRequest
		for _, jr := range joinRequests {
			if jrMap, ok := jr.(map[string]interface{}); ok {
				result = append(result, MapJoinRequestFirestoreToGo(jrMap))
			}
		}
		return result
	}
	return []models.JoinRequest{}
}
