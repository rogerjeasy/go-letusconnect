package models

import (
	"time"
)

// NotificationType defines the type of notification.
type NotificationType string

// NotificationStatus defines the status of the notification.
type NotificationStatus string

// NotificationPriority defines the priority level of the notification.
type NotificationPriority string

// Define constants for NotificationType
const (
	NotificationTypeMessage          NotificationType = "message"
	NotificationTypeEvent            NotificationType = "event"
	NotificationTypeReminder         NotificationType = "reminder"
	NotificationTypeSystem           NotificationType = "system"
	NotificationTypeCustom           NotificationType = "custom"
	NotificationTypeJob              NotificationType = "job"
	NotificationTypeProject          NotificationType = "project"
	NotificationTypeMentorship       NotificationType = "mentorship"
	NotificationTypeCollaboration    NotificationType = "collaboration"
	NotificationTypeResource         NotificationType = "resource"
	NotificationTypeTask             NotificationType = "task"
	NotificationTypeGoal             NotificationType = "goal"
	NotificationTypeSurvey           NotificationType = "survey"
	NotificationTypeFeedback         NotificationType = "feedback"
	NotificationTypeAnnouncement     NotificationType = "announcement"
	NotificationTypeMeeting          NotificationType = "meeting"
	NotificationTypeRequest          NotificationType = "request"
	NotificationTypeApproval         NotificationType = "approval"
	NotificationTypeReview           NotificationType = "review"
	NotificationTypeFeedbackLoop     NotificationType = "feedback_loop"
	NotificationTypeFeedbackForm     NotificationType = "feedback_form"
	NotificationTypeFeedbackPeer     NotificationType = "feedback_peer"
	NotificationTypeFeedback360      NotificationType = "feedback_360"
	NotificationTypeFeedbackSelf     NotificationType = "feedback_self"
	NotificationTypeFeedbackTeam     NotificationType = "feedback_team"
	NotificationTypeFeedbackManager  NotificationType = "feedback_manager"
	NotificationTypeFeedbackCustom   NotificationType = "feedback_custom"
	NotificationTypeFeedbackRequest  NotificationType = "feedback_request"
	NotificationTypeFeedbackResponse NotificationType = "feedback_response"
	NotificationTypeFeedbackReview   NotificationType = "feedback_review"
	NotificationTypeFeedbackReminder NotificationType = "feedback_reminder"
	NotificationTypeFeedbackSummary  NotificationType = "feedback_summary"
	NotificationTypeFeedbackReport   NotificationType = "feedback_report"
	NotificationTypeFeedbackAnalysis NotificationType = "feedback_analysis"
	NotificationTypeFeedbackAction   NotificationType = "feedback_action"
	NotificationTypeFeedbackGoal     NotificationType = "feedback_goal"
	NotificationTypeNewUser          NotificationType = "new_user"
	NotificationTypeNewMentor        NotificationType = "new_mentor"
	NotificationTypeNewMentee        NotificationType = "new_mentee"
	NotificationTypeNewRequest       NotificationType = "new_request"
	NotificationTypeNewReview        NotificationType = "new_review"
	NotificationTypeNewFeedback      NotificationType = "new_feedback"
)

// Define constants for NotificationStatus
const (
	NotificationStatusUnread   NotificationStatus = "unread"
	NotificationStatusRead     NotificationStatus = "read"
	NotificationStatusHidden   NotificationStatus = "hidden"
	NotificationStatusArchived NotificationStatus = "archived"
	NotificationStatusDeleted  NotificationStatus = "deleted"
)

// Define constants for NotificationPriority
const (
	NotificationPriorityLow    NotificationPriority = "low"
	NotificationPriorityNormal NotificationPriority = "normal"
	NotificationPriorityHigh   NotificationPriority = "high"
	NotificationPriorityUrgent NotificationPriority = "urgent"
)

// NotificationAction represents an actionable button or link in a notification.
type NotificationAction struct {
	Label      string `json:"label" firestore:"label"`
	URL        string `json:"url" firestore:"url"`
	IsPrimary  bool   `json:"isPrimary" firestore:"is_primary"`
	ActionType string `json:"actionType" firestore:"action_type"`
}

// EntityReference represents a reference to a related entity
type EntityReference struct {
	ID   string `json:"id" firestore:"id"`
	Type string `json:"type" firestore:"type"`
}

// Notification defines a unified structure for notifications.
type Notification struct {
	ID              string                 `json:"id" firestore:"id"`
	UserID          string                 `json:"userId" firestore:"user_id"`
	ActorID         string                 `json:"actorId,omitempty" firestore:"actor_id,omitempty"`
	ActorName       string                 `json:"actorName,omitempty" firestore:"actor_name,omitempty"`
	ActorType       string                 `json:"actorType,omitempty" firestore:"actor_type,omitempty"`
	Type            NotificationType       `json:"type" firestore:"type"`
	Title           string                 `json:"title" firestore:"title"`
	Content         string                 `json:"content" firestore:"content"`
	Category        string                 `json:"category" firestore:"category"`
	Priority        NotificationPriority   `json:"priority" firestore:"priority"`
	Status          NotificationStatus     `json:"status" firestore:"status"`
	RelatedEntities []EntityReference      `json:"relatedEntities,omitempty" firestore:"related_entities,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty" firestore:"metadata,omitempty"`
	Actions         []NotificationAction   `json:"actions,omitempty" firestore:"actions,omitempty"`
	ReadStatus      map[string]bool        `json:"readStatus" firestore:"read_status"`
	IsArchived      map[string]bool        `json:"isArchived" firestore:"is_archived"`
	IsImportant     bool                   `json:"isImportant" firestore:"is_important"`
	ExpiresAt       *time.Time             `json:"expiresAt,omitempty" firestore:"expires_at,omitempty"`
	CreatedAt       time.Time              `json:"createdAt" firestore:"created_at"`
	UpdatedAt       time.Time              `json:"updatedAt" firestore:"updated_at"`
	SentAt          *time.Time             `json:"sentAt,omitempty" firestore:"sent_at,omitempty"`
	ReadAt          *time.Time             `json:"readAt,omitempty" firestore:"read_at,omitempty"`
	Source          string                 `json:"source,omitempty" firestore:"source,omitempty"`
	Tags            []string               `json:"tags,omitempty" firestore:"tags,omitempty"`
	GroupID         string                 `json:"groupId,omitempty" firestore:"group_id,omitempty"`
	DeliveryChannel string                 `json:"deliveryChannel,omitempty" firestore:"delivery_channel,omitempty"`
	TargetedUsers   []string               `json:"targetedUsers,omitempty" firestore:"targeted_users,omitempty"`
	IsRead          bool                   `json:"isRead,omitempty" firestore:"is_read,omitempty"`
}
