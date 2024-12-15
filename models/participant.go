package models

import (
	"time"
)

type Participant struct {
	UserID         string    `json:"user_id"`
	Role           string    `json:"role"` // e.g., "owner", "collaborator"
	ProfilePicture string    `json:"profile_picture"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	JoinedAt       time.Time `json:"joined_at"`
}
