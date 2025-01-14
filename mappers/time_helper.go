package mappers

import (
	"time"
)

func getFirestoreTimeToGoTime(value interface{}) time.Time {
	if value == nil {
		return time.Time{}
	}

	if t, ok := value.(time.Time); ok {
		return t
	}

	if tsMap, ok := value.(map[string]interface{}); ok {
		if seconds, ok := tsMap["seconds"].(int64); ok {
			if nanos, ok := tsMap["nanoseconds"].(int64); ok {
				return time.Unix(seconds, nanos)
			}
		}
	}

	return time.Time{}
}
