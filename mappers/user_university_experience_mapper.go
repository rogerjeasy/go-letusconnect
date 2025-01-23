package mappers

import (
	"github.com/rogerjeasy/go-letusconnect/models"
)

// MapUserSchoolExperienceFromGoToFirestore converts Go struct to Firestore format
func MapUserSchoolExperienceFromGoToFirestore(schoolExperience *models.UserSchoolExperience) map[string]interface{} {
	if schoolExperience == nil {
		return nil
	}

	universities := make([]map[string]interface{}, len(schoolExperience.Universities))
	for i, uni := range schoolExperience.Universities {
		universities[i] = map[string]interface{}{
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
		}
	}

	return map[string]interface{}{
		"uid":          schoolExperience.UID,
		"created_at":   schoolExperience.CreatedAt,
		"updated_at":   schoolExperience.UpdatedAt,
		"universities": universities,
	}
}

// MapUserSchoolExperienceFromFirestoreToGo converts Firestore format to Go struct
func MapUserSchoolExperienceFromFirestoreToGo(firestoreData map[string]interface{}) *models.UserSchoolExperience {
	if firestoreData == nil {
		return nil
	}

	experience := &models.UserSchoolExperience{
		UID:       getStringValue(firestoreData, "uid"),
		CreatedAt: getFirestoreTimeToGoTime(firestoreData["created_at"]),
		UpdatedAt: getFirestoreTimeToGoTime(firestoreData["updated_at"]),
	}

	if universities, ok := firestoreData["universities"].([]interface{}); ok {
		experience.Universities = make([]models.University, 0, len(universities))
		for _, uni := range universities {
			if uniMap, ok := uni.(map[string]interface{}); ok {
				university := models.University{
					ID:               getStringValue(uniMap, "id"),
					Name:             getStringValue(uniMap, "name"),
					Program:          getStringValue(uniMap, "program"),
					Country:          getStringValue(uniMap, "country"),
					City:             getStringValue(uniMap, "city"),
					StartYear:        getIntValueSafe(uniMap, "start_year"),
					EndYear:          getIntValueSafe(uniMap, "end_year"),
					Degree:           getStringValue(uniMap, "degree"),
					Experience:       getStringValue(uniMap, "experience"),
					Awards:           getStringArrayValue(uniMap, "awards"),
					Extracurriculars: getStringArrayValue(uniMap, "extracurriculars"),
				}
				experience.Universities = append(experience.Universities, university)
			}
		}
	}

	return experience
}

// MapUserSchoolExperienceFromFirestoreToFrontend converts Firestore format to frontend format
func MapUserSchoolExperienceFromFirestoreToFrontend(firestoreData map[string]interface{}) map[string]interface{} {
	if firestoreData == nil {
		return nil
	}

	universities := []map[string]interface{}{}
	if firestoreUnis, ok := firestoreData["universities"].([]interface{}); ok {
		for _, uni := range firestoreUnis {
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
		"uid":          getStringValue(firestoreData, "uid"),
		"createdAt":    getTimeValue(firestoreData, "created_at"),
		"updatedAt":    getTimeValue(firestoreData, "updated_at"),
		"universities": universities,
	}
}

// MapUniversityFromFrontendToGo converts frontend format to Go struct
func MapUniversityFromFrontendToGo(frontendData map[string]interface{}) models.University {
	if frontendData == nil {
		return models.University{}
	}

	return models.University{
		ID:               getStringValue(frontendData, "id"),
		Name:             getStringValue(frontendData, "name"),
		Program:          getStringValue(frontendData, "program"),
		Country:          getStringValue(frontendData, "country"),
		City:             getStringValue(frontendData, "city"),
		StartYear:        getIntValueSafe(frontendData, "startYear"),
		EndYear:          getIntValueSafe(frontendData, "endYear"),
		Degree:           getStringValue(frontendData, "degree"),
		Experience:       getStringValue(frontendData, "experience"),
		Awards:           getStringArrayValue(frontendData, "awards"),
		Extracurriculars: getStringArrayValue(frontendData, "extracurriculars"),
	}
}

// MapUserSchoolExperienceFromGoToFrontend converts Go struct to frontend format
func MapUserSchoolExperienceFromGoToFrontend(experience *models.UserSchoolExperience) map[string]interface{} {
	if experience == nil {
		return nil
	}

	universities := make([]map[string]interface{}, len(experience.Universities))
	for i, uni := range experience.Universities {
		universities[i] = map[string]interface{}{
			"id":               uni.ID,
			"name":             uni.Name,
			"program":          uni.Program,
			"country":          uni.Country,
			"city":             uni.City,
			"startYear":        uni.StartYear,
			"endYear":          uni.EndYear,
			"degree":           uni.Degree,
			"experience":       uni.Experience,
			"awards":           uni.Awards,
			"extracurriculars": uni.Extracurriculars,
		}
	}

	return map[string]interface{}{
		"uid":          experience.UID,
		"createdAt":    experience.CreatedAt,
		"updatedAt":    experience.UpdatedAt,
		"universities": universities,
	}
}
