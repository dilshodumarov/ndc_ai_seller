package notification

import (
	"sugurta/internal/entity"
	"sugurta/internal/infrastructure/repository"
	"context"
	"time"
)

// NotificationService is the notification use case implementation
type notificationService struct {
	ctxTimeout    time.Duration
	notificationRepo repository.Notification
}

func NewNotificationService(ctxTimeout time.Duration, n repository.Notification) *notificationService {
	return &notificationService{
		ctxTimeout:       ctxTimeout,
		notificationRepo: n,
	}
}

func (ns *notificationService) Create(ctx context.Context, req *entity.CreateNotificationRequest) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, ns.ctxTimeout)
	defer cancel()

	// tracing can be added here if needed

	return ns.notificationRepo.Create(ctx, req)
}

func (ns *notificationService) Get(ctx context.Context, id string) (*entity.GetNotification, error) {
	ctx, cancel := context.WithTimeout(ctx, ns.ctxTimeout)
	defer cancel()

	// tracing can be added here if needed

	return ns.notificationRepo.Get(ctx, id)
}

func (ns *notificationService) Update(ctx context.Context, req *entity.UpdateNotificationRequest) error {
	ctx, cancel := context.WithTimeout(ctx, ns.ctxTimeout)
	defer cancel()

	// tracing can be added here if needed

	return ns.notificationRepo.Update(ctx, req)
}

func (ns *notificationService) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, ns.ctxTimeout)
	defer cancel()

	// tracing can be added here if needed

	return ns.notificationRepo.Delete(ctx, id)
}

func (ns *notificationService) List(ctx context.Context, req entity.ListNotificationRequest) (*entity.ListNotification, error) {
	ctx, cancel := context.WithTimeout(ctx, ns.ctxTimeout)
	defer cancel()

	// tracing can be added here if needed

	return ns.notificationRepo.List(ctx, req)
}

func (ns *notificationService) MarkAsRead(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, ns.ctxTimeout)
	defer cancel()

	// tracing can be added here if needed

	return ns.notificationRepo.MarkAsRead(ctx, id)
}
