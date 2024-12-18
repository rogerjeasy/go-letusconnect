package models

import "time"

// Message represents a chat message
type Message struct {
	ID         string    `json:"id" firestore:"id"`
	SenderID   string    `json:"senderId" firestore:"sender_id"`
	ReceiverID string    `json:"receiverId" firestore:"receiver_id"`
	Content    string    `json:"content" firestore:"content"`
	CreatedAt  time.Time `json:"createdAt" firestore:"created_at"`
}
