package mappers

import (
	"time"

	"github.com/rogerjeasy/go-letusconnect/models"
)

func MapPollFrontendToGo(data map[string]interface{}) models.Poll {
	options := []models.PollOption{}
	if opts, ok := data["options"].([]interface{}); ok {
		for _, opt := range opts {
			if optionMap, ok := opt.(map[string]interface{}); ok {
				options = append(options, models.PollOption{
					ID:        getStringValue(optionMap, "id"),
					Text:      getStringValue(optionMap, "text"),
					VoteCount: int(getFloatValue(optionMap, "voteCount")),
				})
			}
		}
	}

	return models.Poll{
		ID:                 getStringValue(data, "id"),
		Question:           getStringValue(data, "question"),
		Options:            options,
		CreatedBy:          getStringValue(data, "createdBy"),
		CreatedAt:          getTimeValue(data, "createdAt"),
		ExpiresAt:          getOptionalTimeValue(data, "expiresAt"),
		AllowMultipleVotes: getBoolValue(data, "allowMultipleVotes"),
		Votes:              getVotesMap(data, "votes"),
		IsClosed:           getBoolValue(data, "isClosed"),
	}
}

func MapPollGoToFirestore(poll models.Poll) map[string]interface{} {
	options := []map[string]interface{}{}
	for _, opt := range poll.Options {
		options = append(options, map[string]interface{}{
			"id":         opt.ID,
			"text":       opt.Text,
			"vote_count": opt.VoteCount,
		})
	}

	return map[string]interface{}{
		"id":                   poll.ID,
		"question":             poll.Question,
		"options":              options,
		"created_by":           poll.CreatedBy,
		"created_at":           poll.CreatedAt,
		"expires_at":           poll.ExpiresAt,
		"allow_multiple_votes": poll.AllowMultipleVotes,
		"votes":                poll.Votes,
		"is_closed":            poll.IsClosed,
	}
}

func MapPollFirestoreToFrontend(data map[string]interface{}) map[string]interface{} {
	options := []map[string]interface{}{}
	if opts, ok := data["options"].([]interface{}); ok {
		for _, opt := range opts {
			if optionMap, ok := opt.(map[string]interface{}); ok {
				options = append(options, map[string]interface{}{
					"id":        getStringValue(optionMap, "id"),
					"text":      getStringValue(optionMap, "text"),
					"voteCount": int(getFloatValue(optionMap, "vote_count")),
				})
			}
		}
	}

	return map[string]interface{}{
		"id":                 getStringValue(data, "id"),
		"question":           getStringValue(data, "question"),
		"options":            options,
		"createdBy":          getStringValue(data, "created_by"),
		"createdAt":          getTimeValue(data, "created_at").Format(time.RFC3339),
		"expiresAt":          dereferenceTime(getOptionalTimeValue(data, "expires_at")),
		"allowMultipleVotes": getBoolValue(data, "allow_multiple_votes"),
		"votes":              getVotesMap(data, "votes"),
		"isClosed":           getBoolValue(data, "is_closed"),
	}
}

func MapPollFirestoreToGo(data map[string]interface{}) models.Poll {
	options := []models.PollOption{}
	if opts, ok := data["options"].([]interface{}); ok {
		for _, opt := range opts {
			if optionMap, ok := opt.(map[string]interface{}); ok {
				options = append(options, models.PollOption{
					ID:        getStringValue(optionMap, "id"),
					Text:      getStringValue(optionMap, "text"),
					VoteCount: int(getFloatValue(optionMap, "vote_count")),
				})
			}
		}
	}

	return models.Poll{
		ID:                 getStringValue(data, "id"),
		Question:           getStringValue(data, "question"),
		Options:            options,
		CreatedBy:          getStringValue(data, "created_by"),
		CreatedAt:          getTimeValue(data, "created_at"),
		ExpiresAt:          getOptionalTimeValue(data, "expires_at"),
		AllowMultipleVotes: getBoolValue(data, "allow_multiple_votes"),
		Votes:              getVotesMap(data, "votes"),
		IsClosed:           getBoolValue(data, "is_closed"),
	}
}

func getVotesMap(data map[string]interface{}, key string) map[string][]string {
	votes := make(map[string][]string)
	if votesData, ok := data[key].(map[string]interface{}); ok {
		for userID, options := range votesData {
			if optionIDs, ok := options.([]interface{}); ok {
				for _, optionID := range optionIDs {
					votes[userID] = append(votes[userID], optionID.(string))
				}
			}
		}
	}
	return votes
}
