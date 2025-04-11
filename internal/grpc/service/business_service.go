package service

// import (
// 	"context"

// 	v1 "sugurta/genproto/business_service"
// 	"sugurta/internal/entity"
// 	"sugurta/internal/usecase"

// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// )

// // BusinessService is implementation of BusinessServiceServer
// type BusinessService struct {
// 	v1.UnimplementedBusinessServiceServer
// 	useCase usecase.UseCase
// }

// // NewBusinessService creates a new BusinessService
// func NewBusinessService(useCase usecase.UseCase) *BusinessService {
// 	return &BusinessService{
// 		useCase: useCase,
// 	}
// }

// // CreateBusiness creates a new business
// func (s *BusinessService) CreateBusiness(ctx context.Context, req *v1.CreateBusinessRequest) (*v1.Empty, error) {
// 	if req.Business == nil {
// 		return nil, status.Error(codes.InvalidArgument, "business is required")
// 	}

// 	err := s.useCase.CreateBusiness(ctx, entity.CreateBusinessRequest{
// 		Name:        req.Business.Name,
// 		Description: req.Business.Description,
// 	})
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &v1.Empty{}, nil
// }

// // GetBusiness gets a business by ID
// func (s *BusinessService) GetBusiness(ctx context.Context, req *v1.GetBusinessRequest) (*v1.Business, error) {
// 	if req.Id == "" {
// 		return nil, status.Error(codes.InvalidArgument, "id is required")
// 	}

// 	business, err := s.useCase.GetBusiness(ctx, req.Id)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &v1.Business{
// 		Id:          business.ID,
// 		Name:        business.Name,
// 		Description: business.Description,
// 		CreatedAt:   business.CreatedAt,
// 		UpdatedAt:   business.UpdatedAt,
// 	}, nil
// }

// // GetAllBusinesses gets all businesses with pagination
// func (s *BusinessService) GetAllBusinesses(ctx context.Context, req *v1.GetAllBusinessesRequest) (*v1.GetAllBusinessesResponse, error) {
// 	limit := int(req.Limit)
// 	offset := int(req.Offset)
// 	if limit <= 0 {
// 		limit = 10
// 	}
// 	if offset < 0 {
// 		offset = 0
// 	}

// 	businesses, err := s.useCase.GetAllBusinesses(ctx, entity.GetAllBusinessesRequest{
// 		Limit: limit,
// 		Page:  offset,
// 	})
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	response := &v1.GetAllBusinessesResponse{
// 		Total:      int32(businesses.Count),
// 		Businesses: make([]*v1.Business, 0, len(businesses.Businesses)),
// 	}

// 	for _, business := range businesses.Businesses {
// 		response.Businesses = append(response.Businesses, &v1.Business{
// 			Id:          business.ID,
// 			Name:        business.Name,
// 			Description: business.Description,
// 			CreatedAt:   business.CreatedAt,
// 			UpdatedAt:   business.UpdatedAt,
// 		})
// 	}

// 	return response, nil
// }

// // UpdateBusiness updates a business
// func (s *BusinessService) UpdateBusiness(ctx context.Context, req *v1.UpdateBusinessRequest) (*v1.Empty, error) {
// 	if req.Business == nil {
// 		return nil, status.Error(codes.InvalidArgument, "business is required")
// 	}

// 	err := s.useCase.UpdateBusiness(ctx, entity.UpdateBusinessRequest{
// 		ID:          req.Business.Id,
// 		Name:        req.Business.Name,
// 		Description: req.Business.Description,
// 	})
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &v1.Empty{}, nil
// }

// // DeleteBusiness deletes a business
// func (s *BusinessService) DeleteBusiness(ctx context.Context, req *v1.DeleteBusinessRequest) (*v1.Empty, error) {
// 	if req.Id == "" {
// 		return nil, status.Error(codes.InvalidArgument, "id is required")
// 	}

// 	err := s.useCase.DeleteBusiness(ctx, req.Id, req.UserId)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &v1.Empty{}, nil
// }
