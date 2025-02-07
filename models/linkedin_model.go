package models

type LinkedInUser struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	ProfilePicture string `json:"profilePicture,omitempty"`
}

type LinkedInJobApplication struct {
	ApplyTimestamp    int64  `json:"applyTimestamp" firestore:"apply_timestamp"`
	CompanyID         string `json:"companyId" firestore:"company_id"`
	CompanyName       string `json:"companyName" firestore:"company_name"`
	JobID             string `json:"jobId" firestore:"job_id"`
	JobTitle          string `json:"jobTitle" firestore:"job_title"`
	ApplicationStatus string `json:"applicationStatus" firestore:"application_status"`
}
