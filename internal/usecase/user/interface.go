package user

import (
	"sugurta/internal/entity"
	"context"
)

type User interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	Get(ctx context.Context, params map[string]string) (*entity.User, error)
	List(ctx context.Context, limit, offset uint64, filter map[string]string) ([]*entity.User, error)
	Update(ctx context.Context, investor *entity.User) error
	UpdatePassword(ctx context.Context, investor *entity.UpdatePasswordRequest) error
	Delete(ctx context.Context, guid string) error
}
