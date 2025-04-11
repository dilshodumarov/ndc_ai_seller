package service

// import (
// 	"context"

// 	v1 "sugurta/genproto/product_service"
// 	"sugurta/internal/entity"
// 	"sugurta/internal/usecase"

// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// )

// // ProductService is implementation of ProductServiceServer
// type ProductService struct {
// 	v1.UnimplementedProductServiceServer
// 	useCase usecase.UseCase
// }

// // NewProductService creates a new ProductService
// func NewProductService(useCase usecase.UseCase) *ProductService {
// 	return &ProductService{
// 		useCase: useCase,
// 	}
// }

// // CreateProduct creates a new product
// func (s *ProductService) CreateProduct(ctx context.Context, req *v1.CreateProductRequest) (*v1.Empty, error) {
// 	if req.Product == nil {
// 		return nil, status.Error(codes.InvalidArgument, "product is required")
// 	}

// 	err := s.useCase.CreateProduct(ctx, entity.CreateProductRequest{
// 		Name:        req.Product.Name,
// 		Description: req.Product.Description,
// 		Cost:        int(req.Product.Price),
// 		CategoryID:  req.Product.CategoryId,
// 	})
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &v1.Empty{}, nil
// }

// // GetProduct gets a product by ID
// func (s *ProductService) GetProduct(ctx context.Context, req *v1.GetProductRequest) (*v1.Product, error) {
// 	if req.Id == "" {
// 		return nil, status.Error(codes.InvalidArgument, "id is required")
// 	}

// 	product, err := s.useCase.GetProduct(ctx, req.Id)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	result := &v1.Product{
// 		Id:          product.ID,
// 		Name:        product.Name,
// 		Description: product.Description,
// 		Price:       float64(product.Cost),
// 		CategoryId:  product.CategoryID,
// 		CreatedAt    time.Time:   product.CreatedAt    time.Time,
// 		UpdatedAt    time.Time:   product.UpdatedAt    time.Time,
// 	}

// 	// if product. != "" {
// 	// 	result.Category = &v1.Category{
// 	// 		Id:          product.Category.ID,
// 	// 		Name:        product.Category.Name,
// 	// 		Description: product.Category.Description,
// 	// 		CreatedAt    time.Time:   product.Category.CreatedAt    time.Time,
// 	// 		UpdatedAt    time.Time:   product.Category.UpdatedAt    time.Time,
// 	// 	}
// 	// }

// 	return result, nil
// }

// // UpdateProduct updates a product
// func (s *ProductService) UpdateProduct(ctx context.Context, req *v1.UpdateProductRequest) (*v1.Empty, error) {
// 	if req.Product == nil {
// 		return nil, status.Error(codes.InvalidArgument, "product is required")
// 	}

// 	err := s.useCase.UpdateProduct(ctx, entity.UpdateProductRequest{
// 		ID:          req.Product.Id,
// 		Name:        req.Product.Name,
// 		Description: req.Product.Description,
// 		Cost:        int(req.Product.Price),
// 		CategoryID:  req.Product.CategoryId,
// 	})
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &v1.Empty{}, nil
// }

// // DeleteProduct deletes a product
// func (s *ProductService) DeleteProduct(ctx context.Context, req *v1.DeleteProductRequest) (*v1.Empty, error) {
// 	if req.Id == "" {
// 		return nil, status.Error(codes.InvalidArgument, "id is required")
// 	}

// 	err := s.useCase.DeleteProduct(ctx, req.Id)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &v1.Empty{}, nil
// }
