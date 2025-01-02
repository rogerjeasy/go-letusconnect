package mappers

import (
	"time"

	"github.com/rogerjeasy/go-letusconnect/models"
)

func MapConnectionsFrontendToGo(data map[string]interface{}) models.UserConnections {
	return models.UserConnections{
		ID:              getStringValue(data, "id"),
		UID:             getStringValue(data, "uid"),
		Connections:     mapConnectionsMapFrontendToGo(data["connections"]),
		PendingRequests: mapRequestsMapFrontendToGo(data["pendingRequests"]),
		SentRequests:    mapSentRequestsMapFrontendToGo(data["sentRequests"]),
	}
}

func MapConnectionsGoToFirestore(conn models.UserConnections) map[string]interface{} {
	return map[string]interface{}{
		"id":               conn.ID,
		"uid":              conn.UID,
		"connections":      mapConnectionsMapGoToFirestore(conn.Connections),
		"pending_requests": mapRequestsMapGoToFirestore(conn.PendingRequests),
		"sent_requests":    mapSentRequestsMapGoToFirestore(conn.SentRequests),
	}
}

func MapConnectionsFirestoreToFrontend(data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":              getStringValue(data, "id"),
		"uid":             getStringValue(data, "uid"),
		"connections":     mapConnectionsMapFirestoreToFrontend(data["connections"]),
		"pendingRequests": mapRequestsMapFirestoreToFrontend(data["pending_requests"]),
		"sentRequests":    mapSentRequestsMapFirestoreToFrontend(data["sent_requests"]),
	}
}

func MapConnectionsFirestoreToGo(data map[string]interface{}) models.UserConnections {
	return models.UserConnections{
		ID:              getStringValue(data, "id"),
		UID:             getStringValue(data, "uid"),
		Connections:     mapConnectionsMapFirestoreToGo(data["connections"]),
		PendingRequests: mapRequestsMapFirestoreToGo(data["pending_requests"]),
		SentRequests:    mapSentRequestsMapFirestoreToGo(data["sent_requests"]),
	}
}

func MapConnectionRequestFrontendToGo(data map[string]interface{}) models.ConnectionRequest {
	return models.ConnectionRequest{
		ToUID:   getStringValue(data, "toUid"),
		Message: getStringValue(data, "message"),
		Status:  "pending",
	}
}

func MapConnectionsGoToFrontend(conn models.UserConnections) map[string]interface{} {
	return map[string]interface{}{
		"id":              conn.ID,
		"uid":             conn.UID,
		"connections":     mapConnectionsMapGoToFrontend(conn.Connections),
		"pendingRequests": mapRequestsMapGoToFrontend(conn.PendingRequests),
		"sentRequests":    mapSentRequestsMapGoToFrontend(conn.SentRequests),
	}
}

func mapConnectionsMapGoToFrontend(conns map[string]models.Connection) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range conns {
		result[k] = map[string]interface{}{
			"targetUid":  v.TargetUID,
			"targetName": v.TargetName,
			"sentAt":     v.SentAt,
			"acceptedAt": v.AcceptedAt,
			"status":     v.Status,
		}
	}
	return result
}

func mapRequestsMapGoToFrontend(reqs map[string]models.ConnectionRequest) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range reqs {
		result[k] = map[string]interface{}{
			"fromUid":  v.FromUID,
			"fromName": v.FromName,
			"toUid":    v.ToUID,
			"sentAt":   v.SentAt,
			"message":  v.Message,
			"status":   v.Status,
		}
	}
	return result
}

func mapSentRequestsMapGoToFrontend(reqs map[string]models.SentRequest) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range reqs {
		result[k] = map[string]interface{}{
			"toUid":    v.ToUID,
			"sentAt":   v.SentAt,
			"message":  v.Message,
			"status":   v.Status,
			"accepted": v.Accepted,
		}
	}
	return result
}

func mapConnectionsMapFrontendToGo(data interface{}) map[string]models.Connection {
	result := make(map[string]models.Connection)
	if conns, ok := data.(map[string]interface{}); ok {
		for k, v := range conns {
			if connMap, ok := v.(map[string]interface{}); ok {
				result[k] = models.Connection{
					TargetUID:  getStringValue(connMap, "targetUid"),
					TargetName: getStringValue(connMap, "targetName"),
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
					FromUID:  getStringValue(reqMap, "fromUid"),
					FromName: getStringValue(reqMap, "fromName"),
					ToUID:    getStringValue(reqMap, "toUid"),
					SentAt:   getTimeValue(reqMap, "sentAt"),
					Message:  getStringValue(reqMap, "message"),
					Status:   getStringValue(reqMap, "status"),
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
			"target_name": v.TargetName,
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
			"from_uid":  v.FromUID,
			"from_name": v.FromName,
			"to_uid":    v.ToUID,
			"sent_at":   v.SentAt,
			"message":   v.Message,
			"status":    v.Status,
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
					"targetName": getStringValue(connMap, "target_name"),
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
					"fromUid":  getStringValue(reqMap, "from_uid"),
					"fromName": getStringValue(reqMap, "from_name"),
					"toUid":    getStringValue(reqMap, "to_uid"),
					"sentAt":   getTimeValue(reqMap, "sent_at"),
					"message":  getStringValue(reqMap, "message"),
					"status":   getStringValue(reqMap, "status"),
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
					TargetName: getStringValue(connMap, "target_name"),
					SentAt:     connMap["sent_at"].(time.Time),
					AcceptedAt: connMap["accepted_at"].(time.Time),
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
					FromUID:  getStringValue(reqMap, "from_uid"),
					FromName: getStringValue(reqMap, "from_name"),
					ToUID:    getStringValue(reqMap, "to_uid"),
					SentAt:   reqMap["sent_at"].(time.Time),
					Message:  getStringValue(reqMap, "message"),
					Status:   getStringValue(reqMap, "status"),
				}
			}
		}
	}
	return result
}

func mapSentRequestsMapFrontendToGo(data interface{}) map[string]models.SentRequest {
	result := make(map[string]models.SentRequest)
	if reqs, ok := data.(map[string]interface{}); ok {
		for k, v := range reqs {
			if reqMap, ok := v.(map[string]interface{}); ok {
				result[k] = models.SentRequest{
					ToUID:    getStringValue(reqMap, "toUid"),
					SentAt:   getTimeValue(reqMap, "sentAt"),
					Message:  getStringValue(reqMap, "message"),
					Status:   getStringValue(reqMap, "status"),
					Accepted: getTimeValue(reqMap, "accepted"),
				}
			}
		}
	}
	return result
}

func mapSentRequestsMapGoToFirestore(reqs map[string]models.SentRequest) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range reqs {
		result[k] = map[string]interface{}{
			"to_uid":   v.ToUID,
			"sent_at":  v.SentAt,
			"message":  v.Message,
			"status":   v.Status,
			"accepted": v.Accepted,
		}
	}
	return result
}

func mapSentRequestsMapFirestoreToFrontend(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	if reqs, ok := data.(map[string]interface{}); ok {
		for k, v := range reqs {
			if reqMap, ok := v.(map[string]interface{}); ok {
				result[k] = map[string]interface{}{
					"toUid":    getStringValue(reqMap, "to_uid"),
					"sentAt":   getTimeValue(reqMap, "sent_at"),
					"message":  getStringValue(reqMap, "message"),
					"status":   getStringValue(reqMap, "status"),
					"accepted": getTimeValue(reqMap, "accepted"),
				}
			}
		}
	}
	return result
}

func mapSentRequestsMapFirestoreToGo(data interface{}) map[string]models.SentRequest {
	result := make(map[string]models.SentRequest)
	if reqs, ok := data.(map[string]interface{}); ok {
		for k, v := range reqs {
			if reqMap, ok := v.(map[string]interface{}); ok {
				result[k] = models.SentRequest{
					ToUID:    getStringValue(reqMap, "to_uid"),
					SentAt:   getTimeValue(reqMap, "sent_at"),
					Message:  getStringValue(reqMap, "message"),
					Status:   getStringValue(reqMap, "status"),
					Accepted: getTimeValue(reqMap, "accepted"),
				}
			}
		}
	}
	return result
}
