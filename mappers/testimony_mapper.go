package mappers

import (
	"time"

	"github.com/rogerjeasy/go-letusconnect/models"
)

// Testimonial Mappers
func MapTestimonialFrontendToGo(data map[string]interface{}) models.Testimonial {
	return models.Testimonial{
		ID:          getStringValue(data, "id"),
		Type:        models.TestimonialType(getStringValue(data, "type")),
		UserID:      getStringValue(data, "userId"),
		Title:       getStringValue(data, "title"),
		Content:     getStringValue(data, "content"),
		MediaURLs:   getStringArrayValue(data, "mediaUrls"),
		Tags:        getStringArrayValue(data, "tags"),
		CreatedAt:   getTimeValue(data, "createdAt"),
		UpdatedAt:   getTimeValue(data, "updatedAt"),
		IsPublished: getBoolValue(data, "isPublished"),
		Likes:       getIntValueSafe(data, "likes"),
	}
}

func MapTestimonialGoToFirestore(testimonial models.Testimonial) map[string]interface{} {
	return map[string]interface{}{
		"id":           testimonial.ID,
		"type":         string(testimonial.Type),
		"user_id":      testimonial.UserID,
		"title":        testimonial.Title,
		"content":      testimonial.Content,
		"media_urls":   testimonial.MediaURLs,
		"tags":         testimonial.Tags,
		"created_at":   testimonial.CreatedAt,
		"updated_at":   testimonial.UpdatedAt,
		"is_published": testimonial.IsPublished,
		"likes":        testimonial.Likes,
	}
}

func MapTestimonialFirestoreToGo(data map[string]interface{}) models.Testimonial {
	return models.Testimonial{
		ID:          getStringValue(data, "id"),
		Type:        models.TestimonialType(getStringValue(data, "type")),
		UserID:      getStringValue(data, "user_id"),
		Title:       getStringValue(data, "title"),
		Content:     getStringValue(data, "content"),
		MediaURLs:   getStringArrayValue(data, "media_urls"),
		Tags:        getStringArrayValue(data, "tags"),
		CreatedAt:   getFirestoreTimeToGoTime(data["created_at"]),
		UpdatedAt:   getFirestoreTimeToGoTime(data["updated_at"]),
		IsPublished: getBoolValue(data, "is_published"),
		Likes:       getIntValueSafe(data, "likes"),
	}
}

func MapTestimonialGoToFrontend(testimonial models.Testimonial) map[string]interface{} {
	return map[string]interface{}{
		"id":          testimonial.ID,
		"type":        string(testimonial.Type),
		"userId":      testimonial.UserID,
		"title":       testimonial.Title,
		"content":     testimonial.Content,
		"mediaUrls":   testimonial.MediaURLs,
		"tags":        testimonial.Tags,
		"createdAt":   testimonial.CreatedAt.Format(time.RFC3339),
		"updatedAt":   testimonial.UpdatedAt.Format(time.RFC3339),
		"isPublished": testimonial.IsPublished,
		"likes":       testimonial.Likes,
	}
}

// Alumni Testimonial Mappers
func MapAlumniTestimonialFrontendToGo(data map[string]interface{}) models.AlumniTestimonial {
	return models.AlumniTestimonial{
		Testimonial:      MapTestimonialFrontendToGo(data),
		GraduationYear:   getIntValueSafe(data, "graduationYear"),
		CurrentPosition:  getStringValue(data, "currentPosition"),
		CurrentCompany:   getStringValue(data, "currentCompany"),
		CareerHighlights: getStringArrayValue(data, "careerHighlights"),
		ProgramImpact:    getStringValue(data, "programImpact"),
		IndustryField:    getStringValue(data, "industryField"),
	}
}

func MapAlumniTestimonialGoToFirestore(testimonial models.AlumniTestimonial) map[string]interface{} {
	base := MapTestimonialGoToFirestore(testimonial.Testimonial)
	base["graduation_year"] = testimonial.GraduationYear
	base["current_position"] = testimonial.CurrentPosition
	base["current_company"] = testimonial.CurrentCompany
	base["career_highlights"] = testimonial.CareerHighlights
	base["program_impact"] = testimonial.ProgramImpact
	base["industry_field"] = testimonial.IndustryField
	return base
}

func MapAlumniTestimonialFirestoreToGo(data map[string]interface{}) models.AlumniTestimonial {
	return models.AlumniTestimonial{
		Testimonial:      MapTestimonialFirestoreToGo(data),
		GraduationYear:   getIntValueSafe(data, "graduation_year"),
		CurrentPosition:  getStringValue(data, "current_position"),
		CurrentCompany:   getStringValue(data, "current_company"),
		CareerHighlights: getStringArrayValue(data, "career_highlights"),
		ProgramImpact:    getStringValue(data, "program_impact"),
		IndustryField:    getStringValue(data, "industry_field"),
	}
}

func MapAlumniTestimonialGoToFrontend(testimonial models.AlumniTestimonial) map[string]interface{} {
	base := MapTestimonialGoToFrontend(testimonial.Testimonial)
	base["graduationYear"] = testimonial.GraduationYear
	base["currentPosition"] = testimonial.CurrentPosition
	base["currentCompany"] = testimonial.CurrentCompany
	base["careerHighlights"] = testimonial.CareerHighlights
	base["programImpact"] = testimonial.ProgramImpact
	base["industryField"] = testimonial.IndustryField
	return base
}

// Student Spotlight Mappers
func MapStudentSpotlightFrontendToGo(data map[string]interface{}) models.StudentSpotlight {
	return models.StudentSpotlight{
		Testimonial:        MapTestimonialFrontendToGo(data),
		CurrentSemester:    getIntValueSafe(data, "currentSemester"),
		ExpectedGraduation: getTimeValue(data, "expectedGraduation"),
		ResearchTopics:     getStringArrayValue(data, "researchTopics"),
		Projects:           MapProjectsArrayFrontendToGo(data, "projects"),
		Achievements:       getStringArrayValue(data, "achievements"),
	}
}

func MapStudentSpotlightGoToFirestore(spotlight models.StudentSpotlight) map[string]interface{} {
	base := MapTestimonialGoToFirestore(spotlight.Testimonial)
	base["current_semester"] = spotlight.CurrentSemester
	base["expected_graduation"] = spotlight.ExpectedGraduation
	base["research_topics"] = spotlight.ResearchTopics
	// base["projects"] = MapProjectGoToFirestore(spotlight.Projects)
	base["achievements"] = spotlight.Achievements
	return base
}

func MapStudentSpotlightFirestoreToGo(data map[string]interface{}) models.StudentSpotlight {
	return models.StudentSpotlight{
		Testimonial:        MapTestimonialFirestoreToGo(data),
		CurrentSemester:    getIntValueSafe(data, "current_semester"),
		ExpectedGraduation: getFirestoreTimeToGoTime(data["expected_graduation"]),
		ResearchTopics:     getStringArrayValue(data, "research_topics"),
		// Projects:           MapProjectFirestoreToGo(data, "projects"),
		Achievements: getStringArrayValue(data, "achievements"),
	}
}

func MapStudentSpotlightGoToFrontend(spotlight models.StudentSpotlight) map[string]interface{} {
	base := MapTestimonialGoToFrontend(spotlight.Testimonial)
	base["currentSemester"] = spotlight.CurrentSemester
	base["expectedGraduation"] = spotlight.ExpectedGraduation.Format(time.RFC3339)
	base["researchTopics"] = spotlight.ResearchTopics
	// base["projects"] = MapProjectsArrayGoToFrontend(spotlight.Projects)
	base["achievements"] = spotlight.Achievements
	return base
}

// Project Array Mappers
func MapProjectsArrayFrontendToGo(data map[string]interface{}, key string) []models.Project {
	var projects []models.Project
	if arr, ok := data[key].([]interface{}); ok {
		for _, item := range arr {
			if projectMap, ok := item.(map[string]interface{}); ok {
				projects = append(projects, MapProjectFrontendToGo(projectMap))
			}
		}
	}
	return projects
}
