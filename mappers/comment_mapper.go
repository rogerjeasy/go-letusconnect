package mappers

import (
	"time"

	"github.com/rogerjeasy/go-letusconnect/models"
)

// 1. Map Comment from frontend format to Go struct
func MapCommentFrontendToGo(data map[string]interface{}) models.Comment {
	return models.Comment{
		UserID:    getStringValue(data, "userId"),
		UserName:  getStringValue(data, "userName"),
		Content:   getStringValue(data, "content"),
		CreatedAt: getTimeValue(data, "createdAt"),
	}
}

// 2. Map Comment from Go struct to Firestore format
func MapCommentGoToFirestore(comment models.Comment) map[string]interface{} {
	return map[string]interface{}{
		"user_id":    comment.UserID,
		"user_name":  comment.UserName,
		"content":    comment.Content,
		"created_at": comment.CreatedAt,
	}
}

// 3. Map Comment from Firestore format to frontend format
func MapCommentsFirestoreToClient(data interface{}) []map[string]interface{} {
	var result []map[string]interface{}

	// Handle Firestore data as []interface{}
	if comments, ok := data.([]interface{}); ok {
		for _, c := range comments {
			if commentMap, ok := c.(map[string]interface{}); ok {
				// Map from Firestore to Go struct
				comment := MapCommentFirestoreToGo(commentMap)

				// Map from Go struct to client format
				convertedComment := MapCommentGoToFrontend(comment)
				result = append(result, convertedComment)
			}
		}
	}

	return result
}

// 4. Map Comment from Firestore format to Go struct
func MapCommentFirestoreToGo(data map[string]interface{}) models.Comment {
	return models.Comment{
		UserID:    getStringValue(data, "user_id"),
		UserName:  getStringValue(data, "user_name"),
		Content:   getStringValue(data, "content"),
		CreatedAt: getTimeValue(data, "created_at"),
	}
}

func mapCommentsArrayToFrontend(data interface{}) []map[string]interface{} {
	var result []map[string]interface{}

	// Step 1: Convert Firestore data to Go structs
	var comments []models.Comment

	// Handle slice of interface{} (Firestore returns this type)
	if rawComments, ok := data.([]interface{}); ok {
		for _, c := range rawComments {
			if commentMap, ok := c.(map[string]interface{}); ok {
				comment := MapCommentFirestoreToGo(commentMap)
				comments = append(comments, comment)
			}
		}
	}

	// Handle slice of map[string]interface{}
	if rawComments, ok := data.([]map[string]interface{}); ok {
		for _, commentMap := range rawComments {
			comment := MapCommentFirestoreToGo(commentMap)
			comments = append(comments, comment)
		}
	}

	// Step 2: Convert Go structs to frontend format
	for _, comment := range comments {
		convertedComment := MapCommentGoToFrontend(comment)
		result = append(result, convertedComment)
	}

	return result
}

func MapCommentGoToFrontend(comment models.Comment) map[string]interface{} {
	return map[string]interface{}{
		"userId":    comment.UserID,
		"userName":  comment.UserName,
		"content":   comment.Content,
		"createdAt": comment.CreatedAt.Format(time.RFC3339),
	}
}
