package order

import (
	"context"
	"sugurta/internal/entity"
	"sugurta/internal/infrastructure/repository"
	"time"
)

type orderService struct {
	ctxTimeout   time.Duration
	orderService repository.Order
}

func NewRoleService(ctxTimeout time.Duration, order repository.Order) Order {
	return &orderService{
		ctxTimeout:   ctxTimeout,
		orderService: order,
	}
}

// CreateOrder -.
func (o *orderService) Create(ctx context.Context, order *entity.CreateOrderRequest) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, o.ctxTimeout)
	defer cancel()

	id, err := o.orderService.Create(ctx, order)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (o *orderService) Get(ctx context.Context, id string) (*entity.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, o.ctxTimeout)
	defer cancel()

	order, err := o.orderService.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (o *orderService) List(ctx context.Context, filter *entity.OrderFilter, limit, offset uint64) (*entity.GetAllOrdersResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, o.ctxTimeout)
	defer cancel()

	orders, err := o.orderService.List(ctx, filter, limit, offset)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *orderService) Update(ctx context.Context, order *entity.OrderUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, o.ctxTimeout)
	defer cancel()

	return o.orderService.Update(ctx, order)
}

func (o *orderService) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, o.ctxTimeout)
	defer cancel()

	return o.orderService.Delete(ctx, id)
}
