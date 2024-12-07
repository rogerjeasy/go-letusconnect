package mappers

import (
	"github.com/rogerjeasy/go-letusconnect/models"
)

// MapUserSchoolExperienceFrontendToBackend maps frontend UserSchoolExperience to Firestore-compatible format
func MapUserSchoolExperienceFrontendToBackend(schoolExperience *models.UserSchoolExperience) map[string]interface{} {
	universities := []map[string]interface{}{}

	for _, uni := range schoolExperience.Universities {
		universities = append(universities, map[string]interface{}{
			"id":               uni.ID,
			"name":             uni.Name,
			"program":          uni.Program,
			"country":          uni.Country,
			"city":             uni.City,
			"start_year":       uni.StartYear,
			"end_year":         uni.EndYear,
			"degree":           uni.Degree,
			"experience":       uni.Experience,
			"awards":           uni.Awards,
			"extracurriculars": uni.Extracurriculars,
		})
	}

	// Return formatted map
	return map[string]interface{}{
		"uid":          schoolExperience.UID,
		"created_at":   schoolExperience.CreatedAt,
		"updated_at":   schoolExperience.UpdatedAt,
		"universities": universities,
	}
}

// MapUserSchoolExperienceBackendToFrontend maps Firestore UserSchoolExperience to frontend format
func MapUserSchoolExperienceBackendToFrontend(backendData map[string]interface{}) map[string]interface{} {
	universities := []map[string]interface{}{}

	if backendUnis, ok := backendData["universities"].([]interface{}); ok {
		for _, uni := range backendUnis {
			if uniMap, ok := uni.(map[string]interface{}); ok {
				universities = append(universities, map[string]interface{}{
					"id":               getStringValue(uniMap, "id"),
					"name":             getStringValue(uniMap, "name"),
					"program":          getStringValue(uniMap, "program"),
					"country":          getStringValue(uniMap, "country"),
					"city":             getStringValue(uniMap, "city"),
					"startYear":        getIntValueSafe(uniMap, "start_year"),
					"endYear":          getIntValueSafe(uniMap, "end_year"),
					"degree":           getStringValue(uniMap, "degree"),
					"experience":       getStringValue(uniMap, "experience"),
					"awards":           getStringArrayValue(uniMap, "awards"),
					"extracurriculars": getStringArrayValue(uniMap, "extracurriculars"),
				})
			}
		}
	}

	return map[string]interface{}{
		"uid":          getStringValue(backendData, "uid"),
		"createdAt":    getStringValue(backendData, "created_at"),
		"updatedAt":    getStringValue(backendData, "updated_at"),
		"universities": universities,
	}
}

// MapFrontendToUniversity maps frontend data to University model
func MapFrontendToUniversity(data map[string]interface{}) models.University {

	return models.University{
		ID:               getStringValue(data, "id"),
		Name:             getStringValue(data, "name"),
		Program:          getStringValue(data, "program"),
		Country:          getStringValue(data, "country"),
		City:             getStringValue(data, "city"),
		StartYear:        getIntValueSafe(data, "startYear"),
		EndYear:          getIntValueSafe(data, "endYear"),
		Degree:           getStringValue(data, "degree"),
		Experience:       getStringValue(data, "experience"),
		Awards:           getStringArrayValue(data, "awards"),
		Extracurriculars: getStringArrayValue(data, "extracurriculars"),
	}
}
