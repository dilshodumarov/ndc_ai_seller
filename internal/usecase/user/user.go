package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"sugurta/internal/entity"
	"sugurta/internal/infrastructure/repository"
	"sugurta/internal/pkg/helper"

	"github.com/google/uuid"
)

// UseCase is the User use case implementation
type userService struct {
	ctxTimeout time.Duration
	userRepo   repository.User
}

func NewUserService(ctxTimeout time.Duration, u repository.User) User {
	return &userService{
		ctxTimeout: ctxTimeout,
		userRepo:   u,
	}
}

// CreateAdmin creates a new admin user
func (u *userService) Create(ctx context.Context, user *entity.User) (*entity.User, error) {

	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	var (
		roleData entity.Role
	)

	// Validate admin data
	if user.Email == "" || user.Password == "" || user.FirstName == "" {
		return nil, errors.New("required fields missing")
	}

	if user.RoleID == "" {

		// get default role
		// WARNING > need to get role data and set role id to user data	for creating

		roleData.ID = uuid.NewString()

		user.RoleID = roleData.ID
	}

	// Create admin in database
	userRes, err := u.userRepo.Create(ctx, user)
	if err != nil {
		// a.log.Error(fmt.Sprintf("CreateAdmin - error creating admin: %v", err))
		return nil, fmt.Errorf("failed to create admin: %w", err)
	}

	userRes.RoleData = roleData

	return userRes, nil
}

// GetUserByEmail gets a user by email
func (u *userService) Get(ctx context.Context, params map[string]string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	admin, err := u.userRepo.Get(ctx, params)
	if err != nil {
		// a.log.Error(fmt.Sprintf("GetUserByEmail - error getting user by email: %v", err))
		return nil, err
	}

	return admin, nil
}

func (u *userService) GetByIDs(ctx context.Context, id string) ([]*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	admin, err := u.userRepo.GetByIDs(ctx, id)
	if err != nil {
		// a.log.Error(fmt.Sprintf("GetUserByEmail - error getting user by email: %v", err))
		return nil, err
	}

	return admin, nil
}

// GetUserByEmail gets a user by email
func (u *userService) List(ctx context.Context, filter entity.UserFilter) ([]*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	user, err := u.userRepo.List(ctx, filter)
	if err != nil {
		// a.log.Error(fmt.Sprintf("GetUserByID - error getting user by id: %v", err))
		return nil, err
	}

	return user, nil
}

// UpdateUser updates a user
func (u *userService) Update(ctx context.Context, user *entity.User) error {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	return u.userRepo.Update(ctx, user)
}

func (u *userService) UpdatePassword(ctx context.Context, user *entity.UpdatePasswordRequest) error {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	return u.userRepo.UpdatePassword(ctx, user)
}

// VerifyUser verifies a user with a verification code
func (u *userService) VerifyUser(ctx context.Context, email, code string) error {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	// Validate input
	if email == "" || code == "" {
		return errors.New("email and code are required")
	}

	// In a real implementation, we would check the code against
	// a saved code in the cache or database
	// This is done in the handler level with the in-memory cache

	return nil
}

// Login Userenticates a user and returns user data
func (u *userService) Login(ctx context.Context, email, password string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	// Validate input
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	// Get user by email
	user, err := u.userRepo.Get(ctx, map[string]string{
		"email": email,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid email or password")
		}
		// a.log.Error(fmt.Sprintf("Login - error retrieving user: %v", err))
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	// Check password
	err = helper.CheckPassword(password, user.Password)
	if err != nil {
		// a.log.Warn(fmt.Sprintf("Login - invalid password for user %s", email))
		return nil, errors.New("invalid email or password")
	}

	// Return admin data with password cleared
	user.Password = ""
	return user, nil
}

// DeleteUser deletes a user
func (u *userService) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	return u.userRepo.Delete(ctx, id)
}

func (u *userService) ListClients(ctx context.Context, filter entity.ClientFilter) (*entity.ListClients, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	user, err := u.userRepo.ListClients(ctx, filter)
	if err != nil {
		// a.log.Error(fmt.Sprintf("GetUserByID - error getting user by id: %v", err))
		return nil, err
	}

	return user, nil
}

func (u *userService) GetClientByID(ctx context.Context, id string) (*entity.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	user, err := u.userRepo.GetClientByID(ctx, id)
	if err != nil {
		// a.log.Error(fmt.Sprintf("GetUserByID - error getting user by id: %v", err))
		return nil, err
	}
	return user, nil
}

func (u *userService) BlockUser(ctx context.Context, req entity.BlockUser) error {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	err := u.userRepo.BlockUser(ctx, req)
	if err != nil {
		// a.log.Error(fmt.Sprintf("GetUserByID - error getting user by id: %v", err))
		return err
	}
	return nil
}

func (u *userService) PauzChat(ctx context.Context, req entity.PauzeChat) error {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	err := u.userRepo.PauzChat(ctx, req)
	if err != nil {
		// a.log.Error(fmt.Sprintf("GetUserByID - error getting user by id: %v", err))
		return err
	}
	return nil
}
