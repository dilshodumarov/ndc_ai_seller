package repository

import (
	"sugurta/internal/entity"
	"context"
)

type Order interface {
	Create(ctx context.Context, order *entity.CreateOrderRequest) (string, error)
	Get(ctx context.Context, id string) (*entity.Order, error)
	List(ctx context.Context, filter *entity.OrderFilter, limit, offset uint64) (*entity.GetAllOrdersResponse, error)
	Update(ctx context.Context, order *entity.OrderUpdate) error
	Delete(ctx context.Context, id string) error
	// OrderProducts
	// CreateOrderProducts(ctx context.Context, op entity.OrderProducts) error
	// GetOrderProducts(ctx context.Context, id string) (*entity.OrderProducts, error)
	// UpdateOrderProducts(ctx context.Context, op entity.OrderProducts) error
	// DeleteOrderProducts(ctx context.Context, id string) error
}
