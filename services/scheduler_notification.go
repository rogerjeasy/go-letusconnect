package services

import (
	"context"

	"github.com/rogerjeasy/go-letusconnect/models"
)

type SchedulerNotificationService struct {
	scheduler *NotificationScheduler
}

func NewSchedulerNotificationService(scheduler *NotificationScheduler) *SchedulerNotificationService {
	return &SchedulerNotificationService{
		scheduler: scheduler,
	}
}

func (s *SchedulerNotificationService) ScheduleNotification(ctx context.Context, req *models.NotificationRequest) error {
	return s.scheduler.ScheduleNotification(ctx, req)
}

func (s *SchedulerNotificationService) CancelNotification(ctx context.Context, notificationID string) error {
	return s.scheduler.CancelNotification(ctx, notificationID)
}

func (s *SchedulerNotificationService) GetNotifications(ctx context.Context, userID string) ([]models.Notification, error) {
	return s.scheduler.GetNotifications(ctx, userID)
}
