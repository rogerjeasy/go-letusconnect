package models

import (
	"time"
)

type InvitedUser struct {
	UserID         string    `json:"user_id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	ProfilePicture string    `json:"profile_picture"`
	Role           string    `json:"role"` // e.g., "owner", "collaborator"
	JoinedAt       time.Time `json:"joined_at"`
}
