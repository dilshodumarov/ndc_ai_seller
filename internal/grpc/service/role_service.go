package service

// import (
// 	"context"

// 	v1 "sugurta/genproto/auth_service"
// 	"sugurta/internal/entity"
// 	"sugurta/internal/usecase"

// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// )

// // RoleService is implementation of RoleServiceServer
// type RoleService struct {
// 	v1.UnimplementedRoleServiceServer
// 	useCase usecase.UseCase
// }

// // NewRoleService creates a new RoleService
// func NewRoleService(useCase usecase.UseCase) *RoleService {
// 	return &RoleService{
// 		useCase: useCase,
// 	}
// }

// // CreateRole creates a new role
// func (s *RoleService) CreateRole(ctx context.Context, req *v1.CreateRoleRequest) (*v1.Empty, error) {
// 	if req.Role == nil {
// 		return nil, status.Error(codes.InvalidArgument, "role is required")
// 	}

// 	err := s.useCase.CreateRole(ctx, entity.CreateRoleRequest{
// 		Name: req.Role.Name,
// 	})
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &v1.Empty{}, nil
// }

// // GetRoleByID gets a role by ID
// func (s *RoleService) GetRoleByID(ctx context.Context, req *v1.GetRoleByIDRequest) (*v1.Role, error) {
// 	if req.Id == "" {
// 		return nil, status.Error(codes.InvalidArgument, "id is required")
// 	}

// 	role, err := s.useCase.GetRoleByID(ctx, req.Id)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &v1.Role{
// 		Id:   role.ID,
// 		Name: role.Name,

// 		CreatedAt    time.Time: role.CreatedAt    time.Time,
// 		UpdatedAt    time.Time: role.UpdatedAt    time.Time,
// 	}, nil
// }

// // GetRoleByName gets a role by name
// func (s *RoleService) GetRoleByName(ctx context.Context, req *v1.GetRoleByNameRequest) (*v1.Role, error) {
// 	if req.Name == "" {
// 		return nil, status.Error(codes.InvalidArgument, "name is required")
// 	}

// 	role, err := s.useCase.GetRoleByName(ctx, req.Name)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &v1.Role{
// 		Id:   role.ID,
// 		Name: role.Name,

// 		CreatedAt    time.Time: role.CreatedAt    time.Time,
// 		UpdatedAt    time.Time: role.UpdatedAt    time.Time,
// 	}, nil
// }

// // GetRoles gets all roles with pagination
// func (s *RoleService) GetRoles(ctx context.Context, req *v1.GetRolesRequest) (*v1.RoleListResponse, error) {
// 	limit := int(req.Limit)
// 	offset := int(req.Offset)
// 	if limit <= 0 {
// 		limit = 10
// 	}
// 	if offset < 0 {
// 		offset = 0
// 	}

// 	roles, err := s.useCase.GetRoles(ctx, offset, limit)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	response := &v1.RoleListResponse{
// 		Total: int32(roles.Count),
// 		Roles: make([]*v1.Role, 0, len(roles.Items)),
// 	}

// 	for _, role := range roles.Items {
// 		response.Roles = append(response.Roles, &v1.Role{
// 			Id:   role.ID,
// 			Name: role.Name,

// 			CreatedAt    time.Time: role.CreatedAt    time.Time,
// 			UpdatedAt    time.Time: role.UpdatedAt    time.Time,
// 		})
// 	}

// 	return response, nil
// }

// // UpdateRole updates a role
// func (s *RoleService) UpdateRole(ctx context.Context, req *v1.UpdateRoleRequest) (*v1.Empty, error) {
// 	if req.Role == nil {
// 		return nil, status.Error(codes.InvalidArgument, "role is required")
// 	}

// 	err := s.useCase.UpdateRole(ctx, entity.UpdateRoleRequest{
// 		ID:   req.Role.Id,
// 		Name: req.Role.Name,
// 	})
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &v1.Empty{}, nil
// }

// // DeleteRole deletes a role
// func (s *RoleService) DeleteRole(ctx context.Context, req *v1.DeleteRoleRequest) (*v1.Empty, error) {
// 	if req.Id == "" {
// 		return nil, status.Error(codes.InvalidArgument, "id is required")
// 	}

// 	err := s.useCase.DeleteRole(ctx, req.Id)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &v1.Empty{}, nil
// }
