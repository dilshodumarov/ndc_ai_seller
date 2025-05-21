package order

import (
	"sugurta/internal/entity"
	"context"
)

type Order interface {
	Create(ctx context.Context, order *entity.CreateOrderRequest) (string, error)
	Get(ctx context.Context,id string) (*entity.Order, error)
	List(ctx context.Context, filter *entity.OrderFilter, limit, offset uint64) (*entity.GetAllOrdersResponse, error)
	Update(ctx context.Context, order *entity.OrderUpdate) error
	Delete(ctx context.Context, id string) error
	GetProductsByOrderID(ctx context.Context, orderID string) ([]entity.OrderProductBuOrderID, error)
}
