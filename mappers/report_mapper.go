package mappers

import "github.com/rogerjeasy/go-letusconnect/models"

func MapReportFirestoreToGo(data map[string]interface{}) models.Report {
	return models.Report{
		ID:          getStringValue(data, "id"),
		ReporterID:  getStringValue(data, "reporter_id"),
		MessageID:   getStringValue(data, "message_id"),
		GroupChatID: getStringValue(data, "group_chat_id"),
		Reason:      getStringValue(data, "reason"),
		Description: getStringValue(data, "description"),
		Status:      getStringValue(data, "status"),
		CreatedAt:   getTimeValue(data, "created_at"),
		UpdatedAt:   getTimeValue(data, "updated_at"),
		ResolvedBy:  getOptionalStringPointer(data, "resolved_by"),
	}
}

func MapReportGoToFirestore(report models.Report) map[string]interface{} {
	return map[string]interface{}{
		"id":            report.ID,
		"reporter_id":   report.ReporterID,
		"message_id":    report.MessageID,
		"group_chat_id": report.GroupChatID,
		"reason":        report.Reason,
		"description":   report.Description,
		"status":        report.Status,
		"created_at":    report.CreatedAt,
		"updated_at":    report.UpdatedAt,
		"resolved_by":   report.ResolvedBy,
	}
}

func MapReportsArrayFromFirestore(data []interface{}) []models.Report {
	reports := []models.Report{}
	for _, item := range data {
		if reportMap, ok := item.(map[string]interface{}); ok {
			reports = append(reports, MapReportFirestoreToGo(reportMap))
		}
	}
	return reports
}

func MapReportsArrayToFirestore(reports []models.Report) []map[string]interface{} {
	result := []map[string]interface{}{}
	for _, report := range reports {
		result = append(result, MapReportGoToFirestore(report))
	}
	return result
}

func getOptionalStringPointer(data map[string]interface{}, key string) *string {
	if value, ok := data[key].(string); ok {
		return &value
	}
	return nil
}
