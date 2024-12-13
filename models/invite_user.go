package models

import (
	"time"
)

type InvitedUser struct {
	UserID   string    `json:"user_id"`
	Role     string    `json:"role"` // e.g., "owner", "collaborator"
	JoinedAt time.Time `json:"joined_at"`
}
