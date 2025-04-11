package order

import (
	"sugurta/internal/entity"
	"context"
)

type Order interface {
	Create(ctx context.Context, order *entity.CreateOrderRequest) (string, error)
	Get(ctx context.Context, params map[string]string) (*entity.Order, error)
	List(ctx context.Context, limit, offset uint64, filter map[string]string) (*entity.GetAllOrdersResponse, error)
	Update(ctx context.Context, order *entity.Order) error
	Delete(ctx context.Context, id string) error
}
