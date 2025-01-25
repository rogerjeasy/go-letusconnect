package mappers

import (
	"time"

	"github.com/rogerjeasy/go-letusconnect/models"
)

// Frontend to Go
func MapGroupFrontendToGo(data map[string]interface{}) models.Group {
	return models.Group{
		ID:            getStringValue(data, "id"),
		Name:          getStringValue(data, "name"),
		Description:   getStringValue(data, "description"),
		ImageURL:      getStringValue(data, "imageUrl"),
		Category:      MapGroupCategoryFrontendToGo(getMapValue(data, "category")),
		ActivityLevel: getStringValue(data, "activityLevel"),
		CreatedAt:     getTimeValue(data, "createdAt"),
		UpdatedAt:     getTimeValue(data, "updatedAt"),
		Members:       MapMembersArrayFrontendToGo(data, "members"),
		Topics:        MapTopicsArrayFrontendToGo(data, "topics"),
		Events:        MapEventsArrayFrontendToGo(data, "events"),
		Resources:     MapResourcesArrayFrontendToGo(data, "resources"),
		Admins:        MapAdminsArrayFrontendToGo(data, "admins"),
		Privacy:       getStringValue(data, "privacy"),
		Rules:         MapRulesArrayFrontendToGo(data, "rules"),
		Featured:      getBoolValue(data, "featured"),
		Size:          getStringValue(data, "size"),
	}
}

// Go to Firestore
func MapGroupGoToFirestore(group models.Group) map[string]interface{} {
	return map[string]interface{}{
		"id":             group.ID,
		"name":           group.Name,
		"description":    group.Description,
		"image_url":      group.ImageURL,
		"category":       MapGroupCategoryGoToFirestore(group.Category),
		"activity_level": group.ActivityLevel,
		"created_at":     group.CreatedAt,
		"updated_at":     group.UpdatedAt,
		"members":        MapMembersArrayGoToFirestore(group.Members),
		"topics":         MapTopicsArrayGoToFirestore(group.Topics),
		"events":         MapEventsArrayGoToFirestore(group.Events),
		"resources":      MapResourcesArrayGoToFirestore(group.Resources),
		"admins":         MapAdminsArrayGoToFirestore(group.Admins),
		"privacy":        group.Privacy,
		"rules":          MapRulesArrayGoToFirestore(group.Rules),
		"featured":       group.Featured,
		"size":           group.Size,
	}
}

// Firestore to Go
func MapGroupFirestoreToGo(data map[string]interface{}) models.Group {
	return models.Group{
		ID:            getStringValue(data, "id"),
		Name:          getStringValue(data, "name"),
		Description:   getStringValue(data, "description"),
		ImageURL:      getStringValue(data, "image_url"),
		Category:      MapGroupCategoryFirestoreToGo(getMapValue(data, "category")),
		ActivityLevel: getStringValue(data, "activity_level"),
		CreatedAt:     getFirestoreTimeToGoTime(data["created_at"]),
		UpdatedAt:     getFirestoreTimeToGoTime(data["updated_at"]),
		Members:       MapMembersArrayFirestoreToGo(data, "members"),
		Topics:        MapTopicsArrayFirestoreToGo(data, "topics"),
		Events:        MapEventsArrayFirestoreToGo(data, "events"),
		Resources:     MapResourcesArrayFirestoreToGo(data, "resources"),
		Admins:        MapAdminsArrayFirestoreToGo(data, "admins"),
		Privacy:       getStringValue(data, "privacy"),
		Rules:         MapRulesArrayFirestoreToGo(data, "rules"),
		Featured:      getBoolValue(data, "featured"),
		Size:          getStringValue(data, "size"),
	}
}

// Go to Frontend
func MapGroupGoToFrontend(group models.Group) map[string]interface{} {
	return map[string]interface{}{
		"id":            group.ID,
		"name":          group.Name,
		"description":   group.Description,
		"imageUrl":      group.ImageURL,
		"category":      MapGroupCategoryGoToFrontend(group.Category),
		"activityLevel": group.ActivityLevel,
		"createdAt":     group.CreatedAt.Format(time.RFC3339),
		"updatedAt":     group.UpdatedAt.Format(time.RFC3339),
		"members":       MapMembersArrayGoToFrontend(group.Members),
		"topics":        MapTopicsArrayGoToFrontend(group.Topics),
		"events":        MapEventsArrayGoToFrontend(group.Events),
		"resources":     MapResourcesArrayGoToFrontend(group.Resources),
		"admins":        MapAdminsArrayGoToFrontend(group.Admins),
		"privacy":       group.Privacy,
		"rules":         MapRulesArrayGoToFrontend(group.Rules),
		"featured":      group.Featured,
		"size":          group.Size,
	}
}

// Category Mappers
func MapGroupCategoryFrontendToGo(data map[string]interface{}) models.GroupCategory {
	return models.GroupCategory{
		Name:        getStringValue(data, "name"),
		Icon:        getStringValue(data, "icon"),
		Description: getStringValue(data, "description"),
		Count:       getIntValueSafe(data, "count"),
	}
}

func MapGroupCategoryGoToFirestore(category models.GroupCategory) map[string]interface{} {
	return map[string]interface{}{
		"name":        category.Name,
		"icon":        category.Icon,
		"description": category.Description,
		"count":       category.Count,
	}
}

func MapGroupCategoryFirestoreToGo(data map[string]interface{}) models.GroupCategory {
	return models.GroupCategory{
		Name:        getStringValue(data, "name"),
		Icon:        getStringValue(data, "icon"),
		Description: getStringValue(data, "description"),
		Count:       getIntValueSafe(data, "count"),
	}
}

func MapGroupCategoryGoToFrontend(category models.GroupCategory) map[string]interface{} {
	return map[string]interface{}{
		"name":        category.Name,
		"icon":        category.Icon,
		"description": category.Description,
		"count":       category.Count,
	}
}

// Array Mappers for Related Entities
func MapMembersArrayFrontendToGo(data map[string]interface{}, key string) []models.Member {
	var members []models.Member
	if arr, ok := data[key].([]interface{}); ok {
		for _, item := range arr {
			if memberMap, ok := item.(map[string]interface{}); ok {
				members = append(members, models.Member{
					UserID:   getStringValue(memberMap, "userId"),
					GroupID:  getStringValue(memberMap, "groupId"),
					JoinedAt: getTimeValue(memberMap, "joinedAt"),
					Role:     getStringValue(memberMap, "role"),
					Status:   getStringValue(memberMap, "status"),
				})
			}
		}
	}
	return members
}

// Topics Array Mappers
func MapTopicsArrayFrontendToGo(data map[string]interface{}, key string) []models.Topic {
	var topics []models.Topic
	if arr, ok := data[key].([]interface{}); ok {
		for _, item := range arr {
			if topicMap, ok := item.(map[string]interface{}); ok {
				topics = append(topics, models.Topic{
					ID:          getStringValue(topicMap, "id"),
					Name:        getStringValue(topicMap, "name"),
					Color:       getStringValue(topicMap, "color"),
					Description: getStringValue(topicMap, "description"),
				})
			}
		}
	}
	return topics
}

func MapTopicsArrayGoToFirestore(topics []models.Topic) []map[string]interface{} {
	var result []map[string]interface{}
	for _, topic := range topics {
		result = append(result, map[string]interface{}{
			"id":          topic.ID,
			"name":        topic.Name,
			"color":       topic.Color,
			"description": topic.Description,
		})
	}
	return result
}

func MapTopicsArrayFirestoreToGo(data map[string]interface{}, key string) []models.Topic {
	var topics []models.Topic
	if arr, ok := data[key].([]interface{}); ok {
		for _, item := range arr {
			if topicMap, ok := item.(map[string]interface{}); ok {
				topics = append(topics, models.Topic{
					ID:          getStringValue(topicMap, "id"),
					Name:        getStringValue(topicMap, "name"),
					Color:       getStringValue(topicMap, "color"),
					Description: getStringValue(topicMap, "description"),
				})
			}
		}
	}
	return topics
}

func MapTopicsArrayGoToFrontend(topics []models.Topic) []map[string]interface{} {
	var result []map[string]interface{}
	for _, topic := range topics {
		result = append(result, map[string]interface{}{
			"id":          topic.ID,
			"name":        topic.Name,
			"color":       topic.Color,
			"description": topic.Description,
		})
	}
	return result
}

// Events Array Mappers
func MapEventsArrayFrontendToGo(data map[string]interface{}, key string) []models.Event {
	var events []models.Event
	if arr, ok := data[key].([]interface{}); ok {
		for _, item := range arr {
			if eventMap, ok := item.(map[string]interface{}); ok {
				events = append(events, models.Event{
					ID:          getStringValue(eventMap, "id"),
					GroupID:     getStringValue(eventMap, "groupId"),
					Title:       getStringValue(eventMap, "title"),
					Description: getStringValue(eventMap, "description"),
					StartTime:   getTimeValue(eventMap, "startTime"),
					EndTime:     getTimeValue(eventMap, "endTime"),
					Location:    getStringValue(eventMap, "location"),
					Type:        getStringValue(eventMap, "type"),
					Attendees:   MapAttendeesArrayFrontendToGo(eventMap, "attendees"),
				})
			}
		}
	}
	return events
}

func MapEventsArrayGoToFirestore(events []models.Event) []map[string]interface{} {
	var result []map[string]interface{}
	for _, event := range events {
		result = append(result, map[string]interface{}{
			"id":          event.ID,
			"group_id":    event.GroupID,
			"title":       event.Title,
			"description": event.Description,
			"start_time":  event.StartTime,
			"end_time":    event.EndTime,
			"location":    event.Location,
			"type":        event.Type,
			"attendees":   MapAttendeesArrayGoToFirestore(event.Attendees),
		})
	}
	return result
}

func MapEventsArrayFirestoreToGo(data map[string]interface{}, key string) []models.Event {
	var events []models.Event
	if arr, ok := data[key].([]interface{}); ok {
		for _, item := range arr {
			if eventMap, ok := item.(map[string]interface{}); ok {
				events = append(events, models.Event{
					ID:          getStringValue(eventMap, "id"),
					GroupID:     getStringValue(eventMap, "group_id"),
					Title:       getStringValue(eventMap, "title"),
					Description: getStringValue(eventMap, "description"),
					StartTime:   getTimeValue(eventMap, "start_time"),
					EndTime:     getTimeValue(eventMap, "end_time"),
					Location:    getStringValue(eventMap, "location"),
					Type:        getStringValue(eventMap, "type"),
					Attendees:   MapAttendeesArrayFirestoreToGo(eventMap, "attendees"),
				})
			}
		}
	}
	return events
}

func MapEventsArrayGoToFrontend(events []models.Event) []map[string]interface{} {
	var result []map[string]interface{}
	for _, event := range events {
		result = append(result, map[string]interface{}{
			"id":          event.ID,
			"groupId":     event.GroupID,
			"title":       event.Title,
			"description": event.Description,
			"startTime":   event.StartTime.Format(time.RFC3339),
			"endTime":     event.EndTime.Format(time.RFC3339),
			"location":    event.Location,
			"type":        event.Type,
			"attendees":   MapAttendeesArrayGoToFrontend(event.Attendees),
		})
	}
	return result
}

// Resources Array Mappers
func MapResourcesArrayFrontendToGo(data map[string]interface{}, key string) []models.Resource {
	var resources []models.Resource
	if arr, ok := data[key].([]interface{}); ok {
		for _, item := range arr {
			if resourceMap, ok := item.(map[string]interface{}); ok {
				resources = append(resources, models.Resource{
					ID:          getStringValue(resourceMap, "id"),
					GroupID:     getStringValue(resourceMap, "groupId"),
					Title:       getStringValue(resourceMap, "title"),
					Type:        getStringValue(resourceMap, "type"),
					URL:         getStringValue(resourceMap, "url"),
					Description: getStringValue(resourceMap, "description"),
					AddedBy:     getStringValue(resourceMap, "addedBy"),
					AddedAt:     getTimeValue(resourceMap, "addedAt"),
				})
			}
		}
	}
	return resources
}

func MapResourcesArrayGoToFirestore(resources []models.Resource) []map[string]interface{} {
	var result []map[string]interface{}
	for _, resource := range resources {
		result = append(result, map[string]interface{}{
			"id":          resource.ID,
			"group_id":    resource.GroupID,
			"title":       resource.Title,
			"type":        resource.Type,
			"url":         resource.URL,
			"description": resource.Description,
			"added_by":    resource.AddedBy,
			"added_at":    resource.AddedAt,
		})
	}
	return result
}

func MapResourcesArrayFirestoreToGo(data map[string]interface{}, key string) []models.Resource {
	var resources []models.Resource
	if arr, ok := data[key].([]interface{}); ok {
		for _, item := range arr {
			if resourceMap, ok := item.(map[string]interface{}); ok {
				resources = append(resources, models.Resource{
					ID:          getStringValue(resourceMap, "id"),
					GroupID:     getStringValue(resourceMap, "group_id"),
					Title:       getStringValue(resourceMap, "title"),
					Type:        getStringValue(resourceMap, "type"),
					URL:         getStringValue(resourceMap, "url"),
					Description: getStringValue(resourceMap, "description"),
					AddedBy:     getStringValue(resourceMap, "added_by"),
					AddedAt:     getTimeValue(resourceMap, "added_at"),
				})
			}
		}
	}
	return resources
}

func MapResourcesArrayGoToFrontend(resources []models.Resource) []map[string]interface{} {
	var result []map[string]interface{}
	for _, resource := range resources {
		result = append(result, map[string]interface{}{
			"id":          resource.ID,
			"groupId":     resource.GroupID,
			"title":       resource.Title,
			"type":        resource.Type,
			"url":         resource.URL,
			"description": resource.Description,
			"addedBy":     resource.AddedBy,
			"addedAt":     resource.AddedAt.Format(time.RFC3339),
		})
	}
	return result
}

// Rules Array Mappers
func MapRulesArrayFrontendToGo(data map[string]interface{}, key string) []models.Rule {
	var rules []models.Rule
	if arr, ok := data[key].([]interface{}); ok {
		for _, item := range arr {
			if ruleMap, ok := item.(map[string]interface{}); ok {
				rules = append(rules, models.Rule{
					ID:          getStringValue(ruleMap, "id"),
					GroupID:     getStringValue(ruleMap, "groupId"),
					Title:       getStringValue(ruleMap, "title"),
					Description: getStringValue(ruleMap, "description"),
					Order:       getIntValueSafe(ruleMap, "order"),
				})
			}
		}
	}
	return rules
}

func MapRulesArrayGoToFirestore(rules []models.Rule) []map[string]interface{} {
	var result []map[string]interface{}
	for _, rule := range rules {
		result = append(result, map[string]interface{}{
			"id":          rule.ID,
			"group_id":    rule.GroupID,
			"title":       rule.Title,
			"description": rule.Description,
			"order":       rule.Order,
		})
	}
	return result
}

func MapRulesArrayFirestoreToGo(data map[string]interface{}, key string) []models.Rule {
	var rules []models.Rule
	if arr, ok := data[key].([]interface{}); ok {
		for _, item := range arr {
			if ruleMap, ok := item.(map[string]interface{}); ok {
				rules = append(rules, models.Rule{
					ID:          getStringValue(ruleMap, "id"),
					GroupID:     getStringValue(ruleMap, "group_id"),
					Title:       getStringValue(ruleMap, "title"),
					Description: getStringValue(ruleMap, "description"),
					Order:       getIntValueSafe(ruleMap, "order"),
				})
			}
		}
	}
	return rules
}

func MapRulesArrayGoToFrontend(rules []models.Rule) []map[string]interface{} {
	var result []map[string]interface{}
	for _, rule := range rules {
		result = append(result, map[string]interface{}{
			"id":          rule.ID,
			"groupId":     rule.GroupID,
			"title":       rule.Title,
			"description": rule.Description,
			"order":       rule.Order,
		})
	}
	return result
}

// Admins Array Mappers
func MapAdminsArrayFrontendToGo(data map[string]interface{}, key string) []*models.User {
	var admins []*models.User
	if arr, ok := data[key].([]interface{}); ok {
		for _, item := range arr {
			if adminMap, ok := item.(map[string]interface{}); ok {
				user := MapFrontendToUser(adminMap)
				admins = append(admins, &user)
			}
		}
	}
	return admins
}

func MapAdminsArrayGoToFirestore(admins []*models.User) []map[string]interface{} {
	var result []map[string]interface{}
	for _, admin := range admins {
		result = append(result, MapUserFrontendToBackend(admin))
	}
	return result
}

func MapAdminsArrayFirestoreToGo(data map[string]interface{}, key string) []*models.User {
	var admins []*models.User
	if arr, ok := data[key].([]interface{}); ok {
		for _, item := range arr {
			if adminMap, ok := item.(map[string]interface{}); ok {
				user := MapBackendToUser(adminMap)
				admins = append(admins, &user)
			}
		}
	}
	return admins
}

func MapAdminsArrayGoToFrontend(admins []*models.User) []map[string]interface{} {
	var result []map[string]interface{}
	for _, admin := range admins {
		result = append(result, MapUserToFrontend(admin))
	}
	return result
}

// Attendees Array Mappers
func MapAttendeesArrayFrontendToGo(data map[string]interface{}, key string) []*models.User {
	var attendees []*models.User
	if arr, ok := data[key].([]interface{}); ok {
		for _, item := range arr {
			if attendeeMap, ok := item.(map[string]interface{}); ok {
				user := MapFrontendToUser(attendeeMap)
				attendees = append(attendees, &user)
			}
		}
	}
	return attendees
}

func MapAttendeesArrayGoToFirestore(attendees []*models.User) []map[string]interface{} {
	var result []map[string]interface{}
	for _, attendee := range attendees {
		result = append(result, MapUserFrontendToBackend(attendee))
	}
	return result
}

func MapAttendeesArrayFirestoreToGo(data map[string]interface{}, key string) []*models.User {
	var attendees []*models.User
	if arr, ok := data[key].([]interface{}); ok {
		for _, item := range arr {
			if attendeeMap, ok := item.(map[string]interface{}); ok {
				user := MapBackendToUser(attendeeMap)
				attendees = append(attendees, &user)
			}
		}
	}
	return attendees
}

func MapAttendeesArrayGoToFrontend(attendees []*models.User) []map[string]interface{} {
	var result []map[string]interface{}
	for _, attendee := range attendees {
		result = append(result, MapUserToFrontend(attendee))
	}
	return result
}

// Members Array Mappers
func MapMembersArrayGoToFirestore(members []models.Member) []map[string]interface{} {
	var result []map[string]interface{}
	for _, member := range members {
		result = append(result, map[string]interface{}{
			"user_id":   member.UserID,
			"group_id":  member.GroupID,
			"joined_at": member.JoinedAt,
			"role":      member.Role,
			"status":    member.Status,
		})
	}
	return result
}

func MapMembersArrayFirestoreToGo(data map[string]interface{}, key string) []models.Member {
	var members []models.Member
	if arr, ok := data[key].([]interface{}); ok {
		for _, item := range arr {
			if memberMap, ok := item.(map[string]interface{}); ok {
				members = append(members, models.Member{
					UserID:   getStringValue(memberMap, "user_id"),
					GroupID:  getStringValue(memberMap, "group_id"),
					JoinedAt: getTimeValue(memberMap, "joined_at"),
					Role:     getStringValue(memberMap, "role"),
					Status:   getStringValue(memberMap, "status"),
				})
			}
		}
	}
	return members
}

func MapMembersArrayGoToFrontend(members []models.Member) []map[string]interface{} {
	var result []map[string]interface{}
	for _, member := range members {
		result = append(result, map[string]interface{}{
			"userId":   member.UserID,
			"groupId":  member.GroupID,
			"joinedAt": member.JoinedAt.Format(time.RFC3339),
			"role":     member.Role,
			"status":   member.Status,
		})
	}
	return result
}
