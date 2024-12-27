package models

import "time"

// NotificationType defines the type of notification.
type NotificationType string

// NotificationStatus defines the status of the notification.
const (
	NotificationTypeMessage  NotificationType = "message"
	NotificationTypeEvent    NotificationType = "event"
	NotificationTypeReminder NotificationType = "reminder"
	NotificationTypeSystem   NotificationType = "system"
	NotificationTypeCustom   NotificationType = "custom"
)

type NotificationStatus string

const (
	NotificationStatusUnread NotificationStatus = "unread"
	NotificationStatusRead   NotificationStatus = "read"
	NotificationStatusHidden NotificationStatus = "hidden"
)

// NotificationAction represents an actionable button or link in a notification.
type NotificationAction struct {
	Label     string `json:"label" firestore:"label"`
	URL       string `json:"url" firestore:"url"`
	IsPrimary bool   `json:"isPrimary" firestore:"is_primary"`
}

// Notification defines a unified structure for notifications.
type Notification struct {
	ID              string                 `json:"id" firestore:"id"`
	UserID          string                 `json:"userId" firestore:"user_id"`
	ActorID         string                 `json:"actorId,omitempty" firestore:"actor_id,omitempty"`
	ActorName       string                 `json:"actorName,omitempty" firestore:"actor_name,omitempty"`
	Type            NotificationType       `json:"type" firestore:"type"`
	Title           string                 `json:"title" firestore:"title"`
	Content         string                 `json:"content" firestore:"content"`
	Category        string                 `json:"category" firestore:"category"`
	Priority        string                 `json:"priority" firestore:"priority"`
	Status          NotificationStatus     `json:"status" firestore:"status"`
	RelatedEntityID string                 `json:"relatedEntityId,omitempty" firestore:"related_entity_id,omitempty"`
	RelatedEntity   map[string]interface{} `json:"relatedEntity,omitempty" firestore:"related_entity,omitempty"`
	Metadata        map[string]string      `json:"metadata,omitempty" firestore:"metadata,omitempty"`
	Actions         []NotificationAction   `json:"actions,omitempty" firestore:"actions,omitempty"`
	IsRead          bool                   `json:"isRead" firestore:"is_read"`
	IsArchived      bool                   `json:"isArchived" firestore:"is_archived"`
	IsImportant     bool                   `json:"isImportant" firestore:"is_important"`
	CreatedAt       time.Time              `json:"createdAt" firestore:"created_at"`
	UpdatedAt       time.Time              `json:"updatedAt" firestore:"updated_at"`
}
