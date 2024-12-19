package mappers

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/rogerjeasy/go-letusconnect/models"
)

func getJoinRequestsArray(data map[string]interface{}, key string) []models.JoinRequest {
	if value, ok := data[key].([]interface{}); ok {
		var joinRequests []models.JoinRequest
		for _, v := range value {
			if joinRequestData, ok := v.(map[string]interface{}); ok {
				joinRequest := MapJoinRequestFrontendToGo(joinRequestData)
				joinRequests = append(joinRequests, joinRequest)
			}
		}
		return joinRequests
	}
	return []models.JoinRequest{}
}

func getTasksArray(data map[string]interface{}, key string) []models.Task {
	if value, ok := data[key].([]interface{}); ok {
		var tasks []models.Task
		for _, v := range value {
			if taskData, ok := v.(map[string]interface{}); ok {
				task := models.Task{
					ID:          getStringValue(taskData, "id"),
					Description: getStringValue(taskData, "description"),
					Status:      getStringValue(taskData, "status"),
					DueDate:     getTimeValue(taskData, "dueDate"),
				}
				tasks = append(tasks, task)
			}
		}
		return tasks
	}
	return []models.Task{}
}

func getCommentsArray(data map[string]interface{}, key string) []models.Comment {
	if value, ok := data[key].([]interface{}); ok {
		var comments []models.Comment
		for _, v := range value {
			if commentData, ok := v.(map[string]interface{}); ok {
				comment := models.Comment{
					UserID:    getStringValue(commentData, "userId"),
					UserName:  getStringValue(commentData, "userName"),
					Content:   getStringValue(commentData, "content"),
					CreatedAt: getTimeValue(commentData, "createdAt"),
				}
				comments = append(comments, comment)
			}
		}
		return comments
	}
	return []models.Comment{}
}

func getAttachmentsArray(data map[string]interface{}, key string) []models.Attachment {
	if value, ok := data[key].([]interface{}); ok {
		var attachments []models.Attachment
		for _, v := range value {
			if attachmentData, ok := v.(map[string]interface{}); ok {
				attachment := models.Attachment{
					FileName:   getStringValue(attachmentData, "fileName"),
					URL:        getStringValue(attachmentData, "url"),
					UploadedAt: getTimeValue(attachmentData, "uploadedAt"),
				}
				attachments = append(attachments, attachment)
			}
		}
		return attachments
	}
	return []models.Attachment{}
}

func getFeedbacksArray(data map[string]interface{}, key string) []models.Feedback {
	if value, ok := data[key].([]interface{}); ok {
		var feedbacks []models.Feedback
		for _, v := range value {
			if feedbackData, ok := v.(map[string]interface{}); ok {
				feedback := models.Feedback{
					UserID:    getStringValue(feedbackData, "userId"),
					Rating:    getIntValueSafe(feedbackData, "rating"),
					Comment:   getStringValue(feedbackData, "comment"),
					CreatedAt: getTimeValue(feedbackData, "createdAt"),
				}
				feedbacks = append(feedbacks, feedback)
			}
		}
		return feedbacks
	}
	return []models.Feedback{}
}

func mapJoinRequestsArrayToFirestore(joinRequests []models.JoinRequest) []map[string]interface{} {
	var result []map[string]interface{}
	for _, jr := range joinRequests {
		result = append(result, MapJoinRequestGoToFirestore(jr))
	}
	return result
}

func mapTasksArrayToFirestore(tasks []models.Task) []map[string]interface{} {
	var result []map[string]interface{}
	for _, task := range tasks {
		result = append(result, MapTaskGoToFirestore(task))
	}
	return result
}

func mapCommentsArrayToFirestore(comments []models.Comment) []map[string]interface{} {
	var result []map[string]interface{}
	for _, comment := range comments {
		result = append(result, MapCommentGoToFirestore(comment))
	}
	return result
}

func mapAttachmentsArrayToFirestore(attachments []models.Attachment) []map[string]interface{} {
	var result []map[string]interface{}
	for _, attachment := range attachments {
		result = append(result, MapAttachmentGoToFirestore(attachment))
	}
	return result
}

func mapFeedbacksArrayToFirestore(feedbacks []models.Feedback) []map[string]interface{} {
	var result []map[string]interface{}
	for _, feedback := range feedbacks {
		result = append(result, MapFeedbackGoToFirestore(feedback))
	}
	return result
}

// Helper function to safely get string values

func getStringValue(data map[string]interface{}, key string) string {
	if value, ok := data[key]; ok {
		if value == nil {
			return ""
		}

		if strVal, isString := value.(string); isString {
			return strVal
		} else {
			fmt.Printf("Key '%s' is of unexpected type %T with value: %+v\n", key, value, value)
		}
	} else {
		fmt.Printf("Key '%s' not found in data\n", key)
	}

	return ""
}

func getTimeValue(data map[string]interface{}, key string) time.Time {
	if value, ok := data[key]; ok {
		if timeStr, isString := value.(string); isString {
			parsedTime, err := time.Parse(time.RFC3339, timeStr)
			if err != nil {
				log.Printf("Error parsing time for key %s: %v", key, err)
				return time.Time{} // return zero value of time.Time in case of an error
			}
			return parsedTime
		}
	}
	return time.Time{} // return zero value of time.Time if key not found or type mismatch
}

// Helper function to safely get integer values
func getIntValueSafe(data map[string]interface{}, key string) int {
	if val, ok := data[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case int64:
			return int(v)
		case float64:
			return int(v)
		case string:
			if intValue, err := strconv.Atoi(v); err == nil {
				return intValue
			}
		}
	}
	return 0
}

func getInvitedUsersArray(data map[string]interface{}, key string) []models.InvitedUser {
	result := []models.InvitedUser{}

	if value, ok := data[key]; ok {
		// Ensure the value is a slice of interfaces
		if participants, ok := value.([]interface{}); ok {
			for _, v := range participants {
				if userMap, ok := v.(map[string]interface{}); ok {
					user := MapInvitedUserFrontendToGo(userMap)
					result = append(result, user)
				} else {
					fmt.Println("Error: participant is not a map[string]interface{}")
				}
			}
		} else {
			fmt.Println("Error: participants is not a []interface{}")
		}
	} else {
		fmt.Println("Error: participants key not found in data")
	}

	return result
}

func getParticipantsArray(data map[string]interface{}, key string) []models.Participant {
	result := []models.Participant{}

	if value, ok := data[key]; ok {
		// Ensure the value is a slice of interfaces
		if participants, ok := value.([]interface{}); ok {
			for _, v := range participants {
				if userMap, ok := v.(map[string]interface{}); ok {
					user := MapParticipantFrontendToGo(userMap)
					result = append(result, user)
				} else {
					fmt.Println("Error: participant is not a map[string]interface{}")
				}
			}
		} else {
			fmt.Println("Error: participants is not a []interface{}")
		}
	} else {
		fmt.Println("Error: participants key not found in data")
	}

	return result
}

// Converts a slice of InvitedUser structs to Firestore format
func mapInvitedUsersArrayToFirestore(users []models.InvitedUser) []map[string]interface{} {
	var result []map[string]interface{}
	for _, user := range users {
		result = append(result, MapInvitedUserGoToFirestore(user))
	}
	return result
}

// Converts a slice of Participant structs to Firestore format
func mapParticipantsArrayToFirestore(users []models.Participant) []map[string]interface{} {
	var result []map[string]interface{}
	for _, user := range users {
		result = append(result, MapParticipantGoToFirestore(user))
	}
	return result
}

// getStringArrayValue safely extracts a slice of strings from a map
func getStringArrayValue(data map[string]interface{}, key string) []string {
	if value, ok := data[key]; ok {
		// Handle case where value is a slice of interfaces
		if interfaceSlice, ok := value.([]interface{}); ok {
			var result []string
			for _, item := range interfaceSlice {
				if str, ok := item.(string); ok {
					result = append(result, str)
				} else {
					fmt.Printf("Warning: item in %s is not a string: %+v\n", key, item)
				}
			}
			return result
		}

		// Handle case where value is a slice of strings directly
		if stringSlice, ok := value.([]string); ok {
			return stringSlice
		}

		// Log if the value type is unexpected
		fmt.Printf("Warning: unexpected type for %s: %+v\n", key, value)
	}
	return []string{}
}

// func getMapValue(data map[string]interface{}, key string) map[string]bool {
// 	if val, ok := data[key].(map[string]interface{}); ok {
// 		readStatus := make(map[string]bool)
// 		for k, v := range val {
// 			if b, ok := v.(bool); ok {
// 				readStatus[k] = b
// 			}
// 		}
// 		return readStatus
// 	}
// 	return nil
// }

// func getStringPointer(data map[string]interface{}, key string) *string {
// 	if val, ok := data[key].(string); ok {
// 		return &val
// 	}
// 	return nil
// }

// func getSliceValue(data map[string]interface{}, key string) []string {
// 	if val, ok := data[key].([]interface{}); ok {
// 		var result []string
// 		for _, v := range val {
// 			if str, ok := v.(string); ok {
// 				result = append(result, str)
// 			}
// 		}
// 		return result
// 	}
// 	return nil
// }

// func getMapReactions(data map[string]interface{}, key string) map[string]int {
// 	if val, ok := data[key].(map[string]interface{}); ok {
// 		reactions := make(map[string]int)
// 		for k, v := range val {
// 			if i, ok := v.(float64); ok {
// 				reactions[k] = int(i)
// 			}
// 		}
// 		return reactions
// 	}
// 	return nil
// }

// Helper function to safely get a boolean value from a map
func getBoolValue(data map[string]interface{}, key string) bool {
	if val, ok := data[key].(bool); ok {
		return val
	}
	return false
}

// Helper function to safely get a map of read status
func getReadStatusMap(data map[string]interface{}, key string) map[string]bool {
	if val, ok := data[key].(map[string]interface{}); ok {
		readStatus := make(map[string]bool)
		for k, v := range val {
			if status, ok := v.(bool); ok {
				readStatus[k] = status
			}
		}
		return readStatus
	}
	return nil
}

// Helper function to safely get a map of reactions
func getReactionsMap(data map[string]interface{}, key string) map[string]int {
	if val, ok := data[key].(map[string]interface{}); ok {
		reactions := make(map[string]int)
		for k, v := range val {
			if count, ok := v.(float64); ok {
				reactions[k] = int(count)
			}
		}
		return reactions
	}
	return nil
}

// Helper function to safely get an optional string pointer
func getOptionalStringValue(data map[string]interface{}, key string) *string {
	if val, ok := data[key].(string); ok {
		return &val
	}
	return nil
}

// Helper function to get time value as a string
func getTimeStringValue(data map[string]interface{}, key string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return time.Now().Format(time.RFC3339)
}

// getDirectMessagesArrayFromFrontend converts frontend direct messages array to Go struct array
func getDirectMessagesArrayFromFrontend(data map[string]interface{}, key string) []models.DirectMessage {
	var directMessages []models.DirectMessage

	if messages, exists := data[key]; exists {
		if messageArray, ok := messages.([]interface{}); ok {
			for _, msg := range messageArray {
				if msgMap, ok := msg.(map[string]interface{}); ok {
					directMessages = append(directMessages, MapDirectMessageFrontendToGo(msgMap))
				}
			}
		}
	}
	return directMessages
}

// mapDirectMessagesArrayToFirestore converts Go struct direct messages array to Firestore format
func mapDirectMessagesArrayToFirestore(messages []models.DirectMessage) []map[string]interface{} {
	var result []map[string]interface{}

	for _, msg := range messages {
		result = append(result, MapDirectMessageGoToFirestore(msg))
	}
	return result
}

// getDirectMessagesArrayFromFirestore converts Firestore direct messages array to frontend format
func getDirectMessagesArrayFromFirestore(data map[string]interface{}, key string) []map[string]interface{} {
	var result []map[string]interface{}

	if messages, exists := data[key]; exists {
		if messageArray, ok := messages.([]interface{}); ok {
			for _, msg := range messageArray {
				if msgMap, ok := msg.(map[string]interface{}); ok {
					result = append(result, MapDirectMessageFirestoreToFrontend(msgMap))
				}
			}
		}
	}
	return result
}

// getDirectMessagesArrayFromFirestoreToGo converts Firestore direct messages array to Go struct format
func getDirectMessagesArrayFromFirestoreToGo(data map[string]interface{}, key string) []models.DirectMessage {
	var directMessages []models.DirectMessage

	if messages, exists := data[key]; exists {
		if messageArray, ok := messages.([]interface{}); ok {
			for _, msg := range messageArray {
				if msgMap, ok := msg.(map[string]interface{}); ok {
					directMessages = append(directMessages, MapDirectMessageFirestoreToGo(msgMap))
				}
			}
		}
	}
	return directMessages
}
