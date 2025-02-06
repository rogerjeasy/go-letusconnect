package models

import (
	"time"
)

// JobStatus represents the different states a job application can be in
type JobStatus string

const (
	Applied       JobStatus = "Applied"
	Interviewing  JobStatus = "Interviewing"
	OfferReceived JobStatus = "OfferReceived"
	Rejected      JobStatus = "Rejected"
	Withdrawn     JobStatus = "Withdrawn"
)

// InterviewType defines whether the interview is virtual or onsite
type InterviewType string

const (
	Virtual InterviewType = "Virtual"
	Onsite  InterviewType = "Onsite"
)

// ReminderType defines the types of reminders a user can set
type ReminderType string

const (
	None         ReminderType = "None"
	Email        ReminderType = "Email"
	Push         ReminderType = "Push Notification"
	SMS          ReminderType = "SMS"
	AllReminders ReminderType = "All"
)

// InterviewRound represents a scheduled interview round
type InterviewRound struct {
	RoundNumber   int           `json:"round_number" firestore:"round_number"`
	Date          time.Time     `json:"date" firestore:"date"`
	Time          string        `json:"time" firestore:"time"`
	Location      string        `json:"location" firestore:"location"`
	InterviewType InterviewType `json:"interview_type" firestore:"interview_type"`
	Interviewer   string        `json:"interviewer" firestore:"interviewer"`
	Description   string        `json:"description" firestore:"description"`
	Reminder      ReminderType  `json:"reminder" firestore:"reminder"`
	MeetingLink   string        `json:"meeting_link,omitempty" firestore:"meeting_link,omitempty"`
	Notes         string        `json:"notes" firestore:"notes"`
}

// Job represents a job application tracked by the user
type Job struct {
	ID              string           `json:"id" firestore:"id"`
	UserID          string           `json:"user_id" firestore:"user_id"`
	Company         string           `json:"company" firestore:"company"`
	Position        string           `json:"position" firestore:"position"`
	Location        string           `json:"location" firestore:"location"`
	ApplicationDate time.Time        `json:"application_date" firestore:"application_date"`
	Status          JobStatus        `json:"status" firestore:"status"`
	SalaryRange     string           `json:"salary_range,omitempty" firestore:"salary_range,omitempty"`
	JobType         string           `json:"job_type" firestore:"job_type"`
	JobDescription  string           `json:"job_description,omitempty" firestore:"job_description,omitempty"`
	JobPostLink     string           `json:"job_post_link,omitempty" firestore:"job_post_link,omitempty"`
	CompanyWebsite  string           `json:"company_website,omitempty" firestore:"company_website,omitempty"`
	Referral        string           `json:"referral,omitempty" firestore:"referral,omitempty"`
	Interviews      []InterviewRound `json:"interviews,omitempty" firestore:"interviews,omitempty"`
	OfferDetails    string           `json:"offer_details,omitempty" firestore:"offer_details,omitempty"`
	RejectionReason string           `json:"rejection_reason,omitempty" firestore:"rejection_reason,omitempty"`
	FollowUpDate    time.Time        `json:"follow_up_date,omitempty" firestore:"follow_up_date,omitempty"`
	CompanyRating   int              `json:"company_rating,omitempty" firestore:"company_rating,omitempty"`
	CreatedAt       time.Time        `json:"created_at" firestore:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at" firestore:"updated_at"`
}
