package mappers

import (
	"time"

	"github.com/rogerjeasy/go-letusconnect/models"
)

// ========================= GroupSettings Mappers =========================

// 1. MapGroupSettingsFrontendToGo maps frontend GroupSettings data to Go struct format
func MapGroupSettingsFrontendToGo(data map[string]interface{}) models.GroupSettings {
	return models.GroupSettings{
		AllowFileSharing:  getBoolValue(data, "allowFileSharing"),
		AllowPinning:      getBoolValue(data, "allowPinning"),
		AllowReactions:    getBoolValue(data, "allowReactions"),
		AllowReplies:      getBoolValue(data, "allowReplies"),
		MuteNotifications: getBoolValue(data, "muteNotifications"),
		OnlyAdminsCanPost: getBoolValue(data, "onlyAdminsCanPost"),
	}
}

// 2. MapGroupSettingsGoToFirestore maps Go struct GroupSettings data to Firestore format
func MapGroupSettingsGoToFirestore(settings models.GroupSettings) map[string]interface{} {
	return map[string]interface{}{
		"allow_file_sharing":   settings.AllowFileSharing,
		"allow_pinning":        settings.AllowPinning,
		"allow_reactions":      settings.AllowReactions,
		"allow_replies":        settings.AllowReplies,
		"mute_notifications":   settings.MuteNotifications,
		"only_admins_can_post": settings.OnlyAdminsCanPost,
	}
}

// 3. MapGroupSettingsFirestoreToFrontend maps Firestore GroupSettings data to frontend format
func MapGroupSettingsFirestoreToFrontend(data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"allowFileSharing":  getBoolValue(data, "allow_file_sharing"),
		"allowPinning":      getBoolValue(data, "allow_pinning"),
		"allowReactions":    getBoolValue(data, "allow_reactions"),
		"allowReplies":      getBoolValue(data, "allow_replies"),
		"muteNotifications": getBoolValue(data, "mute_notifications"),
		"onlyAdminsCanPost": getBoolValue(data, "only_admins_can_post"),
	}
}

// 4. MapGroupSettingsFirestoreToGo maps Firestore GroupSettings data to Go struct format
func MapGroupSettingsFirestoreToGo(data map[string]interface{}) models.GroupSettings {
	return models.GroupSettings{
		AllowFileSharing:  getBoolValue(data, "allow_file_sharing"),
		AllowPinning:      getBoolValue(data, "allow_pinning"),
		AllowReactions:    getBoolValue(data, "allow_reactions"),
		AllowReplies:      getBoolValue(data, "allow_replies"),
		MuteNotifications: getBoolValue(data, "mute_notifications"),
		OnlyAdminsCanPost: getBoolValue(data, "only_admins_can_post"),
	}
}

// ========================= GroupChat Mappers =========================

// 1. MapGroupChatFrontendToGo maps frontend GroupChat data to Go struct format
func MapGroupChatFrontendToGo(data map[string]interface{}) models.GroupChat {
	return models.GroupChat{
		ID:             getStringValue(data, "id"),
		ProjectID:      getStringValue(data, "projectId"),
		CreatedByUID:   getStringValue(data, "createdByUid"),
		CreatedByName:  getStringValue(data, "createdByName"),
		Name:           getStringValue(data, "name"),
		Description:    dereferenceString(getOptionalStringValue(data, "description"), ""),
		Participants:   GetParticipantsArray(data, "participants"),
		Messages:       getBaseMessagesArrayFromFrontend(data, "messages"),
		PinnedMessages: getStringArrayValue(data, "pinnedMessages"),
		IsArchived:     getBoolValue(data, "isArchived"),
		Notifications:  getNotificationsMap(data, "notifications"),
		CreatedAt:      getTimeValue(data, "createdAt"),
		UpdatedAt:      getTimeValue(data, "updatedAt"),
		ReadStatus:     getReadStatusMap(data, "readStatus"),
		GroupSettings:  MapGroupSettingsFrontendToGo(getMapValue(data, "groupSettings")),
	}
}

// 2. MapGroupChatGoToFirestore maps Go struct GroupChat data to Firestore format
func MapGroupChatGoToFirestore(chat models.GroupChat) map[string]interface{} {
	return map[string]interface{}{
		"id":              chat.ID,
		"project_id":      chat.ProjectID,
		"created_by_uid":  chat.CreatedByUID,
		"created_by_name": chat.CreatedByName,
		"name":            chat.Name,
		"description":     chat.Description,
		"participants":    MapParticipantsArrayToFirestore(chat.Participants),
		"messages":        mapBaseMessagesArrayToFirestore(chat.Messages),
		"pinned_messages": chat.PinnedMessages,
		"is_archived":     chat.IsArchived,
		"notifications":   chat.Notifications,
		"created_at":      chat.CreatedAt,
		"updated_at":      chat.UpdatedAt,
		"read_status":     chat.ReadStatus,
		"group_settings":  MapGroupSettingsGoToFirestore(chat.GroupSettings),
	}
}

// 3. MapGroupChatFirestoreToFrontend maps Firestore GroupChat data to frontend format
func MapGroupChatFirestoreToFrontend(data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":             getStringValue(data, "id"),
		"projectId":      getStringValue(data, "project_id"),
		"createdByUid":   getStringValue(data, "created_by_uid"),
		"createdByName":  getStringValue(data, "created_by_name"),
		"name":           getStringValue(data, "name"),
		"description":    dereferenceString(getOptionalStringValue(data, "description"), ""),
		"participants":   mapParticipantsArrayToFrontend(GetParticipantsGoArray(data, "participants")),
		"messages":       mapBaseMessagesArrayToFrontend(GetBaseMessagesArrayFromFirestore(data, "messages")),
		"pinnedMessages": getStringArrayValue(data, "pinned_messages"),
		"isArchived":     getBoolValue(data, "is_archived"),
		"notifications":  getNotificationsMap(data, "notifications"),
		"createdAt":      getTimeValue(data, "created_at").Format(time.RFC3339),
		"updatedAt":      getTimeValue(data, "updated_at").Format(time.RFC3339),
		"readStatus":     getReadStatusMap(data, "read_status"),
		"groupSettings":  MapGroupSettingsFirestoreToFrontend(getMapValue(data, "group_settings")),
	}
}

// 4. MapGroupChatFirestoreToGo maps Firestore GroupChat data to Go struct format
func MapGroupChatFirestoreToGo(data map[string]interface{}) models.GroupChat {
	return models.GroupChat{
		ID:             getStringValue(data, "id"),
		ProjectID:      getStringValue(data, "project_id"),
		CreatedByUID:   getStringValue(data, "created_by_uid"),
		CreatedByName:  getStringValue(data, "created_by_name"),
		Name:           getStringValue(data, "name"),
		Description:    dereferenceString(getOptionalStringValue(data, "description"), ""),
		Participants:   GetParticipantsGoArray(data, "participants"),
		Messages:       GetBaseMessagesArrayFromFirestore(data, "messages"),
		PinnedMessages: getStringArrayValue(data, "pinned_messages"),
		IsArchived:     getBoolValue(data, "is_archived"),
		Notifications:  getNotificationsMap(data, "notifications"),
		CreatedAt:      getTimeValue(data, "created_at"),
		UpdatedAt:      getTimeValue(data, "updated_at"),
		ReadStatus:     getReadStatusMap(data, "read_status"),
		GroupSettings:  MapGroupSettingsFirestoreToGo(getMapValue(data, "group_settings")),
	}
}

// MapGroupChatGoToFrontend maps Go struct GroupChat data to frontend format
func MapGroupChatGoToFrontend(chat models.GroupChat) map[string]interface{} {
	return map[string]interface{}{
		"id":             chat.ID,
		"projectId":      chat.ProjectID,
		"createdByUid":   chat.CreatedByUID,
		"createdByName":  chat.CreatedByName,
		"name":           chat.Name,
		"description":    chat.Description,
		"participants":   mapParticipantsArrayToFrontend(chat.Participants),
		"messages":       mapBaseMessagesArrayToFrontend(chat.Messages),
		"pinnedMessages": chat.PinnedMessages,
		"isArchived":     chat.IsArchived,
		"notifications":  chat.Notifications,
		"createdAt":      chat.CreatedAt.Format(time.RFC3339),
		"updatedAt":      chat.UpdatedAt.Format(time.RFC3339),
		"readStatus":     chat.ReadStatus,
		"groupSettings":  MapGroupSettingsGoToFrontend(chat.GroupSettings),
	}
}

// MapGroupSettingsGoToFrontend maps Go struct GroupSettings to frontend format
func MapGroupSettingsGoToFrontend(settings models.GroupSettings) map[string]interface{} {
	return map[string]interface{}{
		"allowFileSharing":  settings.AllowFileSharing,
		"allowPinning":      settings.AllowPinning,
		"allowReactions":    settings.AllowReactions,
		"allowReplies":      settings.AllowReplies,
		"muteNotifications": settings.MuteNotifications,
		"onlyAdminsCanPost": settings.OnlyAdminsCanPost,
	}
}
