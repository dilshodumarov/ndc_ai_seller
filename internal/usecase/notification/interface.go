package notification

import (
	"context"
	"sugurta/internal/entity"
)

type Notification interface {
	Create(ctx context.Context, req *entity.CreateNotificationRequest) (string, error)
	Get(ctx context.Context, id string) (*entity.GetNotification, error)
	Update(ctx context.Context, req *entity.UpdateNotificationRequest) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, req entity.ListNotificationRequest) (*entity.ListNotification, error)
	MarkAsRead(ctx context.Context, id string) error
}
