package category

import (
	"sugurta/internal/entity"
	"context"
)

type Category interface {
	Create(ctx context.Context, data *entity.CreateCategoryRequest) error
	Get(ctx context.Context, id string) (*entity.Category, error)
	List(ctx context.Context, limit, offset uint64, filter map[string]string) (*entity.GetAllCategoriesResponse, error)
	Update(ctx context.Context, role *entity.UpdateCategoryRequest) error
	Delete(ctx context.Context, id string) error
}
