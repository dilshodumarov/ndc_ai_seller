package category

import (
	"context"
	"sugurta/internal/entity"
	"sugurta/internal/infrastructure/repository"
	"time"
)

type categoryService struct {
	ctxTimeout   time.Duration
	categoryRepo repository.Category
}

func NewCategoryService(ctxTimeout time.Duration, categoryRepo repository.Category) Category {

	return &categoryService{
		ctxTimeout:   ctxTimeout,
		categoryRepo: categoryRepo,
	}
}

// CreateCategory -.
func (c *categoryService) Create(ctx context.Context, category *entity.CreateCategoryRequest) error {
	ctx, cancel := context.WithTimeout(ctx, c.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	err := c.categoryRepo.Create(ctx, category)
	if err != nil {
		return err
	}

	return nil
}

// GetCategory -.
func (c *categoryService) Get(ctx context.Context, id string) (*entity.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, c.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	category, err := c.categoryRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (c *categoryService) List(ctx context.Context, filter entity.CategoryFilter) (*entity.GetAllCategoriesResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	category, err := c.categoryRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// UpdateCategory -.
func (c *categoryService) Update(ctx context.Context, category *entity.UpdateCategoryRequest) error {
	ctx, cancel := context.WithTimeout(ctx, c.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	err := c.categoryRepo.Update(ctx, category)
	if err != nil {
		return err
	}

	return nil
}

// DeleteCategory -.
func (c *categoryService) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, c.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	err := c.categoryRepo.Delete(ctx, id)
	if err != nil {
		// c.log.Error("categoryService - DeleteCategory - s.prodct.DeleteCategory: ", err.Error())
		return err
	}

	return nil
}
