package models

import (
	"time"
)

// TestimonialType represents the type of testimonial
type TestimonialType string

// const (
// 	AlumniTestimonial TestimonialType = "ALUMNI"
// 	StudentSpotlight  TestimonialType = "STUDENT"
// )

// Testimonial represents the base structure for all testimonials
type Testimonial struct {
	ID          string          `json:"id" firestore:"id"`
	Type        TestimonialType `json:"type" firestore:"type"`
	UserID      string          `json:"userId" firestore:"user_id"`
	Title       string          `json:"title" firestore:"title"`
	Content     string          `json:"content" firestore:"content"`
	MediaURLs   []string        `json:"mediaUrls" firestore:"media_urls"`
	Tags        []string        `json:"tags" firestore:"tags"`
	CreatedAt   time.Time       `json:"createdAt" firestore:"created_at"`
	UpdatedAt   time.Time       `json:"updatedAt" firestore:"updated_at"`
	IsPublished bool            `json:"isPublished" firestore:"is_published"`
	Likes       int             `json:"likes" firestore:"likes"`
}

// AlumniTestimonial extends Testimonial with alumni-specific fields
type AlumniTestimonial struct {
	Testimonial
	GraduationYear   int      `json:"graduationYear" firestore:"graduation_year"`
	CurrentPosition  string   `json:"currentPosition" firestore:"current_position"`
	CurrentCompany   string   `json:"currentCompany" firestore:"current_company"`
	CareerHighlights []string `json:"careerHighlights" firestore:"career_highlights"`
	ProgramImpact    string   `json:"programImpact" firestore:"program_impact"`
	IndustryField    string   `json:"industryField" firestore:"industry_field"`
}

// StudentSpotlight extends Testimonial with student-specific fields
type StudentSpotlight struct {
	Testimonial
	CurrentSemester    int       `json:"currentSemester" firestore:"current_semester"`
	ExpectedGraduation time.Time `json:"expectedGraduation" firestore:"expected_graduation"`
	ResearchTopics     []string  `json:"researchTopics" firestore:"research_topics"`
	Projects           []Project `json:"projects" firestore:"projects"`
	Achievements       []string  `json:"achievements" firestore:"achievements"`
}

// FirestoreCollections defines the collection names for Firestore
const (
	TestimonialsCollection = "testimonials"
	CommentsCollection     = "comments"
)
