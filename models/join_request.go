package models

import (
	"time"
)

type JoinRequest struct {
	UserID      string    `json:"user_id"`
	UserName    string    `json:"user_name"`
	Message     string    `json:"message"`
	RequestedAt time.Time `json:"requested_at"`
	Status      string    `json:"status"` // "pending", "accepted", "rejected"
}
