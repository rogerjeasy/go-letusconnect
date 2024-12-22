package models

import "time"

// PollOption represents an individual option in a poll
type PollOption struct {
	ID        string `json:"id" firestore:"id"`                // Unique identifier for the option
	Text      string `json:"text" firestore:"text"`            // Option text
	VoteCount int    `json:"voteCount" firestore:"vote_count"` // Number of votes for the option
}

// Poll represents a poll created in a group chat
type Poll struct {
	ID                 string              `json:"id" firestore:"id"`                                    // Unique identifier for the poll
	Question           string              `json:"question" firestore:"question"`                        // Poll question
	Options            []PollOption        `json:"options" firestore:"options"`                          // Poll options
	CreatedBy          string              `json:"createdBy" firestore:"created_by"`                     // User ID of the creator
	CreatedAt          time.Time           `json:"createdAt" firestore:"created_at"`                     // Time of poll creation
	ExpiresAt          *time.Time          `json:"expiresAt,omitempty" firestore:"expires_at,omitempty"` // Poll expiration time (optional)
	AllowMultipleVotes bool                `json:"allowMultipleVotes" firestore:"allow_multiple_votes"`  // Whether multiple votes are allowed
	Votes              map[string][]string `json:"votes" firestore:"votes"`                              // Map of user IDs to voted option IDs
	IsClosed           bool                `json:"isClosed" firestore:"is_closed"`                       // Whether the poll is closed
}
