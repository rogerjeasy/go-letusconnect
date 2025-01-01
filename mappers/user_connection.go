package mappers

import (
	"github.com/rogerjeasy/go-letusconnect/models"
)

// MapConnectionsFrontendToGo maps frontend to Go struct
func MapConnectionsFrontendToGo(data map[string]interface{}) models.UserConnections {
	return models.UserConnections{
		ID:              getStringValue(data, "id"),
		UID:             getStringValue(data, "uid"),
		Connections:     mapConnectionsMapFrontendToGo(data["connections"]),
		PendingRequests: mapRequestsMapFrontendToGo(data["pendingRequests"]),
	}
}

// MapConnectionsGoToFirestore maps Go struct to Firestore
func MapConnectionsGoToFirestore(conn models.UserConnections) map[string]interface{} {
	return map[string]interface{}{
		"id":               conn.ID,
		"uid":              conn.UID,
		"connections":      mapConnectionsMapGoToFirestore(conn.Connections),
		"pending_requests": mapRequestsMapGoToFirestore(conn.PendingRequests),
	}
}

// MapConnectionsFirestoreToFrontend maps Firestore to frontend
func MapConnectionsFirestoreToFrontend(data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":              getStringValue(data, "id"),
		"uid":             getStringValue(data, "uid"),
		"connections":     mapConnectionsMapFirestoreToFrontend(data["connections"]),
		"pendingRequests": mapRequestsMapFirestoreToFrontend(data["pending_requests"]),
	}
}

func MapConnectionRequestFrontendToGo(data map[string]interface{}) models.ConnectionRequest {
	return models.ConnectionRequest{
		ToUID:   getStringValue(data, "toUid"),
		Message: getStringValue(data, "message"),
		Status:  "pending",
	}
}

// MapConnectionsFirestoreToGo maps Firestore to Go struct
func MapConnectionsFirestoreToGo(data map[string]interface{}) models.UserConnections {
	return models.UserConnections{
		ID:              getStringValue(data, "id"),
		UID:             getStringValue(data, "uid"),
		Connections:     mapConnectionsMapFirestoreToGo(data["connections"]),
		PendingRequests: mapRequestsMapFirestoreToGo(data["pending_requests"]),
	}
}

// Helper functions for mapping nested maps
func mapConnectionsMapFrontendToGo(data interface{}) map[string]models.Connection {
	result := make(map[string]models.Connection)
	if conns, ok := data.(map[string]interface{}); ok {
		for k, v := range conns {
			if connMap, ok := v.(map[string]interface{}); ok {
				result[k] = models.Connection{
					TargetUID:  getStringValue(connMap, "targetUid"),
					SentAt:     getTimeValue(connMap, "sentAt"),
					AcceptedAt: getTimeValue(connMap, "acceptedAt"),
					Status:     getStringValue(connMap, "status"),
				}
			}
		}
	}
	return result
}

func mapRequestsMapFrontendToGo(data interface{}) map[string]models.ConnectionRequest {
	result := make(map[string]models.ConnectionRequest)
	if reqs, ok := data.(map[string]interface{}); ok {
		for k, v := range reqs {
			if reqMap, ok := v.(map[string]interface{}); ok {
				result[k] = models.ConnectionRequest{
					FromUID: getStringValue(reqMap, "fromUid"),
					ToUID:   getStringValue(reqMap, "toUid"),
					SentAt:  getTimeValue(reqMap, "sentAt"),
					Message: getStringValue(reqMap, "message"),
					Status:  getStringValue(reqMap, "status"),
				}
			}
		}
	}
	return result
}

func mapConnectionsMapGoToFirestore(conns map[string]models.Connection) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range conns {
		result[k] = map[string]interface{}{
			"target_uid":  v.TargetUID,
			"sent_at":     v.SentAt,
			"accepted_at": v.AcceptedAt,
			"status":      v.Status,
		}
	}
	return result
}

func mapRequestsMapGoToFirestore(reqs map[string]models.ConnectionRequest) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range reqs {
		result[k] = map[string]interface{}{
			"from_uid": v.FromUID,
			"to_uid":   v.ToUID,
			"sent_at":  v.SentAt,
			"message":  v.Message,
			"status":   v.Status,
		}
	}
	return result
}

func mapConnectionsMapFirestoreToFrontend(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	if conns, ok := data.(map[string]interface{}); ok {
		for k, v := range conns {
			if connMap, ok := v.(map[string]interface{}); ok {
				result[k] = map[string]interface{}{
					"targetUid":  getStringValue(connMap, "target_uid"),
					"sentAt":     getTimeValue(connMap, "sent_at"),
					"acceptedAt": getTimeValue(connMap, "accepted_at"),
					"status":     getStringValue(connMap, "status"),
				}
			}
		}
	}
	return result
}

func mapRequestsMapFirestoreToFrontend(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	if reqs, ok := data.(map[string]interface{}); ok {
		for k, v := range reqs {
			if reqMap, ok := v.(map[string]interface{}); ok {
				result[k] = map[string]interface{}{
					"fromUid": getStringValue(reqMap, "from_uid"),
					"toUid":   getStringValue(reqMap, "to_uid"),
					"sentAt":  getTimeValue(reqMap, "sent_at"),
					"message": getStringValue(reqMap, "message"),
					"status":  getStringValue(reqMap, "status"),
				}
			}
		}
	}
	return result
}

func mapConnectionsMapFirestoreToGo(data interface{}) map[string]models.Connection {
	result := make(map[string]models.Connection)
	if conns, ok := data.(map[string]interface{}); ok {
		for k, v := range conns {
			if connMap, ok := v.(map[string]interface{}); ok {
				result[k] = models.Connection{
					TargetUID:  getStringValue(connMap, "target_uid"),
					SentAt:     getTimeValue(connMap, "sent_at"),
					AcceptedAt: getTimeValue(connMap, "accepted_at"),
					Status:     getStringValue(connMap, "status"),
				}
			}
		}
	}
	return result
}

func mapRequestsMapFirestoreToGo(data interface{}) map[string]models.ConnectionRequest {
	result := make(map[string]models.ConnectionRequest)
	if reqs, ok := data.(map[string]interface{}); ok {
		for k, v := range reqs {
			if reqMap, ok := v.(map[string]interface{}); ok {
				result[k] = models.ConnectionRequest{
					FromUID: getStringValue(reqMap, "from_uid"),
					ToUID:   getStringValue(reqMap, "to_uid"),
					SentAt:  getTimeValue(reqMap, "sent_at"),
					Message: getStringValue(reqMap, "message"),
					Status:  getStringValue(reqMap, "status"),
				}
			}
		}
	}
	return result
}
