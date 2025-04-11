package service

// import (
// 	"context"

// 	v1 "sugurta/genproto/auth_service"
// 	"sugurta/internal/entity"
// 	"sugurta/internal/usecase"

// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// )

// // UserService is implementation of UserServiceServer
// type UserService struct {
// 	v1.UnimplementedUserServiceServer
// 	useCase usecase.UseCase
// }

// // NewUserService creates a new UserService
// func NewUserService(useCase usecase.UseCase) *UserService {
// 	return &UserService{
// 		useCase: useCase,
// 	}
// }

// // CreateUser creates a new user
// func (s *UserService) CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*v1.User, error) {
// 	if req.User == nil {
// 		return nil, status.Error(codes.InvalidArgument, "user is required")
// 	}

// 	user := entity.User{
// 		FirstName:   req.User.FirstName,
// 		LastName:    req.User.LastName,
// 		Email:       req.User.Email,
// 		PhoneNumber: req.User.PhoneNumber,
// 		Password:    req.User.Password,
// 		IsActive:    req.User.IsActive,
// 		RoleID:      req.User.RoleId,
// 	}

// 	result, err := s.useCase.CreateUser(ctx, user)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	response := &v1.User{
// 		Id:          result.ID,
// 		FirstName:   result.FirstName,
// 		LastName:    result.LastName,
// 		Email:       result.Email,
// 		PhoneNumber: result.PhoneNumber,
// 		IsActive:    result.IsActive,
// 		RoleId:      result.RoleID,
// 		CreatedAt    time.Time:   result.CreatedAt    time.Time,
// 		UpdatedAt    time.Time:   result.UpdatedAt    time.Time,
// 	}

// 	if result.RoleData.ID != "" {
// 		response.RoleData = &v1.Role{
// 			Id:        result.RoleData.ID,
// 			Name:      result.RoleData.Name,
// 			CreatedAt    time.Time: result.RoleData.CreatedAt    time.Time,
// 			UpdatedAt    time.Time: result.RoleData.UpdatedAt    time.Time,
// 		}
// 	}

// 	return response, nil
// }

// // GetUserByID gets a user by ID
// func (s *UserService) GetUserByID(ctx context.Context, req *v1.GetUserByIDRequest) (*v1.User, error) {
// 	if req.Id == "" {
// 		return nil, status.Error(codes.InvalidArgument, "id is required")
// 	}

// 	user, err := s.useCase.GetUserByEmail(ctx, req.Id)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	response := &v1.User{
// 		Id:          user.ID,
// 		FirstName:   user.FirstName,
// 		LastName:    user.LastName,
// 		Email:       user.Email,
// 		PhoneNumber: user.PhoneNumber,
// 		IsActive:    user.IsActive,
// 		RoleId:      user.RoleID,
// 		CreatedAt    time.Time:   user.CreatedAt    time.Time,
// 		UpdatedAt    time.Time:   user.UpdatedAt    time.Time,
// 	}

// 	if user.RoleData.ID != "" {
// 		response.RoleData = &v1.Role{
// 			Id:          user.RoleData.ID,
// 			Name:        user.RoleData.Name,

// 			CreatedAt    time.Time:   user.RoleData.CreatedAt    time.Time,
// 			UpdatedAt    time.Time:   user.RoleData.UpdatedAt    time.Time,
// 		}
// 	}

// 	return response, nil
// }

// // GetUserByEmail gets a user by email
// func (s *UserService) GetUserByEmail(ctx context.Context, req *v1.GetUserByEmailRequest) (*v1.User, error) {
// 	if req.Email == "" {
// 		return nil, status.Error(codes.InvalidArgument, "email is required")
// 	}

// 	user, err := s.useCase.GetUserByEmail(ctx, req.Email)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	response := &v1.User{
// 		Id:          user.ID,
// 		FirstName:   user.FirstName,
// 		LastName:    user.LastName,
// 		Email:       user.Email,
// 		PhoneNumber: user.PhoneNumber,
// 		IsActive:    user.IsActive,
// 		RoleId:      user.RoleID,
// 		CreatedAt    time.Time:   user.CreatedAt    time.Time,
// 		UpdatedAt    time.Time:   user.UpdatedAt    time.Time,
// 	}

// 	if user.RoleData.ID != "" {
// 		response.RoleData = &v1.Role{
// 			Id:          user.RoleData.ID,
// 			Name:        user.RoleData.Name,

// 			CreatedAt    time.Time:   user.RoleData.CreatedAt    time.Time,
// 			UpdatedAt    time.Time:   user.RoleData.UpdatedAt    time.Time,
// 		}
// 	}

// 	return response, nil
// }

// // UpdatePassword updates a user's password
// func (s *UserService) UpdatePassword(ctx context.Context, req *v1.UpdatePasswordRequest) (*v1.Empty, error) {
// 	if req.Email == "" || req.Password == "" {
// 		return nil, status.Error(codes.InvalidArgument, "email and password are required")
// 	}

// 	err := s.useCase.UpdatePassword(ctx, entity.UpdatePassword{
// 		Email:    req.Email,
// 		Password: req.Password,
// 	})
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &v1.Empty{}, nil
// }

// // UpdateUser updates a user
// func (s *UserService) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (*v1.Empty, error) {
// 	if req.User == nil {
// 		return nil, status.Error(codes.InvalidArgument, "user is required")
// 	}

// 	err := s.useCase.UpdateUser(ctx, entity.User{
// 		ID:          req.User.Id,
// 		FirstName:   req.User.FirstName,
// 		LastName:    req.User.LastName,
// 		Email:       req.User.Email,
// 		PhoneNumber: req.User.PhoneNumber,
// 		IsActive:    req.User.IsActive,
// 		RoleID:      req.User.RoleId,
// 	})
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &v1.Empty{}, nil
// }

// // DeleteUser deletes a user
// func (s *UserService) DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) (*v1.Empty, error) {
// 	if req.Id == "" {
// 		return nil, status.Error(codes.InvalidArgument, "id is required")
// 	}

// 	err := s.useCase.DeleteUser(ctx, req.Id)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &v1.Empty{}, nil
// }
