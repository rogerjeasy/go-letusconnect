package models

import (
	"time"
)

type Task struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // "todo", "in_progress", "done"
	DueDate     time.Time `json:"due_date"`
}
