package models

import (
	"time"
)

type MessageConversation struct {
	ID        string    `json:"id" bson:"id"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	Message   string    `json:"message" bson:"message"`
	Response  string    `json:"response" bson:"response"`
	Role      string    `json:"role" bson:"role"` // "user" or "assistant"
}

type Conversation struct {
	ID        string                `json:"id" bson:"id"`
	UserID    string                `json:"user_id" bson:"user_id"`
	Title     string                `json:"title" bson:"title"`
	CreatedAt time.Time             `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time             `json:"updated_at" bson:"updated_at"`
	Messages  []MessageConversation `json:"messages" bson:"messages"`
}
