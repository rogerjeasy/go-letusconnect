package mappers

import (
	"log"
	"time"

	"github.com/rogerjeasy/go-letusconnect/models"
)

// ========================= BaseMessage Mappers =========================

// 1. MapBaseMessageFrontendToGo maps frontend BaseMessage data to Go struct format
func MapBaseMessageFrontendToGo(data map[string]interface{}) models.BaseMessage {
	return models.BaseMessage{
		ID:          getStringValue(data, "id"),
		SenderID:    getStringValue(data, "senderId"),
		SenderName:  getStringValue(data, "senderName"),
		Content:     getStringValue(data, "content"),
		CreatedAt:   getStringValue(data, "createdAt"),
		UpdatedAt:   getTimeStringValue(data, "updatedAt"),
		ReadStatus:  getReadStatusMap(data, "readStatus"),
		IsDeleted:   getBoolValue(data, "isDeleted"),
		Attachments: getStringArrayValue(data, "attachments"),
		Reactions:   getReactionsMap(data, "reactions"),
		MessageType: getStringValue(data, "messageType"),
		ReplyToID:   getOptionalStringValue(data, "replyToId"),
		IsPinned:    getBoolValue(data, "isPinned"),
		Priority:    getStringValue(data, "priority"),
	}
}

// 2. MapBaseMessageGoToFirestore maps Go struct BaseMessage data to Firestore format
func MapBaseMessageGoToFirestore(message models.BaseMessage) map[string]interface{} {
	return map[string]interface{}{
		"id":           message.ID,
		"sender_id":    message.SenderID,
		"sender_name":  message.SenderName,
		"content":      message.Content,
		"created_at":   message.CreatedAt,
		"updated_at":   message.UpdatedAt,
		"read_status":  message.ReadStatus,
		"is_deleted":   message.IsDeleted,
		"attachments":  message.Attachments,
		"reactions":    message.Reactions,
		"message_type": message.MessageType,
		"reply_to_id":  message.ReplyToID,
		"is_pinned":    message.IsPinned,
		"priority":     message.Priority,
	}
}

// 3. MapBaseMessageFirestoreToFrontend maps Firestore BaseMessage data to frontend format
func MapBaseMessageFirestoreToFrontend(data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":          getStringValue(data, "id"),
		"senderId":    getStringValue(data, "sender_id"),
		"senderName":  getStringValue(data, "sender_name"),
		"content":     getStringValue(data, "content"),
		"createdAt":   getStringValue(data, "created_at"),
		"updatedAt":   getOptionalStringValue(data, "updated_at"),
		"readStatus":  getReadStatusMap(data, "read_status"),
		"isDeleted":   getBoolValue(data, "is_deleted"),
		"attachments": getStringArrayValue(data, "attachments"),
		"reactions":   getReactionsMap(data, "reactions"),
		"messageType": getStringValue(data, "message_type"),
		"replyToId":   getOptionalStringValue(data, "reply_to_id"),
		"isPinned":    getBoolValue(data, "is_pinned"),
		"priority":    getStringValue(data, "priority"),
	}
}

// MapBaseMessageGoToFrontend maps a BaseMessage Go struct to frontend format.
func MapBaseMessageGoToFrontend(message models.BaseMessage) map[string]interface{} {
	return map[string]interface{}{
		"id":          message.ID,
		"senderId":    message.SenderID,
		"senderName":  message.SenderName,
		"content":     message.Content,
		"createdAt":   message.CreatedAt,
		"updatedAt":   message.UpdatedAt,
		"readStatus":  message.ReadStatus,
		"isDeleted":   message.IsDeleted,
		"attachments": message.Attachments,
		"reactions":   message.Reactions,
		"messageType": message.MessageType,
		"replyToId":   message.ReplyToID,
		"isPinned":    message.IsPinned,
		"priority":    message.Priority,
	}
}

// 4. MapBaseMessageFirestoreToGo maps Firestore BaseMessage data to Go struct format
func MapBaseMessageFirestoreToGo(data map[string]interface{}) models.BaseMessage {
	return models.BaseMessage{
		ID:          getStringValue(data, "id"),
		SenderID:    getStringValue(data, "sender_id"),
		SenderName:  getStringValue(data, "sender_name"),
		Content:     getStringValue(data, "content"),
		CreatedAt:   getStringValue(data, "created_at"),
		UpdatedAt:   getTimeStringValue(data, "updated_at"),
		ReadStatus:  getReadStatusMap(data, "read_status"),
		IsDeleted:   getBoolValue(data, "is_deleted"),
		Attachments: getStringArrayValue(data, "attachments"),
		Reactions:   getReactionsMap(data, "reactions"),
		MessageType: getStringValue(data, "message_type"),
		ReplyToID:   getOptionalStringValue(data, "reply_to_id"),
		IsPinned:    getBoolValue(data, "is_pinned"),
		Priority:    getStringValue(data, "priority"),
	}
}

// ========================= DirectMessage Mappers =========================

// 1. MapDirectMessageFrontendToGo maps frontend DirectMessage data to Go struct format
func MapDirectMessageFrontendToGo(data map[string]interface{}) models.DirectMessage {
	return models.DirectMessage{
		BaseMessage: models.BaseMessage{
			ID:          getStringValue(data, "id"),
			SenderID:    getStringValue(data, "senderId"),
			SenderName:  getStringValue(data, "senderName"),
			Content:     getStringValue(data, "content"),
			CreatedAt:   getStringValue(data, "createdAt"),
			UpdatedAt:   getTimeStringValue(data, "updatedAt"),
			ReadStatus:  getReadStatusMap(data, "readStatus"),
			IsDeleted:   getBoolValue(data, "isDeleted"),
			Attachments: getStringArrayValue(data, "attachments"),
			Reactions:   getReactionsMap(data, "reactions"),
			MessageType: getStringValue(data, "messageType"),
			ReplyToID:   getOptionalStringValue(data, "replyToId"),
			IsPinned:    getBoolValue(data, "isPinned"),
			Priority:    getStringValue(data, "priority"),
		},
		ReceiverID: getStringValue(data, "receiverId"),
	}
}

// 2. MapDirectMessageGoToFirestore maps Go struct DirectMessage data to Firestore format
func MapDirectMessageGoToFirestore(message models.DirectMessage) map[string]interface{} {
	return map[string]interface{}{
		"id":           message.ID,
		"sender_id":    message.SenderID,
		"sender_name":  message.SenderName,
		"content":      message.Content,
		"created_at":   message.CreatedAt,
		"updated_at":   message.UpdatedAt,
		"read_status":  message.ReadStatus,
		"is_deleted":   message.IsDeleted,
		"attachments":  message.Attachments,
		"reactions":    message.Reactions,
		"message_type": message.MessageType,
		"reply_to_id":  message.ReplyToID,
		"is_pinned":    message.IsPinned,
		"priority":     message.Priority,
		"receiver_id":  message.ReceiverID,
	}
}

// 3. MapDirectMessageFirestoreToFrontend maps Firestore DirectMessage data to frontend format
func MapDirectMessageFirestoreToFrontend(data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":          getStringValue(data, "id"),
		"senderId":    getStringValue(data, "sender_id"),
		"senderName":  getStringValue(data, "sender_name"),
		"content":     getStringValue(data, "content"),
		"createdAt":   getStringValue(data, "created_at"),
		"updatedAt":   getOptionalStringValue(data, "updated_at"),
		"readStatus":  getReadStatusMap(data, "read_status"),
		"isDeleted":   getBoolValue(data, "is_deleted"),
		"attachments": getStringArrayValue(data, "attachments"),
		"reactions":   getReactionsMap(data, "reactions"),
		"messageType": getStringValue(data, "message_type"),
		"replyToId":   getOptionalStringValue(data, "reply_to_id"),
		"isPinned":    getBoolValue(data, "is_pinned"),
		"priority":    getStringValue(data, "priority"),
		"receiverId":  getStringValue(data, "receiver_id"),
	}
}

// 4. MapDirectMessageFirestoreToGo maps Firestore DirectMessage data to Go struct format
func MapDirectMessageFirestoreToGo(data map[string]interface{}) models.DirectMessage {
	return models.DirectMessage{
		BaseMessage: models.BaseMessage{
			ID:          getStringValue(data, "id"),
			SenderID:    getStringValue(data, "sender_id"),
			SenderName:  getStringValue(data, "sender_name"),
			Content:     getStringValue(data, "content"),
			CreatedAt:   getStringValue(data, "created_at"),
			UpdatedAt:   getTimeStringValue(data, "updated_at"),
			ReadStatus:  getReadStatusMap(data, "read_status"),
			IsDeleted:   getBoolValue(data, "is_deleted"),
			Attachments: getStringArrayValue(data, "attachments"),
			Reactions:   getReactionsMap(data, "reactions"),
			MessageType: getStringValue(data, "message_type"),
			ReplyToID:   getOptionalStringValue(data, "reply_to_id"),
			IsPinned:    getBoolValue(data, "is_pinned"),
			Priority:    getStringValue(data, "priority"),
		},
		ReceiverID: getStringValue(data, "receiver_id"),
	}
}

// ========================= GroupMessage Mappers =========================

// 1. MapGroupMessageFrontendToGo maps frontend GroupMessage data to Go struct format
func MapGroupMessageFrontendToGo(data map[string]interface{}) models.GroupMessage {
	return models.GroupMessage{
		BaseMessage: models.BaseMessage{
			ID:          getStringValue(data, "id"),
			SenderID:    getStringValue(data, "senderId"),
			SenderName:  getStringValue(data, "senderName"),
			Content:     getStringValue(data, "content"),
			CreatedAt:   getStringValue(data, "createdAt"),
			UpdatedAt:   getTimeStringValue(data, "updatedAt"),
			ReadStatus:  getReadStatusMap(data, "readStatus"),
			IsDeleted:   getBoolValue(data, "isDeleted"),
			Attachments: getStringArrayValue(data, "attachments"),
			Reactions:   getReactionsMap(data, "reactions"),
			MessageType: getStringValue(data, "messageType"),
			ReplyToID:   getOptionalStringValue(data, "replyToId"),
			IsPinned:    getBoolValue(data, "isPinned"),
			Priority:    getStringValue(data, "priority"),
		},
		ProjectID: getStringValue(data, "projectId"),
		GroupID:   getOptionalStringValue(data, "groupId"),
	}
}

// 2. MapGroupMessageGoToFirestore maps Go struct GroupMessage data to Firestore format
func MapGroupMessageGoToFirestore(message models.GroupMessage) map[string]interface{} {
	return map[string]interface{}{
		"id":           message.ID,
		"sender_id":    message.SenderID,
		"sender_name":  message.SenderName,
		"content":      message.Content,
		"created_at":   message.CreatedAt,
		"updated_at":   message.UpdatedAt,
		"read_status":  message.ReadStatus,
		"is_deleted":   message.IsDeleted,
		"attachments":  message.Attachments,
		"reactions":    message.Reactions,
		"message_type": message.MessageType,
		"reply_to_id":  message.ReplyToID,
		"is_pinned":    message.IsPinned,
		"priority":     message.Priority,
		"project_id":   message.ProjectID,
		"group_id":     message.GroupID,
	}
}

// 3. MapGroupMessageFirestoreToFrontend maps Firestore GroupMessage data to frontend format
func MapGroupMessageFirestoreToFrontend(data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":          getStringValue(data, "id"),
		"senderId":    getStringValue(data, "sender_id"),
		"senderName":  getStringValue(data, "sender_name"),
		"content":     getStringValue(data, "content"),
		"createdAt":   getStringValue(data, "created_at"),
		"updatedAt":   getOptionalStringValue(data, "updated_at"),
		"readStatus":  getReadStatusMap(data, "read_status"),
		"isDeleted":   getBoolValue(data, "is_deleted"),
		"attachments": getStringArrayValue(data, "attachments"),
		"reactions":   getReactionsMap(data, "reactions"),
		"messageType": getStringValue(data, "message_type"),
		"replyToId":   getOptionalStringValue(data, "reply_to_id"),
		"isPinned":    getBoolValue(data, "is_pinned"),
		"priority":    getStringValue(data, "priority"),
		"projectId":   getStringValue(data, "project_id"),
		"groupId":     getOptionalStringValue(data, "group_id"),
	}
}

// 4. MapGroupMessageFirestoreToGo maps Firestore GroupMessage data to Go struct format
func MapGroupMessageFirestoreToGo(data map[string]interface{}) models.GroupMessage {
	return models.GroupMessage{
		BaseMessage: models.BaseMessage{
			ID:          getStringValue(data, "id"),
			SenderID:    getStringValue(data, "sender_id"),
			SenderName:  getStringValue(data, "sender_name"),
			Content:     getStringValue(data, "content"),
			CreatedAt:   getStringValue(data, "created_at"),
			UpdatedAt:   getTimeStringValue(data, "updated_at"),
			ReadStatus:  getReadStatusMap(data, "read_status"),
			IsDeleted:   getBoolValue(data, "is_deleted"),
			Attachments: getStringArrayValue(data, "attachments"),
			Reactions:   getReactionsMap(data, "reactions"),
			MessageType: getStringValue(data, "message_type"),
			ReplyToID:   getOptionalStringValue(data, "reply_to_id"),
			IsPinned:    getBoolValue(data, "is_pinned"),
			Priority:    getStringValue(data, "priority"),
		},
		ProjectID: getStringValue(data, "project_id"),
		GroupID:   getOptionalStringValue(data, "group_id"),
	}
}

// MapDirectMessageGoToFrontend maps a DirectMessage struct to frontend format
func MapDirectMessageGoToFrontend(message models.DirectMessage) map[string]interface{} {
	return map[string]interface{}{
		"id":          message.ID,
		"senderId":    message.SenderID,
		"senderName":  message.SenderName,
		"receiverId":  message.ReceiverID,
		"content":     message.Content,
		"createdAt":   message.CreatedAt,
		"updatedAt":   message.UpdatedAt,
		"readStatus":  message.ReadStatus,
		"isDeleted":   message.IsDeleted,
		"attachments": message.Attachments,
		"reactions":   message.Reactions,
		"messageType": message.MessageType,
		"replyToId":   message.ReplyToID,
		"isPinned":    message.IsPinned,
		"priority":    message.Priority,
	}
}

// MapGroupMessageGoToFrontend maps a GroupMessage struct to frontend format
func MapGroupMessageGoToFrontend(message models.GroupMessage) map[string]interface{} {
	return map[string]interface{}{
		"id":          message.ID,
		"senderId":    message.SenderID,
		"senderName":  message.SenderName,
		"content":     message.Content,
		"createdAt":   message.CreatedAt,
		"updatedAt":   message.UpdatedAt,
		"readStatus":  message.ReadStatus,
		"isDeleted":   message.IsDeleted,
		"attachments": message.Attachments,
		"reactions":   message.Reactions,
		"messageType": message.MessageType,
		"replyToId":   message.ReplyToID,
		"isPinned":    message.IsPinned,
		"priority":    message.Priority,
		"projectId":   message.ProjectID,
		"groupId":     message.GroupID,
	}
}

// MapMessagesGoToFrontend maps a Messages struct to frontend format
func MapMessagesGoToFrontend(messages models.Messages) map[string]interface{} {
	directMessages := []map[string]interface{}{}
	for _, msg := range messages.DirectMessages {
		directMessages = append(directMessages, MapDirectMessageGoToFrontend(msg))
	}

	return map[string]interface{}{
		"channelId":      messages.ChannelID,
		"directMessages": directMessages,
	}
}

// ========================= Messages Mappers =========================

// 1. MapMessagesFrontendToGo maps frontend Messages data to Go struct format
func MapMessagesFrontendToGo(data map[string]interface{}) models.Messages {
	directMessages := getDirectMessagesArrayFromFrontend(data, "directMessages")

	return models.Messages{
		ChannelID:      getStringValue(data, "channelId"),
		DirectMessages: directMessages,
	}
}

// 2. MapMessagesGoToFirestore maps Go struct Messages data to Firestore format
func MapMessagesGoToFirestore(messages models.Messages) map[string]interface{} {
	return map[string]interface{}{
		"channel_id":      messages.ChannelID,
		"direct_messages": mapDirectMessagesArrayToFirestore(messages.DirectMessages),
	}
}

// 3. MapMessagesFirestoreToFrontend maps Firestore Messages data to frontend format
func MapMessagesFirestoreToFrontend(data map[string]interface{}) map[string]interface{} {
	directMessages := getDirectMessagesArrayFromFirestore(data, "direct_messages")

	return map[string]interface{}{
		"channelId":      getStringValue(data, "channel_id"),
		"directMessages": directMessages,
	}
}

// 4. MapMessagesFirestoreToGo maps Firestore Messages data to Go struct format
func MapMessagesFirestoreToGo(data map[string]interface{}) models.Messages {
	return models.Messages{
		ChannelID:      getStringValue(data, "channel_id"),
		DirectMessages: getDirectMessagesArrayFromFirestoreToGo(data, "direct_messages"),
	}
}

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
