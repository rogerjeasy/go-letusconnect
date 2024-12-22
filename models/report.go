package models

import "time"

// Report represents a report made by a user
type Report struct {
	ID          string    `json:"id" firestore:"id"`
	ReporterID  string    `json:"reporterId" firestore:"reporter_id"`
	MessageID   string    `json:"messageId" firestore:"message_id"`
	ReportedBy  string    `json:"reportedBy" firestore:"reported_by"`
	GroupChatID string    `json:"groupChatId" firestore:"group_chat_id"`
	Reason      string    `json:"reason" firestore:"reason"`
	Description string    `json:"description,omitempty" firestore:"description,omitempty"`
	Status      string    `json:"status" firestore:"status"` // e.g., "pending", "resolved", "rejected"
	CreatedAt   time.Time `json:"createdAt" firestore:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" firestore:"updated_at"`
	ResolvedBy  *string   `json:"resolvedBy,omitempty" firestore:"resolved_by,omitempty"`
}
