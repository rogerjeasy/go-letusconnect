package models

import "time"

type Project struct {
	ID                string        `json:"id"`
	Title             string        `json:"title"`
	Description       string        `json:"description"`
	OwnerID           string        `json:"owner_id"`
	OwnerUsername     string        `json:"owner_username"`
	CollaborationType string        `json:"collaboration_type"` // "public" or "private"
	SkillsNeeded      []string      `json:"skills_needed"`
	Industry          string        `json:"industry"`
	AcademicFields    []string      `json:"academic_fields"`
	Status            string        `json:"status"` // "open", "in_progress", "completed", "archived"
	Participants      []InvitedUser `json:"participants"`
	InvitedUsers      []InvitedUser `json:"invited_users"`
	JoinRequests      []JoinRequest `json:"join_requests"`
	Tasks             []Task        `json:"tasks"`
	Progress          string        `json:"progress"` // e.g., "50%", "Milestone 2/4"
	Comments          []Comment     `json:"comments"`
	ChatRoomID        string        `json:"chat_room_id"`
	Attachments       []Attachment  `json:"attachments"`
	Feedback          []Feedback    `json:"feedback"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
}
