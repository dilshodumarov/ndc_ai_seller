package role

import (
	"sugurta/internal/entity"
	"sugurta/internal/infrastructure/repository"
	"context"
	"errors"
	"time"
)

type roleService struct {
	ctxTimeout time.Duration
	roleRepo   repository.Role
}

func NewRoleService(ctxTimeout time.Duration, roleRepo repository.Role) Role {
	return &roleService{
		ctxTimeout: ctxTimeout,
		roleRepo:   roleRepo,
	}
}

// CreateRole creates a new role
func (a *roleService) Create(ctx context.Context, role *entity.CreateRoleRequest) error {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	return a.roleRepo.Create(ctx, role)
}

// GetRole gets a role by id
func (a *roleService) Get(ctx context.Context, params map[string]string) (*entity.Role, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	return a.roleRepo.Get(ctx, params)
}

func (a *roleService) List(ctx context.Context, page, limit uint64, filter map[string]string) (*entity.RoleListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	return a.roleRepo.List(ctx, page, limit, filter)
}

// UpdateRole updates a role
func (a *roleService) Update(ctx context.Context, role *entity.UpdateRoleRequest) error {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	if role.ID == "" {
		return errors.New("id is required")
	}
	if role.Name == "" {
		return errors.New("name is required")
	}
	if role.ClientTypeId == "" {
		return errors.New("client_type_id is required")
	}
	return a.roleRepo.Update(ctx, role)
}

// DeleteRole deletes a role
func (a *roleService) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, a.ctxTimeout)
	defer cancel()

	// // tracing
	// ctx, span := otlp.Start(ctx, "refreshTokenService", "refreshTokenUsecaseGet")
	// defer span.End()

	return a.roleRepo.Delete(ctx, id)
}
