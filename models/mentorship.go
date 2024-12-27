package models

import "time"

// UserType defines the type of user (Mentor or Mentee).
type UserType string

const (
	UserTypeMentor UserType = "mentor"
	UserTypeMentee UserType = "mentee"
	UserTypeBoth   UserType = "both" // User can be both mentor and mentee
)

// MentorshipRequestStatus defines the status of a mentorship request.
type MentorshipRequestStatus string

const (
	RequestStatusPending  MentorshipRequestStatus = "pending"
	RequestStatusApproved MentorshipRequestStatus = "approved"
	RequestStatusRejected MentorshipRequestStatus = "rejected"
)

// MentorshipRequest represents a request for mentorship.
type MentorshipRequest struct {
	ID        string                  `json:"id" firestore:"id"`
	MentorID  string                  `json:"mentorId" firestore:"mentor_id"`
	MenteeID  string                  `json:"menteeId" firestore:"mentee_id"`
	Status    MentorshipRequestStatus `json:"status" firestore:"status"`
	Message   string                  `json:"message,omitempty" firestore:"message,omitempty"`
	CreatedAt time.Time               `json:"createdAt" firestore:"created_at"`
	UpdatedAt time.Time               `json:"updatedAt" firestore:"updated_at"`
}

// Mentorship defines the relationship between mentors and mentees.
type Mentorship struct {
	ID                string              `json:"id" firestore:"id"`
	MentorID          string              `json:"mentorId,omitempty" firestore:"mentor_id,omitempty"`
	MentorName        string              `json:"mentorName,omitempty" firestore:"mentor_name,omitempty"`
	Mentees           []string            `json:"mentees,omitempty" firestore:"mentees,omitempty"`
	RejectedMentees   []string            `json:"rejectedMentees,omitempty" firestore:"rejected_mentees,omitempty"`
	PendingRequests   []MentorshipRequest `json:"pendingRequests,omitempty" firestore:"pending_requests,omitempty"`
	MenteesMentorship []MenteeMentorship  `json:"menteesMentorship,omitempty" firestore:"mentees_mentorship,omitempty"`
	CreatedAt         time.Time           `json:"createdAt" firestore:"created_at"`
	UpdatedAt         time.Time           `json:"updatedAt" firestore:"updated_at"`
}

// MenteeMentorship defines a mentee's relationship with mentors.
type MenteeMentorship struct {
	MentorID   string    `json:"mentorId" firestore:"mentor_id"`
	MentorName string    `json:"mentorName" firestore:"mentor_name"`
	StartedAt  time.Time `json:"startedAt" firestore:"started_at"`
}

// UserMentorshipProfile represents a user's mentorship relationships (as mentor or mentee).
type UserMentorshipProfile struct {
	UserID      string              `json:"userId" firestore:"user_id"`
	UserType    UserType            `json:"userType" firestore:"user_type"`
	Mentorships []Mentorship        `json:"mentorships" firestore:"mentorships"`
	Requests    []MentorshipRequest `json:"requests" firestore:"requests"`
	CreatedAt   time.Time           `json:"createdAt" firestore:"created_at"`
	UpdatedAt   time.Time           `json:"updatedAt" firestore:"updated_at"`
}
