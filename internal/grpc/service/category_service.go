package service

// import (
// 	"context"

// 	v1 "sugurta/genproto/product_service"
// 	"sugurta/internal/entity"
// 	"sugurta/internal/usecase"

// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// )

// // CategoryService is implementation of CategoryServiceServer
// type CategoryService struct {
// 	v1.UnimplementedCategoryServiceServer
// 	useCase usecase.UseCase
// }

// // NewCategoryService creates a new CategoryService
// func NewCategoryService(useCase usecase.UseCase) *CategoryService {
// 	return &CategoryService{
// 		useCase: useCase,
// 	}
// }

// // CreateCategory creates a new category
// func (s *CategoryService) CreateCategory(ctx context.Context, req *v1.CreateCategoryRequest) (*v1.Empty, error) {
// 	if req.Category == nil {
// 		return nil, status.Error(codes.InvalidArgument, "category is required")
// 	}

// 	err := s.useCase.CreateCategory(ctx, entity.CreateCategoryRequest{
// 		Name: req.Category.Name,
// 	})
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &v1.Empty{}, nil
// }

// // GetCategory gets a category by ID
// func (s *CategoryService) GetCategory(ctx context.Context, req *v1.GetCategoryRequest) (*v1.Category, error) {
// 	if req.Id == "" {
// 		return nil, status.Error(codes.InvalidArgument, "id is required")
// 	}

// 	category, err := s.useCase.GetCategory(ctx, req.Id)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &v1.Category{
// 		Id:        category.ID,
// 		Name:      category.Name,
// 		CreatedAt    time.Time: category.CreatedAt    time.Time,
// 		UpdatedAt    time.Time: category.UpdatedAt    time.Time,
// 	}, nil
// }

// // UpdateCategory updates a category
// func (s *CategoryService) UpdateCategory(ctx context.Context, req *v1.UpdateCategoryRequest) (*v1.Empty, error) {
// 	if req.Category == nil {
// 		return nil, status.Error(codes.InvalidArgument, "category is required")
// 	}

// 	err := s.useCase.UpdateCategory(ctx, entity.Category{
// 		ID:   req.Category.Id,
// 		Name: req.Category.Name,
// 	})
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &v1.Empty{}, nil
// }

// // DeleteCategory deletes a category
// func (s *CategoryService) DeleteCategory(ctx context.Context, req *v1.DeleteCategoryRequest) (*v1.Empty, error) {
// 	if req.Id == "" {
// 		return nil, status.Error(codes.InvalidArgument, "id is required")
// 	}

// 	err := s.useCase.DeleteCategory(ctx, req.Id)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &v1.Empty{}, nil
// }
