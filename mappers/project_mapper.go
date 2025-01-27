package mappers

import (
	"github.com/rogerjeasy/go-letusconnect/models"
)

// 1. MapProjectFrontendToGo maps frontend project data to Go struct format
func MapProjectFrontendToGo(data map[string]interface{}) models.Project {
	return models.Project{
		ID:                   getStringValue(data, "id"),
		Title:                getStringValue(data, "title"),
		Description:          getStringValue(data, "description"),
		OwnerID:              getStringValue(data, "ownerId"),
		OwnerUsername:        getStringValue(data, "ownerUsername"),
		CollaborationType:    getStringValue(data, "collaborationType"),
		SkillsNeeded:         getStringArrayValue(data, "skillsNeeded"),
		Industry:             getStringValue(data, "industry"),
		AcademicFields:       getStringArrayValue(data, "academicFields"),
		Status:               getStringValue(data, "status"),
		Participants:         GetParticipantsArray(data, "participants"),
		RejectedParticipants: getStringArrayValue(data, "rejectedParticipants"),
		InvitedUsers:         getInvitedUsersArray(data, "invitedUsers"),
		JoinRequests:         getJoinRequestsArray(data, "joinRequests"),
		Tasks:                getTasksArray(data, "tasks"),
		Progress:             getStringValue(data, "progress"),
		// Comments:             getCommentsArray(data, "comments"),
		ChatRoomID:  getStringValue(data, "chatRoomId"),
		Attachments: getAttachmentsArray(data, "attachments"),
		Feedback:    getFeedbacksArray(data, "feedback"),
		CreatedAt:   getTimeValue(data, "createdAt"),
		UpdatedAt:   getTimeValue(data, "updatedAt"),
	}
}

// 2. MapProjectGoToFirestore maps Go struct project data to Firestore format
func MapProjectGoToFirestore(project models.Project) map[string]interface{} {
	return map[string]interface{}{
		"id":                    project.ID,
		"title":                 project.Title,
		"description":           project.Description,
		"owner_id":              project.OwnerID,
		"owner_username":        project.OwnerUsername,
		"collaboration_type":    project.CollaborationType,
		"skills_needed":         project.SkillsNeeded,
		"industry":              project.Industry,
		"academic_fields":       project.AcademicFields,
		"status":                project.Status,
		"participants":          MapParticipantsArrayToFirestore(project.Participants),
		"rejected_participants": project.RejectedParticipants,
		"invited_users":         mapInvitedUsersArrayToFirestore(project.InvitedUsers),
		"join_requests":         mapJoinRequestsArrayToFirestore(project.JoinRequests),
		"tasks":                 mapTasksArrayToFirestore(project.Tasks),
		"progress":              project.Progress,
		// "comments":              mapCommentsArrayToFirestore(project.Comments),
		"chat_room_id": project.ChatRoomID,
		"attachments":  mapAttachmentsArrayToFirestore(project.Attachments),
		"feedback":     mapFeedbacksArrayToFirestore(project.Feedback),
		"created_at":   project.CreatedAt,
		"updated_at":   project.UpdatedAt,
	}
}

// 3. MapProjectFirestoreToFrontend maps Firestore project data to frontend format
func MapProjectFirestoreToFrontend(data map[string]interface{}) map[string]interface{} {

	return map[string]interface{}{
		"id":                   getStringValue(data, "id"),
		"title":                getStringValue(data, "title"),
		"description":          getStringValue(data, "description"),
		"ownerId":              getStringValue(data, "owner_id"),
		"ownerUsername":        getStringValue(data, "owner_username"),
		"collaborationType":    getStringValue(data, "collaboration_type"),
		"skillsNeeded":         getStringArrayValue(data, "skills_needed"),
		"industry":             getStringValue(data, "industry"),
		"academicFields":       getStringArrayValue(data, "academic_fields"),
		"status":               getStringValue(data, "status"),
		"participants":         mapParticipantsArrayToFrontend(data["participants"]),
		"rejectedParticipants": getStringArrayValue(data, "rejected_participants"),
		"invitedUsers":         mapInvitedUsersArrayToFrontend(data["invited_users"]),
		"joinRequests":         mapJoinRequestsArrayToFrontend(data["join_requests"]),
		"tasks":                mapTasksArrayToFrontend(data["tasks"]),
		"progress":             getStringValue(data, "progress"),
		// "comments":             mapCommentsArrayToFrontend(data["comments"]),
		"chatRoomId":  getStringValue(data, "chat_room_id"),
		"attachments": mapAttachmentsArrayToFrontend(data["attachments"]),
		"feedback":    mapFeedbacksArrayToFrontend(data["feedback"]),
		"createdAt":   getTimeValue(data, "created_at"),
		"updatedAt":   getTimeValue(data, "updated_at"),
	}
}

// 4. MapProjectFirestoreToGo maps Firestore project data to Go struct format
func MapProjectFirestoreToGo(data map[string]interface{}) models.Project {
	return models.Project{
		ID:                   getStringValue(data, "id"),
		Title:                getStringValue(data, "title"),
		Description:          getStringValue(data, "description"),
		OwnerID:              getStringValue(data, "owner_id"),
		OwnerUsername:        getStringValue(data, "owner_username"),
		CollaborationType:    getStringValue(data, "collaboration_type"),
		SkillsNeeded:         getStringArrayValue(data, "skills_needed"),
		Industry:             getStringValue(data, "industry"),
		AcademicFields:       getStringArrayValue(data, "academic_fields"),
		Status:               getStringValue(data, "status"),
		Participants:         GetParticipantsArray(data, "participants"),
		RejectedParticipants: getStringArrayValue(data, "rejected_participants"),
		InvitedUsers:         getInvitedUsersArray(data, "invited_users"),
		JoinRequests:         getJoinRequestsArray(data, "join_requests"),
		Tasks:                getTasksArray(data, "tasks"),
		Progress:             getStringValue(data, "progress"),
		// Comments:             getCommentsArray(data, "comments"),
		ChatRoomID:  getStringValue(data, "chat_room_id"),
		Attachments: getAttachmentsArray(data, "attachments"),
		Feedback:    getFeedbacksArray(data, "feedback"),
		CreatedAt:   getTimeValue(data, "created_at"),
		UpdatedAt:   getTimeValue(data, "updated_at"),
	}
}
