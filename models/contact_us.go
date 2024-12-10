package models

import "time"

// Status constants for ContactUs
const (
	StatusUnread  = "unread"
	StatusRead    = "read"
	StatusReplied = "replied"
)

// ContactUs represents a contact form submission
type ContactUs struct {
	ID           string    `json:"id,omitempty" firestore:"id,omitempty"`
	Name         string    `json:"name" firestore:"name"`
	Email        string    `json:"email" firestore:"email"`
	Subject      string    `json:"subject" firestore:"subject"`
	Message      string    `json:"message" firestore:"message"`
	Attachments  []string  `json:"attachments,omitempty" firestore:"attachments,omitempty"`
	Status       string    `json:"status" firestore:"status"`
	CreatedAt    time.Time `json:"createdAt" firestore:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" firestore:"updated_at"`
	RepliedBy    string    `json:"repliedBy,omitempty" firestore:"replied_by,omitempty"`
	ReplyMessage string    `json:"replyMessage,omitempty" firestore:"reply_message,omitempty"`
}
