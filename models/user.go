package models

import "time"

// User represents the user model for the application
type User struct {
	UID              string   `json:"uid"`
	Username         string   `json:"username"`
	FirstName        string   `json:"first_name"`
	LastName         string   `json:"last_name"`
	Email            string   `json:"email"`
	PhoneNumber      string   `json:"phone_number"`
	ProfilePicture   string   `json:"profile_picture"`
	Bio              string   `json:"bio"`
	Role             []string `json:"role"`
	GraduationYear   int      `json:"graduation_year"`
	CurrentJobTitle  string   `json:"current_job_title"`
	AreasOfExpertise []string `json:"areas_of_expertise"`
	Interests        []string `json:"interests"`
	LookingForMentor bool     `json:"looking_for_mentor"`
	WillingToMentor  bool     `json:"willing_to_mentor"`
	ConnectionsMade  int      `json:"connections_made"`
	AccountCreatedAt string   `json:"account_creation_date"`
	IsActive         bool     `json:"is_active"`
	IsVerified       bool     `json:"is_verified"`
	Password         string   `json:"password"`
	Program          string   `json:"program"`
	DateOfBirth      string   `json:"date_of_birth"`
	PhoneCode        string   `json:"phone_code"`
	Languages        []string `json:"languages"`
	Skills           []string `json:"skills"`
	Certifications   []string `json:"certifications"`
	Projects         []string `json:"projects"`
	IsOnline         bool     `json:"is_online"`
}

type UserConnections struct {
	ID              string                       `json:"id" bson:"id"`
	UID             string                       `json:"uid" bson:"uid"`
	Connections     map[string]Connection        `json:"connections" bson:"connections"`
	PendingRequests map[string]ConnectionRequest `json:"pending_requests" bson:"pending_requests"`
}

type Connection struct {
	TargetUID  string    `json:"target_uid" bson:"target_uid"`
	TargetName string    `json:"target_name" bson:"target_name"`
	SentAt     time.Time `json:"sent_at" bson:"sent_at"`
	AcceptedAt time.Time `json:"accepted_at" bson:"accepted_at"`
	Status     string    `json:"status" bson:"status"` // active, blocked, etc.
}

type ConnectionRequest struct {
	FromUID  string    `json:"from_uid" bson:"from_uid"`
	FromName string    `json:"from_name" bson:"from_name"`
	ToUID    string    `json:"to_uid" bson:"to_uid"`
	SentAt   time.Time `json:"sent_at" bson:"sent_at"`
	Message  string    `json:"message" bson:"message"`
	Status   string    `json:"status" bson:"status"` // pending, accepted, rejected
}

type UserAddress struct {
	ID          string `json:"id"`
	UID         string `json:"uid"`
	Street      string `json:"street"`
	City        string `json:"city"`
	State       string `json:"state"`
	Country     string `json:"country"`
	PostalCode  int    `json:"postal_code"`
	HouseNumber int    `json:"house_number"`
	Apartment   string `json:"apartment"`
	Region      string `json:"region"`
}

type UserSchoolExperience struct {
	UID          string       `json:"uid"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	Universities []University `json:"universities"`
}

type University struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Program          string   `json:"program"`
	Country          string   `json:"country"`
	City             string   `json:"city"`
	StartYear        int      `json:"start_year"`
	EndYear          int      `json:"end_year"`
	Degree           string   `json:"degree"`
	Experience       string   `json:"experience"`
	Awards           []string `json:"awards"`
	Extracurriculars []string `json:"extracurriculars"`
}

type UserWorkExperience struct {
	ID              string           `json:"id"`
	UID             string           `json:"uid"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	WorkExperiences []WorkExperience `json:"work_experiences"`
}

type WorkExperience struct {
	ID               string    `json:"id"`
	Company          string    `json:"company"`
	Position         string    `json:"position"`
	City             string    `json:"city"`
	Country          string    `json:"country"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	Responsibilities []string  `json:"responsibilities"`
	Achievements     []string  `json:"achievements"`
}

type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
