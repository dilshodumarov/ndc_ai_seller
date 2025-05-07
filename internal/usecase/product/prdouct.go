package product

import (
	"context"
	"time"

	"sugurta/internal/entity"
	"sugurta/internal/infrastructure/repository"
)

// UseCase -.
type productService struct {
	ctxTimeout  time.Duration
	productRepo repository.Product
}

func NewProductService(ctxTimeout time.Duration, productRepo repository.Product) Product {
	return &productService{
		ctxTimeout:  ctxTimeout,
		productRepo: productRepo,
	}
}

// CreateProduct -.
func (p *productService) Create(ctx context.Context, product *entity.CreateProductRequest) (string,error) {
	ctx, cancel := context.WithTimeout(ctx, p.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	return p.productRepo.Create(ctx, product)
}

// GetProduct -.
func (p *productService) Get(ctx context.Context, id string) (*entity.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, p.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	product, err := p.productRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *productService) List(ctx context.Context, filter entity.ProductFilter) (*entity.GetAllProductsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, p.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	products, err := p.productRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	return products, nil
}

// UpdateProduct -.
func (p *productService) Update(ctx context.Context, product *entity.UpdateProductRequest) error {
	ctx, cancel := context.WithTimeout(ctx, p.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	return p.productRepo.Update(ctx, product)
}

// DeleteProduct -.
func (p *productService) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, p.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	return p.productRepo.Delete(ctx, id)
}


func (p *productService) AddPicture(ctx context.Context, image *entity.CreateProductImage) (string, error){
	ctx, cancel := context.WithTimeout(ctx, p.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	return p.productRepo.AddPicture(ctx, image)
}


func (p *productService) DeletePicture(ctx context.Context, id string) (error){
	ctx, cancel := context.WithTimeout(ctx, p.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	return p.productRepo.DeletePicture(ctx, id)
}


// // CreateAttribute -.
// func (p *ProductUseCase) CreateAttribute(ctx context.Context, a entity.Attribute) error {
// 	err := p.productRepo.CreateAttribute(ctx, a)
// 	if err != nil {
// 		return fmt.Errorf("ProductUseCase - CreateAttribute - s.product.CreateAttribute: %w", err)
// 	}

// 	return nil
// }

// // GetAttribute -.
// func (p *ProductUseCase) GetAttribute(ctx context.Context, id string) (entity.Attribute, error) {
// 	attribute, err := p.productRepo.GetAttribute(ctx, id)
// 	if err != nil {
// 		return entity.Attribute{}, fmt.Errorf("ProductUseCase - GetAttribute - s.product.GetAttribute: %w", err)
// 	}

// 	return attribute, nil
// }

// // UpdateAttribute -.
// func (p *ProductUseCase) UpdateAttribute(ctx context.Context, a entity.Attribute) error {
// 	err := p.productRepo.UpdateAttribute(ctx, a)
// 	if err != nil {
// 		return fmt.Errorf("ProductUseCase - UpdateAttribute - s.product.UpdateAttribute: %w", err)
// 	}

// 	return nil
// }

// // DeleteAttribute -.
// func (p *ProductUseCase) DeleteAttribute(ctx context.Context, id string) error {
// 	err := p.productRepo.DeleteAttribute(ctx, id)
// 	if err != nil {
// 		return fmt.Errorf("ProductUseCase - DeleteAttribute - s.product.DeleteAttribute: %w", err)
// 	}

// 	return nil
// }
