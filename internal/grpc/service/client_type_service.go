package service

// import (
// 	"context"

// 	v1 "sugurta/genproto/auth_service"
// 	"sugurta/internal/entity"
// 	"sugurta/internal/usecase"

// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// )

// // ClientTypeService is implementation of ClientTypeServiceServer
// type ClientTypeService struct {
// 	v1.UnimplementedClientTypeServiceServer
// 	useCase usecase.UseCase
// }

// // NewClientTypeService creates a new ClientTypeService
// func NewClientTypeService(useCase usecase.UseCase) *ClientTypeService {
// 	return &ClientTypeService{
// 		useCase: useCase,
// 	}
// }

// // CreateClientType creates a new client type
// func (s *ClientTypeService) CreateClientType(ctx context.Context, req *v1.CreateClientTypeRequest) (*v1.Empty, error) {
// 	if req.ClientType == nil {
// 		return nil, status.Error(codes.InvalidArgument, "client_type is required")
// 	}

// 	err := s.useCase.CreateClientType(ctx, entity.CreateClientTypeRequest{
// 		Name: req.ClientType.Name,
// 	})
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &v1.Empty{}, nil
// }

// // GetClientTypeById gets a client type by ID
// func (s *ClientTypeService) GetClientTypeById(ctx context.Context, req *v1.GetClientTypeByIdRequest) (*v1.ClientType, error) {
// 	if req.Id == "" {
// 		return nil, status.Error(codes.InvalidArgument, "id is required")
// 	}

// 	clientType, err := s.useCase.GetClientTypeById(ctx, req.Id)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &v1.ClientType{
// 		Id:        clientType.ID,
// 		Name:      clientType.Name,
// 		CreatedAt    time.Time: clientType.CreatedAt    time.Time.String(),
// 		UpdatedAt    time.Time: clientType.UpdatedAt    time.Time.String(),
// 	}, nil
// }

// // GetClientTypeByName gets a client type by name
// func (s *ClientTypeService) GetClientTypeByName(ctx context.Context, req *v1.GetClientTypeByNameRequest) (*v1.ClientType, error) {
// 	if req.Name == "" {
// 		return nil, status.Error(codes.InvalidArgument, "name is required")
// 	}

// 	clientType, err := s.useCase.GetClientTypeByName(ctx, req.Name)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &v1.ClientType{
// 		Id:        clientType.ID,
// 		Name:      clientType.Name,
// 		CreatedAt    time.Time: clientType.CreatedAt    time.Time.String(),
// 		UpdatedAt    time.Time: clientType.UpdatedAt    time.Time.String(),
// 	}, nil
// }

// // GetClientTypes gets all client types with pagination
// func (s *ClientTypeService) GetClientTypes(ctx context.Context, req *v1.GetClientTypesRequest) (*v1.ClientTypeListResponse, error) {
// 	limit := int(req.Limit)
// 	offset := int(req.Offset)
// 	if limit <= 0 {
// 		limit = 10
// 	}
// 	if offset < 0 {
// 		offset = 0
// 	}

// 	clientTypes, err := s.useCase.GetClientTypes(ctx, offset, limit)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	response := &v1.ClientTypeListResponse{
// 		Total:       int32(clientTypes.Count),
// 		ClientTypes: make([]*v1.ClientType, 0, len(clientTypes.Items)),
// 	}

// 	for _, clientType := range clientTypes.Items {
// 		response.ClientTypes = append(response.ClientTypes, &v1.ClientType{
// 			Id:        clientType.ID,
// 			Name:      clientType.Name,
// 			CreatedAt    time.Time: clientType.CreatedAt    time.Time.String(),
// 			UpdatedAt    time.Time: clientType.UpdatedAt    time.Time.String(),
// 		})
// 	}

// 	return response, nil
// }

// // UpdateClientType updates a client type
// func (s *ClientTypeService) UpdateClientType(ctx context.Context, req *v1.UpdateClientTypeRequest) (*v1.Empty, error) {
// 	if req.ClientType == nil {
// 		return nil, status.Error(codes.InvalidArgument, "client_type is required")
// 	}

// 	err := s.useCase.UpdateClientType(ctx, entity.UpdateClientType{
// 		ID:   req.ClientType.Id,
// 		Name: req.ClientType.Name,
// 	})
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &v1.Empty{}, nil
// }

// // DeleteClientType deletes a client type
// func (s *ClientTypeService) DeleteClientType(ctx context.Context, req *v1.DeleteClientTypeRequest) (*v1.Empty, error) {
// 	if req.Id == "" {
// 		return nil, status.Error(codes.InvalidArgument, "id is required")
// 	}

// 	err := s.useCase.DeleteClientType(ctx, req.Id)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &v1.Empty{}, nil
// }
