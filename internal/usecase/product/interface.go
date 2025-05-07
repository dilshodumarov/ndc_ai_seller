package product

import (
	"sugurta/internal/entity"
	"context"
)

type Product interface {
	Create(ctx context.Context, product *entity.CreateProductRequest) (string,error)
	Get(ctx context.Context, id string) (*entity.Product, error)
	List(ctx context.Context, filter entity.ProductFilter) (*entity.GetAllProductsResponse, error)
	Update(ctx context.Context, product *entity.UpdateProductRequest) error
	Delete(ctx context.Context, id string) error
	AddPicture(ctx context.Context, image *entity.CreateProductImage) (string, error)
	DeletePicture(ctx context.Context, id string) (error)
}
