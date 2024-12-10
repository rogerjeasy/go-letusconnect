package models

import "time"

// FAQ represents a frequently asked question
type FAQ struct {
	ID        string    `firestore:"id,omitempty" json:"id"`
	Username  string    `firestore:"uid" json:"username"`
	Question  string    `firestore:"question" json:"question"`
	Response  string    `firestore:"response" json:"response"`
	CreatedAt time.Time `firestore:"created_at" json:"createdAt"`
	UpdatedAt time.Time `firestore:"updated_at" json:"updatedAt"`
	CreatedBy string    `firestore:"created_by" json:"createdBy"`
	UpdatedBy string    `firestore:"updated_by" json:"updatedBy,omitempty"`
	Status    string    `firestore:"status" json:"status"` // e.g., "active" or "inactive"
	Category  string    `firestore:"category" json:"category"`
}
