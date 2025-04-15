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

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	id, err := o.orderService.Create(ctx, order)
	if err != nil {
		return "", err
	}

	return id, nil
}

// GetOrder -.
func (o *orderService) Get(ctx context.Context, params map[string]string) (*entity.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, o.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	order, err := o.orderService.Get(ctx, params)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// ListOrders -.
func (o *orderService) List(ctx context.Context, limit, offset uint64, filter map[string]string) (*entity.GetAllOrdersResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, o.ctxTimeout)
	defer cancel()
	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	orders, err := o.orderService.List(ctx, limit, offset, filter)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

// UpdateOrder -.
func (o *orderService) Update(ctx context.Context, order *entity.Order) error {
	ctx, cancel := context.WithTimeout(ctx, o.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	return o.orderService.Update(ctx, order)
}

// DeleteOrder -.
func (o *orderService) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, o.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	return o.orderService.Delete(ctx, id)
}
