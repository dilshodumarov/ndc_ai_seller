package repository

import (
	"sugurta/internal/entity"
	"context"
)

type Product interface {
	Create(ctx context.Context, product *entity.CreateProductRequest) error
	Get(ctx context.Context,id string) (*entity.Product, error)
	List(ctx context.Context, limit, offset uint64, filter map[string]string) (*entity.GetAllProductsResponse, error)
	Update(ctx context.Context, product *entity.UpdateProductRequest) error
	Delete(ctx context.Context, id string) error
}
