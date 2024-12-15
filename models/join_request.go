package models

import (
	"time"
)

type JoinRequest struct {
	UserID         string    `json:"user_id"`
	Username       string    `json:"username"`
	Message        string    `json:"message"`
	ProfilePicture string    `json:"profile_picture"`
	Email          string    `json:"email"`
	RequestedAt    time.Time `json:"requested_at"`
	Status         string    `json:"status"` // "pending", "accepted", "rejected"
}
