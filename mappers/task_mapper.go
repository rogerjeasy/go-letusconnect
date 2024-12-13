package mappers

import (
	"github.com/rogerjeasy/go-letusconnect/models"
)

// 1. Map Task from frontend format to Go struct
func MapTaskFrontendToGo(data map[string]interface{}) models.Task {
	return models.Task{
		ID:          getStringValue(data, "id"),
		Description: getStringValue(data, "description"),
		Status:      getStringValue(data, "status"),
		DueDate:     getTimeValue(data, "dueDate"),
	}
}

// 2. Map Task from Go struct to Firestore format
func MapTaskGoToFirestore(task models.Task) map[string]interface{} {
	return map[string]interface{}{
		"id":          task.ID,
		"description": task.Description,
		"status":      task.Status,
		"due_date":    task.DueDate,
	}
}

// 3. Map Task from Firestore format to frontend format
func MapTaskFirestoreToFrontend(data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":          getStringValue(data, "id"),
		"description": getStringValue(data, "description"),
		"status":      getStringValue(data, "status"),
		"dueDate":     getTimeValue(data, "due_date"),
	}
}

// 4. Map Task from Firestore database format to Go struct format
func MapTaskFirestoreToGo(data map[string]interface{}) models.Task {
	return models.Task{
		ID:          getStringValue(data, "id"),
		Description: getStringValue(data, "description"),
		Status:      getStringValue(data, "status"),
		DueDate:     getTimeValue(data, "due_date"),
	}
}
