package repository

import (
	"sugurta/internal/entity"
	"context"
)

type Client interface {
	Create(ctx context.Context, role *entity.CreateRoleRequest) error
	Get(ctx context.Context, params map[string]string) (*entity.Role, error)
	List(ctx context.Context, limit, offset uint64, filter map[string]string) ([]*entity.Role, error)
	Update(ctx context.Context, role *entity.UpdateRoleRequest) error
	Delete(ctx context.Context, id string) error
}
