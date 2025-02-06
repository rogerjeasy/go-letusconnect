package mappers

import (
	"time"

	"github.com/rogerjeasy/go-letusconnect/models"
)

// MapJobFrontendToGo maps frontend Job data to Go struct format
func MapJobFrontendToGo(data map[string]interface{}) models.Job {
	return models.Job{
		ID:              getStringValue(data, "id"),
		UserID:          getStringValue(data, "user_id"),
		Company:         getStringValue(data, "company"),
		Position:        getStringValue(data, "position"),
		Location:        getStringValue(data, "location"),
		ApplicationDate: getTimeValue(data, "application_date"),
		Status:          models.JobStatus(getStringValue(data, "status")),
		SalaryRange:     getStringValue(data, "salary_range"),
		JobType:         getStringValue(data, "job_type"),
		JobDescription:  getStringValue(data, "job_description"),
		JobPostLink:     getStringValue(data, "job_post_link"),
		CompanyWebsite:  getStringValue(data, "company_website"),
		Referral:        getStringValue(data, "referral"),
		Interviews:      getInterviewRoundsArrayFromFrontend(data, "interviews"),
		OfferDetails:    getStringValue(data, "offer_details"),
		RejectionReason: getStringValue(data, "rejection_reason"),
		FollowUpDate:    getTimeValue(data, "follow_up_date"),
		CompanyRating:   getOptionalIntValue(data, "company_rating"),
		CreatedAt:       getTimeValue(data, "created_at"),
		UpdatedAt:       getTimeValue(data, "updated_at"),
	}
}

// MapInterviewRoundFrontendToGo maps frontend InterviewRound data to Go struct format
func MapInterviewRoundFrontendToGo(data map[string]interface{}) models.InterviewRound {
	return models.InterviewRound{
		RoundNumber:   getIntValue(data, "round_number"),
		Date:          getTimeValue(data, "date"),
		Time:          getStringValue(data, "time"),
		Location:      getStringValue(data, "location"),
		InterviewType: models.InterviewType(getStringValue(data, "interview_type")),
		Interviewer:   getStringValue(data, "interviewer"),
		Description:   getStringValue(data, "description"),
		Reminder:      models.ReminderType(getStringValue(data, "reminder")),
		MeetingLink:   getStringValue(data, "meeting_link"),
		Notes:         getStringValue(data, "notes"),
	}
}

// MapJobGoToFirestore maps Go struct Job data to Firestore format
func MapJobGoToFirestore(job models.Job) map[string]interface{} {
	return map[string]interface{}{
		"id":               job.ID,
		"user_id":          job.UserID,
		"company":          job.Company,
		"position":         job.Position,
		"location":         job.Location,
		"application_date": job.ApplicationDate,
		"status":           string(job.Status),
		"salary_range":     job.SalaryRange,
		"job_type":         job.JobType,
		"job_description":  job.JobDescription,
		"job_post_link":    job.JobPostLink,
		"company_website":  job.CompanyWebsite,
		"referral":         job.Referral,
		"interviews":       mapInterviewRoundsArrayToFirestore(job.Interviews),
		"offer_details":    job.OfferDetails,
		"rejection_reason": job.RejectionReason,
		"follow_up_date":   job.FollowUpDate,
		"company_rating":   job.CompanyRating,
		"created_at":       job.CreatedAt,
		"updated_at":       job.UpdatedAt,
	}
}

// MapInterviewRoundGoToFirestore maps Go struct InterviewRound data to Firestore format
func MapInterviewRoundGoToFirestore(interview models.InterviewRound) map[string]interface{} {
	return map[string]interface{}{
		"round_number":   interview.RoundNumber,
		"date":           interview.Date,
		"time":           interview.Time,
		"location":       interview.Location,
		"interview_type": string(interview.InterviewType),
		"interviewer":    interview.Interviewer,
		"description":    interview.Description,
		"reminder":       string(interview.Reminder),
		"meeting_link":   interview.MeetingLink,
		"notes":          interview.Notes,
	}
}

// MapJobFirestoreToGo maps Firestore Job data to Go struct format
func MapJobFirestoreToGo(data map[string]interface{}) models.Job {
	return models.Job{
		ID:              getStringValue(data, "id"),
		UserID:          getStringValue(data, "user_id"),
		Company:         getStringValue(data, "company"),
		Position:        getStringValue(data, "position"),
		Location:        getStringValue(data, "location"),
		ApplicationDate: data["application_date"].(time.Time),
		Status:          models.JobStatus(getStringValue(data, "status")),
		SalaryRange:     getStringValue(data, "salary_range"),
		JobType:         getStringValue(data, "job_type"),
		JobDescription:  getStringValue(data, "job_description"),
		JobPostLink:     getStringValue(data, "job_post_link"),
		CompanyWebsite:  getStringValue(data, "company_website"),
		Referral:        getStringValue(data, "referral"),
		Interviews:      getInterviewRoundsArrayFromFirestore(data, "interviews"),
		OfferDetails:    getStringValue(data, "offer_details"),
		RejectionReason: getStringValue(data, "rejection_reason"),
		FollowUpDate:    data["follow_up_date"].(time.Time),
		CompanyRating:   getOptionalIntValue(data, "company_rating"),
		CreatedAt:       data["created_at"].(time.Time),
		UpdatedAt:       data["updated_at"].(time.Time),
	}
}

// MapInterviewRoundFirestoreToGo maps Firestore InterviewRound data to Go struct format
func MapInterviewRoundFirestoreToGo(data map[string]interface{}) models.InterviewRound {
	return models.InterviewRound{
		RoundNumber:   getIntValue(data, "round_number"),
		Date:          getTimeValue(data, "date"),
		Time:          getStringValue(data, "time"),
		Location:      getStringValue(data, "location"),
		InterviewType: models.InterviewType(getStringValue(data, "interview_type")),
		Interviewer:   getStringValue(data, "interviewer"),
		Description:   getStringValue(data, "description"),
		Reminder:      models.ReminderType(getStringValue(data, "reminder")),
		MeetingLink:   getStringValue(data, "meeting_link"),
		Notes:         getStringValue(data, "notes"),
	}
}

// MapJobGoToFrontend maps Go struct Job data to frontend format
func MapJobGoToFrontend(job models.Job) map[string]interface{} {
	return map[string]interface{}{
		"id":              job.ID,
		"userId":          job.UserID,
		"company":         job.Company,
		"position":        job.Position,
		"location":        job.Location,
		"applicationDate": job.ApplicationDate.Format(time.RFC3339),
		"status":          string(job.Status),
		"salaryRange":     job.SalaryRange,
		"jobType":         job.JobType,
		"jobDescription":  job.JobDescription,
		"jobPostLink":     job.JobPostLink,
		"companyWebsite":  job.CompanyWebsite,
		"referral":        job.Referral,
		"interviews":      mapInterviewRoundsArrayToFrontend(job.Interviews),
		"offerDetails":    job.OfferDetails,
		"rejectionReason": job.RejectionReason,
		"followUpDate":    job.FollowUpDate,
		"companyRating":   job.CompanyRating,
		"createdAt":       job.CreatedAt,
		"updatedAt":       job.UpdatedAt,
	}
}

// MapInterviewRoundGoToFrontend maps Go struct InterviewRound data to frontend format
func MapInterviewRoundGoToFrontend(interview models.InterviewRound) map[string]interface{} {
	return map[string]interface{}{
		"round_number":   interview.RoundNumber,
		"date":           interview.Date.Format(time.RFC3339),
		"time":           interview.Time,
		"location":       interview.Location,
		"interview_type": string(interview.InterviewType),
		"interviewer":    interview.Interviewer,
		"description":    interview.Description,
		"reminder":       string(interview.Reminder),
		"meeting_link":   interview.MeetingLink,
		"notes":          interview.Notes,
	}
}

// MapInterviewRoundsArrayToFirestore converts a slice of InterviewRound structs to Firestore format
func MapInterviewRoundsArrayToFirestore(rounds []models.InterviewRound) []map[string]interface{} {
	var firestoreRounds []map[string]interface{}

	for _, round := range rounds {
		firestoreRounds = append(firestoreRounds, map[string]interface{}{
			"round_number":   round.RoundNumber,
			"date":           round.Date,
			"time":           round.Time,
			"location":       round.Location,
			"interview_type": string(round.InterviewType),
			"interviewer":    round.Interviewer,
			"description":    round.Description,
			"reminder":       string(round.Reminder),
			"meeting_link":   round.MeetingLink,
			"notes":          round.Notes,
		})
	}

	return firestoreRounds
}

func getIntValue(data map[string]interface{}, key string) int {
	if value, ok := data[key]; ok {
		if num, ok := value.(int); ok {
			return num
		}
	}
	return 0
}

func getOptionalIntValue(data map[string]interface{}, key string) int {
	if value, ok := data[key]; ok {
		if num, ok := value.(int); ok {
			return num
		}
	}
	return 0
}

func getInterviewRoundsArrayFromFrontend(data map[string]interface{}, key string) []models.InterviewRound {
	var rounds []models.InterviewRound
	if value, ok := data[key]; ok {
		if slice, ok := value.([]interface{}); ok {
			for _, item := range slice {
				if m, ok := item.(map[string]interface{}); ok {
					rounds = append(rounds, MapInterviewRoundFrontendToGo(m))
				}
			}
		}
	}
	return rounds
}

func mapInterviewRoundsArrayToFirestore(rounds []models.InterviewRound) []map[string]interface{} {
	var result []map[string]interface{}
	for _, round := range rounds {
		result = append(result, MapInterviewRoundGoToFirestore(round))
	}
	return result
}

func getInterviewRoundsArrayFromFirestore(data map[string]interface{}, key string) []models.InterviewRound {
	var rounds []models.InterviewRound
	if value, ok := data[key]; ok {
		if slice, ok := value.([]interface{}); ok {
			for _, item := range slice {
				if m, ok := item.(map[string]interface{}); ok {
					rounds = append(rounds, MapInterviewRoundFirestoreToGo(m))
				}
			}
		}
	}
	return rounds
}

func mapInterviewRoundsArrayToFrontend(rounds []models.InterviewRound) []map[string]interface{} {
	var result []map[string]interface{}
	for _, round := range rounds {
		result = append(result, MapInterviewRoundGoToFrontend(round))
	}
	return result
}
