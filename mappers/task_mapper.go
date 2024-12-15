package mappers

import (
	"github.com/rogerjeasy/go-letusconnect/models"
)

// 1. Map Task from frontend format to Go struct
func MapTaskFrontendToGo(data map[string]interface{}) models.Task {
	return models.Task{
		ID:          getStringValue(data, "id"),
		Title:       getStringValue(data, "title"),
		Description: getStringValue(data, "description"),
		Status:      getStringValue(data, "status"),
		Priority:    getStringValue(data, "priority"),
		AssignedTo:  getParticipantsArray(data, "assignedTo"),
		DueDate:     getTimeValue(data, "dueDate"),
		CreatedAt:   getTimeValue(data, "createdAt"),
		UpdatedAt:   getTimeValue(data, "updatedAt"),
	}
}

// 2. Map Task from Go struct to Firestore format
func MapTaskGoToFirestore(task models.Task) map[string]interface{} {
	return map[string]interface{}{
		"id":          task.ID,
		"title":       task.Title,
		"description": task.Description,
		"status":      task.Status,
		"priority":    task.Priority,
		"assigned_to": mapParticipantsArrayToFirestore(task.AssignedTo),
		"due_date":    task.DueDate,
		"created_at":  task.CreatedAt,
		"updated_at":  task.UpdatedAt,
	}
}

// 3. Map Task from Firestore format to frontend format
func MapTaskFirestoreToFrontend(data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"id":          getStringValue(data, "id"),
		"title":       getStringValue(data, "title"),
		"description": getStringValue(data, "description"),
		"status":      getStringValue(data, "status"),
		"priority":    getStringValue(data, "priority"),
		"assignedTo":  mapInvitedUsersArrayToFrontend(data["assigned_to"]),
		"dueDate":     getTimeValue(data, "due_date").Format("2006-01-02T15:04:05Z07:00"),
		"createdAt":   getTimeValue(data, "created_at").Format("2006-01-02T15:04:05Z07:00"),
		"updatedAt":   getTimeValue(data, "updated_at").Format("2006-01-02T15:04:05Z07:00"),
	}
}

// 4. Map Task from Firestore database format to Go struct format
func MapTaskFirestoreToGo(data map[string]interface{}) models.Task {
	return models.Task{
		ID:          getStringValue(data, "id"),
		Title:       getStringValue(data, "title"),
		Description: getStringValue(data, "description"),
		Status:      getStringValue(data, "status"),
		Priority:    getStringValue(data, "priority"),
		AssignedTo:  getParticipantsArray(data, "assigned_to"),
		DueDate:     getTimeValue(data, "due_date"),
		CreatedAt:   getTimeValue(data, "created_at"),
		UpdatedAt:   getTimeValue(data, "updated_at"),
	}
}

// 5. Map Task Go struct to frontend format
func MapTaskGoToFrontend(task models.Task) map[string]interface{} {
	return map[string]interface{}{
		"id":          task.ID,
		"title":       task.Title,
		"description": task.Description,
		"status":      task.Status,
		"priority":    task.Priority,
		"assignedTo":  mapInvitedUsersArrayToFrontend(task.AssignedTo),
		"dueDate":     task.DueDate.Format("2006-01-02T15:04:05Z07:00"),
		"createdAt":   task.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		"updatedAt":   task.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// func mapAssignedUsersToFirestore(users []models.InvitedUser) []map[string]interface{} {
// 	var result []map[string]interface{}
// 	if users == nil {
// 		return result
// 	}

// 	for _, user := range users {
// 		result = append(result, map[string]interface{}{
// 			"user_id":   user.UserID,
// 			"role":      user.Role,
// 			"joined_at": user.JoinedAt.Format(time.RFC3339),
// 		})
// 	}
// 	return result
// }
