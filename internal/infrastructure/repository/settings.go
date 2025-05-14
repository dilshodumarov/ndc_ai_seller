package repository

import (
	"context"
	"sugurta/internal/entity"
)

type SettingsStorage interface {
	Create(ctx context.Context, req *entity.CreateOrderStatusRequest) error
	Get(ctx context.Context, guid string) (*entity.OrderStatus, error)
	Update(ctx context.Context, req *entity.UpdateOrderStatusRequest) error
	Delete(ctx context.Context, guid string) error
	List(ctx context.Context, businessID string) ([]*entity.OrderStatus, error)
}
