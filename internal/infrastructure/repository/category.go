package repository

import (
	"sugurta/internal/entity"
	"context"
)

type Category interface {
	Create(ctx context.Context, data *entity.CreateCategoryRequest) error
	Get(ctx context.Context, id string) (*entity.Category, error)
	List(ctx context.Context, filter entity.CategoryFilter) (*entity.GetAllCategoriesResponse, error)
	Update(ctx context.Context, category *entity.UpdateCategoryRequest) error
	Delete(ctx context.Context, id string) error
}
