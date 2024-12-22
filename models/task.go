package models

import (
	"time"
)

type Task struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Status      string        `json:"status"`   // "todo", "in_progress", "done", "archived", "blocked"
	Priority    string        `json:"priority"` // "low", "medium", "high", "critical"
	AssignedTo  []Participant `json:"assigned_to"`
	DueDate     time.Time     `json:"due_date"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	CreatedBy   string        `json:"created_by"`
	UpdateBy    string        `json:"updated_by"`
}
